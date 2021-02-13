package utils

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

const ipCheckerURL = "https://api64.ipify.org"

type APIIPProvider struct {
}

func (ap APIIPProvider) GetIP(ctx context.Context) (string, error) {
	logrus.Debugf("will detect public IP from %s", ipCheckerURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ipCheckerURL, nil)
	if err != nil {
		return "", nil
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer func() {
		e := resp.Body.Close()
		if e != nil {
			logrus.Error(e)
		}
	}()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	logrus.Debugf("got Response: %s", string(bodyBytes))

	return string(bodyBytes), nil
}
