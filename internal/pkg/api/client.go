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

	"github.com/sirupsen/logrus"
)

const (
	maxValidResponseCode = 399
)

type Rport struct {
	BaseURL string
	Auth    Auth
}

func New(baseURL string, a Auth) *Rport {
	return &Rport{BaseURL: baseURL, Auth: a}
}

// Client responsible for calling rport API
type Client struct {
	auth Auth
}

func (c *Client) WithAuth(a Auth) {
	c.auth = a
}

func (c *Client) Call(req *http.Request, target interface{}) error {
	connectionTimeout := 30 * time.Second
	transport := &http.Transport{
		DisableKeepAlives:     true,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: false},
		ResponseHeaderTimeout: connectionTimeout,
	}
	client := http.Client{Transport: transport}
	dump, _ := httputil.DumpRequest(req, true)
	logrus.Debugf("raw request: %s", string(dump))

	if c.auth != nil {
		err := c.auth.AuthRequest(req)
		if err != nil {
			return err
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		if resp != nil && resp.Body != nil {
			closeErr := resp.Body.Close()
			if closeErr != nil {
				logrus.Error(closeErr)
			}
		}
		return fmt.Errorf("request failed with error: %v", err)
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			logrus.Error(closeErr)
		}
	}()
	if target == nil {
		return nil
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err == io.EOF {
		return fmt.Errorf("empty body in the response, status: %d", resp.StatusCode)
	}
	if err != nil {
		return fmt.Errorf("reading of the request body failed with error: %v, status: %d", err, resp.StatusCode)
	}

	logrus.Debugf("Got response: '%s', status code: '%d'", string(respBody), resp.StatusCode)

	if resp.StatusCode > maxValidResponseCode {
		var errResp ErrorResp
		err = json.Unmarshal(respBody, &errResp)
		if err != nil {
			logrus.Warnf("cannot unmarshal error response %s: %v", string(respBody), err)
		}
		return errResp
	}

	err = json.Unmarshal(respBody, target)
	if err != nil {
		return fmt.Errorf("cannot unmarshal response %s to %+v: %v", string(respBody), target, err)
	}

	return nil
}
