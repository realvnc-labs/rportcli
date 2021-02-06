package utils

import (
	"net/url"
	"strings"
)

func RemovePortFromURL(rawURL string) string {
	urlToUse := rawURL
	if !strings.Contains(urlToUse, "://") {
		urlToUse = "http://" + urlToUse
	}
	u, err := url.Parse(urlToUse)
	if err != nil {
		return rawURL
	}
	port := u.Port()
	if port == "" {
		return rawURL
	}

	return strings.Replace(rawURL, ":"+port, "", 1)
}
