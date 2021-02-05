package api

import (
	"net/http"

	http2 "github.com/breathbath/go_utils/utils/http"
)

type Auth interface {
	AuthRequest(r *http.Request) error
}

type BasicAuth struct {
	Login string
	Pass  string
}

func (ba *BasicAuth) AuthRequest(req *http.Request) error {
	basicAuthHeader := http2.BuildBasicAuthString(ba.Login, ba.Pass)
	req.Header.Add("Authorization", "Basic "+basicAuthHeader)

	return nil
}
