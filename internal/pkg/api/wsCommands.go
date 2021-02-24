package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/breathbath/go_utils/utils/env"
	"github.com/breathbath/go_utils/utils/url"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
	"github.com/sirupsen/logrus"
)

const (
	defaultTokenValiditySeconds = 10 * 60
	defaultCmdTimeoutSeconds    = 30
	CommandsWSUri               = "/api/v1/ws/commands"
	HTTPSPrefix                 = "https"
	HTTPPrefix                  = "http"
	WssPrefix                   = "wss"
	WsPrefix                    = "ws"
)

type WsCommand struct {
	Command             string    `json:"command"`
	ClientIds           []string  `json:"client_ids"`
	GroupIds            *[]string `json:"group_ids,omitempty"`
	TimeoutSec          int       `json:"timeout_sec"`
	ExecuteConcurrently bool      `json:"execute_concurrently"`
}

type AuthProvider func() (login, pass string, err error)

type WsCommandURLProvider struct {
	AuthProvider AuthProvider
	BaseURL      string
}

func (wup *WsCommandURLProvider) BuildWsURL(ctx context.Context) (wsURL string, err error) {
	tokenValiditySeconds := env.ReadEnvInt("SESSION_VALIDITY_SECONDS", defaultTokenValiditySeconds)
	login, pass, err := wup.AuthProvider()
	if err != nil {
		return "", err
	}

	token, err := wup.getToken(ctx, wup.BaseURL, login, pass, tokenValiditySeconds)
	if err != nil {
		return "", err
	}

	wsURL = wup.buildWsURL(token, wup.BaseURL)

	return
}

func (wup *WsCommandURLProvider) buildWsURL(token LoginResponse, baseURL string) string {
	baseURL = wup.replaceHTTPWithWsProtocolPrefix(baseURL)
	return url.JoinURL(baseURL, CommandsWSUri) + "?access_token=" + token.Data.Token
}

func (wup *WsCommandURLProvider) replaceHTTPWithWsProtocolPrefix(u string) string {
	if strings.HasPrefix(u, HTTPSPrefix) {
		return strings.Replace(u, HTTPSPrefix, WssPrefix, 1)
	}
	if strings.HasPrefix(u, HTTPPrefix) {
		return strings.Replace(u, HTTPPrefix, WsPrefix, 1)
	}

	return u
}

func (wup *WsCommandURLProvider) getToken(
	ctx context.Context,
	baseURL, login, pass string,
	tokenValiditySeconds int,
) (LoginResponse, error) {
	basicAuth := &utils.BasicAuth{
		Login: login,
		Pass:  pass,
	}

	cl := New(baseURL, basicAuth)
	logResp, err := cl.Login(ctx, tokenValiditySeconds)
	if err != nil {
		return logResp, err
	}

	return logResp, nil
}

type CliReader interface {
	ReadString() (string, error)
}

type ReadWriter interface {
	Read() (msg []byte, err error)
	Write(inputMsg []byte) (n int, err error)
}

type InteractiveCommandExecutor struct {
	ReadWriter      ReadWriter
	UserInputReader CliReader
}

func (icm *InteractiveCommandExecutor) Start(ctx context.Context) error {
	fmt.Println("please provide client ids as comma separated values")
	fmt.Print("-> ")
	clientIDs, err := icm.UserInputReader.ReadString()
	if err != nil {
		return err
	}

	fmt.Println("please provide command to execute on the clients")
	fmt.Print("-> ")
	cmd, err := icm.UserInputReader.ReadString()
	if err != nil {
		return err
	}

	fmt.Println("please provide group ids as comma separated values")
	fmt.Print("-> ")
	groupIDsStr, err := icm.UserInputReader.ReadString()
	if err != nil {
		return err
	}

	wsCmd := WsCommand{
		Command:             cmd,
		ClientIds:           strings.Split(clientIDs, ","),
		TimeoutSec:          defaultCmdTimeoutSeconds,
		ExecuteConcurrently: false,
		GroupIds:            nil,
	}
	if groupIDsStr != "" {
		groupIDs := strings.Split(groupIDsStr, ",")
		wsCmd.GroupIds = &groupIDs
	}

	wsCmdJSON, err := json.Marshal(wsCmd)
	if err != nil {
		return err
	}
	logrus.Infof("will execute %s", string(wsCmdJSON))

	_, err = icm.ReadWriter.Write(wsCmdJSON)
	if err != nil {
		return err
	}

	for {
		msg, err := icm.ReadWriter.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		fmt.Println(string(msg))
	}
}
