package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	io2 "github.com/breathbath/go_utils/v2/pkg/io"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/sirupsen/logrus"
)

const (
	DefaultCmdTimeoutSeconds = 30
	ClientIDs                = "cids"
	Command                  = "command"
	GroupIDs                 = "gids"
	Timeout                  = "timeout"
	ExecConcurrently         = "conc"
	IsFullOutput             = "full-command-response"
	waitingMsg               = "waiting for the command to finish"
)

type CliReader interface {
	ReadString() (string, error)
}

type ReadWriter interface {
	Read() (msg []byte, err error)
	Write(inputMsg []byte) (n int, err error)
	io.Closer
}

type Spinner interface {
	Start(msg string)
	Update(msg string)
	StopSuccess(msg string)
	StopError(msg string)
}

type JobRenderer interface {
	RenderJob(j *models.Job) error
}

type CommandsController struct {
	ReadWriter   ReadWriter
	JobRenderer  JobRenderer
	ClientSearch ClientSearch
}

func (icm *CommandsController) Start(ctx context.Context, params *options.ParameterBag) error {
	defer io2.CloseResourceSecure("read writer", icm.ReadWriter)

	clientIDs := params.ReadString(ClientIDs, "")
	clientName := params.ReadString(ClientNameFlag, "")
	if clientIDs == "" && clientName == "" {
		return errors.New("no client id nor name provided")
	}

	if clientIDs == "" {
		clients, err := icm.ClientSearch.Search(ctx, clientName)
		if err != nil {
			return err
		}

		if len(clients) == 0 {
			return fmt.Errorf("unknown client(s) '%s'", clientName)
		}

		for i := range clients {
			cl := clients[i]
			clientIDs += cl.ID + ","
		}

		clientIDs = strings.Trim(clientIDs, ",")
	}

	wsCmd := icm.buildCommand(params, clientIDs)
	err := icm.sendCommand(wsCmd)
	if err != nil {
		return err
	}

	err = icm.startReading(ctx)

	return err
}

func (icm *CommandsController) buildCommand(params *options.ParameterBag, clientIDs string) models.WsCommand {
	wsCmd := models.WsCommand{
		Command:             params.ReadString(Command, ""),
		ClientIds:           strings.Split(clientIDs, ","),
		TimeoutSec:          params.ReadInt(Timeout, DefaultCmdTimeoutSeconds),
		ExecuteConcurrently: params.ReadBool(ExecConcurrently, false),
		GroupIds:            nil,
	}
	groupIDsStr := params.ReadString(GroupIDs, "")
	if groupIDsStr != "" {
		groupIDsList := strings.Split(groupIDsStr, ",")
		wsCmd.GroupIds = &groupIDsList
	}

	return wsCmd
}

func (icm *CommandsController) sendCommand(wsCmd models.WsCommand) error {
	wsCmdJSON, err := json.Marshal(wsCmd)
	if err != nil {
		return err
	}
	logrus.Debugf("will send %s", string(wsCmdJSON))

	_, err = icm.ReadWriter.Write(wsCmdJSON)
	if err != nil {
		return err
	}

	return nil
}

func (icm *CommandsController) startReading(ctx context.Context) error {
	errsChan := make(chan error, 1)
	msgChan := make(chan []byte, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer close(msgChan)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := icm.ReadWriter.Read()
				if err != nil {
					if err == io.EOF {
						return
					}
					errsChan <- err
				}
				msgChan <- msg
			}
		}
	}()

mainLoop:
	for {
		select {
		case <-sigs:
			break mainLoop
		case msg, ok := <-msgChan:
			if !ok {
				return nil
			}
			err := icm.processRawMessage(msg)
			if err != nil {
				return err
			}
			logrus.Debug(waitingMsg)
		case err := <-errsChan:
			return err
		}
	}

	return nil
}

func (icm *CommandsController) processRawMessage(msg []byte) error {
	var job models.Job
	err := json.Unmarshal(msg, &job)
	if err != nil || job.Jid == "" {
		logrus.Debugf("cannot unmarshal %s to the Job: %v, will try interpret it as an error", string(msg), err)
		var errResp models.ErrorResp
		err = json.Unmarshal(msg, &errResp)
		if err != nil {
			e := fmt.Errorf("cannot recognize command output message: %s, reason: %v", string(msg), err)
			return e
		}
		logrus.Error(errResp)
		return errResp
	}

	logrus.Debugf("received message: '%s'", string(msg))

	err = icm.JobRenderer.RenderJob(&job)
	return err
}
