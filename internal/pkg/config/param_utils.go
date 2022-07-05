package config

import (
	options "github.com/breathbath/go_utils/v2/pkg/config"
)

func ReadAPIUser(params *options.ParameterBag) (user string) {
	user = params.ReadString(APIUser, "")
	if user == "" {
		user = params.ReadString(Login, "")
	}
	return user
}

func ReadAPIPassword(params *options.ParameterBag) (pass string) {
	pass = params.ReadString(APIPassword, "")
	if pass == "" {
		pass = params.ReadString(Password, "")
	}
	return pass
}

func ReadAPIURL(params *options.ParameterBag) (url string) {
	url = ReadAPIURLWithDefault(params, DefaultServerURL)
	return url
}

func ReadAPIURLWithDefault(params *options.ParameterBag, defaultURL string) (url string) {
	url = params.ReadString(APIURL, "")
	if url == "" {
		url = params.ReadString(ServerURL, defaultURL)
	}
	return url
}

func GetNoPromptFlagSpec() (flagSpec ParameterRequirement) {
	return ParameterRequirement{
		Field:       NoPrompt,
		Help:        "Flag to disable prompting when missing values",
		Description: "Never prompt when missing parameters",
		ShortName:   "q",
		Type:        BoolRequirementType,
		Default:     false,
	}
}
