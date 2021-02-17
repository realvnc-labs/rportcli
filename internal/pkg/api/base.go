package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/breathbath/go_utils/utils/env"

	"github.com/sirupsen/logrus"
)

const (
	maxValidResponseCode = 399
	connectionTimeoutSec = 10
)

type Rport struct {
	BaseURL string
	Auth    Auth
}

func New(baseURL string, a Auth) *Rport {
	return &Rport{BaseURL: baseURL, Auth: a}
}

// BaseClient responsible for calling rport API
type BaseClient struct {
	auth Auth
}

func (c *BaseClient) WithAuth(a Auth) {
	c.auth = a
}

func (c *BaseClient) buildClient() *http.Client {
	connectionTimeout := env.ReadEnvInt("CONN_TIMEOUT_SEC", connectionTimeoutSec)
	transport := &http.Transport{
		DisableKeepAlives:     true,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: false},
		ResponseHeaderTimeout: time.Duration(connectionTimeout) * time.Second,
	}
	cl := &http.Client{Transport: transport}

	return cl
}

func (c *BaseClient) Call(req *http.Request, target interface{}) (resp *http.Response, err error) {
	cl := c.buildClient()
	dump, _ := httputil.DumpRequest(req, true)
	logrus.Debugf("raw request: %s", string(dump))

	if c.auth != nil {
		err = c.auth.AuthRequest(req)
		if err != nil {
			return nil, err
		}
	}

	resp, err = cl.Do(req)
	if err != nil {
		return resp, fmt.Errorf("request failed with error: %v", err)
	}

	if target == nil {
		return resp, nil
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err == io.EOF {
		return resp, fmt.Errorf("empty body in the response, status: %d", resp.StatusCode)
	}
	if err != nil {
		return resp, fmt.Errorf("reading of the request body failed with error: %v, status: %d", err, resp.StatusCode)
	}

	logrus.Debugf("Got response: '%s', status code: '%d'", string(respBody), resp.StatusCode)

	if resp.StatusCode > maxValidResponseCode {
		var errResp ErrorResp
		err = json.Unmarshal(respBody, &errResp)
		if err != nil {
			logrus.Warnf("cannot unmarshal error response %s: %v", string(respBody), err)
		}
		return resp, errResp
	}

	err = json.Unmarshal(respBody, target)
	if err != nil {
		return resp, fmt.Errorf("cannot unmarshal response %s to %+v: %v", string(respBody), target, err)
	}

	return resp, nil
}

func closeRespBody(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	closeErr := resp.Body.Close()
	if closeErr != nil {
		logrus.Error(closeErr)
	}
}
