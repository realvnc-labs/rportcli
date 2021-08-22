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

	options "github.com/breathbath/go_utils/v2/pkg/config"
	io2 "github.com/breathbath/go_utils/v2/pkg/io"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/sirupsen/logrus"
)

const (
	DefaultCmdTimeoutSeconds = 30
	ClientIDs                = "cids"
	Command                  = "command"
	Script                   = "script"
	GroupIDs                 = "gids"
	Timeout                  = "timeout"
	ExecConcurrently         = "conc"
	AbortOnError             = "abort"
	Cwd                      = "cwd"
	IsSudo                   = "is_sudo"
	Interpreter              = "interpreter"
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

type ExecutionHelper struct {
	ClientSearch ClientSearch
	JobRenderer  JobRenderer
	ReadWriter   ReadWriter
}

func (eh *ExecutionHelper) execute(ctx context.Context, params *options.ParameterBag, scriptPayload, interpreter string) error {
	if eh.ReadWriter != nil {
		defer io2.CloseResourceSecure("read writer", eh.ReadWriter)
	}

	clientIDs, err := eh.getClientIDs(ctx, params)
	if err != nil {
		return err
	}

	wsCmd := eh.buildExecInput(params, clientIDs, scriptPayload, interpreter)
	err = eh.sendCommand(wsCmd)
	if err != nil {
		return err
	}

	err = eh.startReading(ctx)

	return err
}

func (eh *ExecutionHelper) buildExecInput(
	params *options.ParameterBag,
	clientIDs, scriptPayload, interpreter string,
) *models.WsScriptCommand {
	wsCmd := &models.WsScriptCommand{
		ClientIDs:           strings.Split(clientIDs, ","),
		TimeoutSec:          params.ReadInt(Timeout, DefaultCmdTimeoutSeconds),
		ExecuteConcurrently: params.ReadBool(ExecConcurrently, false),
		GroupIDs:            nil,
		AbortOnError:        params.ReadBool(AbortOnError, false),
		Cwd:                 params.ReadString(Cwd, ""),
		IsSudo:              params.ReadBool(IsSudo, false),
		Interpreter:         interpreter,
	}

	if scriptPayload != "" {
		wsCmd.Script = scriptPayload
	} else {
		wsCmd.Command = params.ReadString(Command, "")
	}

	groupIDsStr := params.ReadString(GroupIDs, "")
	if groupIDsStr != "" {
		groupIDsList := strings.Split(groupIDsStr, ",")
		wsCmd.GroupIDs = groupIDsList
	}

	return wsCmd
}

func (eh *ExecutionHelper) sendCommand(wsCmd *models.WsScriptCommand) error {
	wsCmdJSON, err := json.Marshal(wsCmd)
	if err != nil {
		return err
	}
	logrus.Debugf("will send %s", string(wsCmdJSON))

	_, err = eh.ReadWriter.Write(wsCmdJSON)
	if err != nil {
		return err
	}

	return nil
}

func (eh *ExecutionHelper) getClientIDs(ctx context.Context, params *options.ParameterBag) (clientIDs string, err error) {
	clientIDs = params.ReadString(ClientIDs, "")
	clientName := params.ReadString(ClientNameFlag, "")

	if clientIDs == "" && clientName == "" {
		return "", errors.New("no client id nor name provided")
	}

	if clientIDs == "" {
		clients, err := eh.ClientSearch.Search(ctx, clientName, params)
		if err != nil {
			return "", err
		}

		if len(clients) == 0 {
			return "", fmt.Errorf("unknown client(s) '%s'", clientName)
		}

		for i := range clients {
			cl := clients[i]
			clientIDs += cl.ID + ","
		}

		clientIDs = strings.Trim(clientIDs, ",")
	}

	return clientIDs, nil
}

func (eh *ExecutionHelper) startReading(ctx context.Context) error {
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
				msg, err := eh.ReadWriter.Read()
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
			err := eh.processRawMessage(msg)
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

func (eh *ExecutionHelper) processRawMessage(msg []byte) error {
	var job models.Job
	err := json.Unmarshal(msg, &job)
	if err != nil || job.Jid == "" {
		logrus.Debugf("cannot unmarshal '%s' to the Job: %v, will try interpret it as an error", string(msg), err)
		var errResp models.ErrorResp
		err = json.Unmarshal(msg, &errResp)
		if err != nil {
			e := fmt.Errorf("cannot recognize command output message: %s, reason: %v", string(msg), err)
			return e
		}
		return errResp
	}

	logrus.Debugf("received message: '%s'", string(msg))

	err = eh.JobRenderer.RenderJob(&job)
	return err
}
