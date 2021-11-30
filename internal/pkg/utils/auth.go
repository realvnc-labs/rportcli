package utils

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	http2 "github.com/breathbath/go_utils/v2/pkg/http"
)

type Auth interface {
	AuthRequest(r *http.Request) error
}

type StorageBasicAuth struct {
	AuthProvider func() (login, pass string, err error)
}

func (sba *StorageBasicAuth) AuthRequest(req *http.Request) error {
	login, pass, err := sba.AuthProvider()
	if err != nil {
		return err
	}
	if login == "" || pass == "" {
		return fmt.Errorf("no login or password value provided")
	}

	basicAuthHeader := http2.BuildBasicAuthString(login, pass)
	req.Header.Add("Authorization", "Basic "+basicAuthHeader)

	return nil
}

type BearerAuth struct {
	TokenProvider func() (string, error)
}

func (ba *BearerAuth) AuthRequest(req *http.Request) error {
	token, err := ba.TokenProvider()
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	return nil
}

type FallbackAuth struct {
	PrimaryAuth  Auth
	FallbackAuth Auth
}

func (fa *FallbackAuth) AuthRequest(req *http.Request) error {
	err := fa.PrimaryAuth.AuthRequest(req)
	if err == nil {
		return nil
	}

	return fa.FallbackAuth.AuthRequest(req)
}

func ExtractBasicAuthLoginAndPassFromRequest(r *http.Request) (login, pass string, err error) {
	basicAuthHeader := r.Header.Get("Authorization")
	loginPassBase64 := strings.TrimPrefix(basicAuthHeader, "Basic ")

	loginPassBytes, err := base64.StdEncoding.DecodeString(loginPassBase64)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to decode basic auth header from base64")
	}
	loginPass := string(loginPassBytes)

	loginPassParts := strings.Split(loginPass, ":")
	const expectedLoginPassPartsCount = 2
	if len(loginPassParts) < expectedLoginPassPartsCount {
		return "", "", fmt.Errorf("failed to extract login and passwrod from %s", loginPass)
	}

	return loginPassParts[0], strings.Join(loginPassParts[1:], ":"), nil
}
