package config

import (
	options "github.com/breathbath/go_utils/v2/pkg/config"
)

func ReadApiUser(params *options.ParameterBag) (user string) {
	user = params.ReadString(ApiUser, "")
	if user == "" {
		user = params.ReadString(Login, "")
	}
	return user
}

func ReadApiPassword(params *options.ParameterBag) (pass string) {
	pass = params.ReadString(ApiPassword, "")
	if pass == "" {
		pass = params.ReadString(Password, "")
	}
	return pass
}

func ReadApiURL(params *options.ParameterBag) (url string) {
	url = ReadApiURLWithDefault(params, DefaultServerURL)
	return url
}

func ReadApiURLWithDefault(params *options.ParameterBag, defaultURL string) (url string) {
	url = params.ReadString(ApiURL, "")
	if url == "" {
		url = params.ReadString(ServerURL, defaultURL)
	}
	return url
}
