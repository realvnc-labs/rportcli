package config

import (
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

type YAMLExecuteParams struct {
	Cids                []string `yaml:"cids,omitempty"`
	Name                []string `yaml:"name,omitempty"`
	Names               []string `yaml:"names,omitempty"`
	Search              string   `yaml:"search,omitempty"`
	Command             string   `yaml:"command,omitempty"`
	EmbeddedScript      string   `yaml:"exec,omitempty"`
	Script              string   `yaml:"script,omitempty"`
	Timeout             string   `yaml:"timeout,omitempty"`
	Gids                []string `yaml:"gids,omitempty"`
	Conc                bool     `yaml:"conc,omitempty"`
	FullCommandResponse bool     `yaml:"full-command-response,omitempty"`
	IsSudo              bool     `yaml:"is_sudo,omitempty"`
	Interpreter         string   `yaml:"interpreter,omitempty"`
	AbortOnError        bool     `yaml:"abort,omitempty"`
	Cwd                 string   `yaml:"cwd,omitempty"`
	WriteExecLog        string   `yaml:"writeexeclog,omitempty"`
	ReadExecLog         string   `yaml:"readexeclog,omitempty"`
}

const (
	expectedMaxYAMLParams = 32
)

type UsedFlagsChecker interface {
	ChangedFlag(flagName string) (isFound bool)
}

func GetParamName(paramTag reflect.StructTag) (paramName string) {
	yamlTag := paramTag.Get("yaml")
	tagParts := strings.Split(yamlTag, ",")
	paramName = tagParts[0]
	return paramName
}

func convertToDelimitedString(paramStrings []string) (paramString string) {
	return strings.Join(paramStrings, ",")
}

func GetExecuteParams(ep *YAMLExecuteParams, yFileParams map[string]interface{}, flagsChecker UsedFlagsChecker) map[string]interface{} {
	e := reflect.ValueOf(ep).Elem()

	for i := 0; i < e.NumField(); i++ {
		paramTag := e.Type().Field(i).Tag
		paramName := GetParamName(paramTag)

		// if the param was changed on the command line then ignore in the yaml
		if flagsChecker != nil {
			isFound := flagsChecker.ChangedFlag(paramName)
			if isFound {
				continue
			}
		}

		paramType := e.Field(i).Type()
		paramValue := e.Field(i).Interface()

		if paramType == reflect.TypeOf([]string{}) {
			paramStrings := paramValue.([]string)
			if len(paramStrings) > 0 {
				// convert to comma delimited string, rather than array
				yFileParams[paramName] = convertToDelimitedString(paramStrings)
			}
		} else {
			if paramType == reflect.TypeOf(string("")) {
				trimmedValue := strings.TrimSpace(paramValue.(string))
				if trimmedValue != "" {
					yFileParams[paramName] = trimmedValue
				}
			} else {
				// will set non string types
				yFileParams[paramName] = paramValue
			}
		}
	}
	return yFileParams
}

func ReadYAMLExecuteParams(fileList []string, flagsChecker UsedFlagsChecker) (yFileParams map[string]interface{}, err error) {
	yFileParams = make(map[string]interface{}, expectedMaxYAMLParams)

	for _, filename := range fileList {
		f := strings.TrimSpace(filename)

		contents, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}

		executeParams := &YAMLExecuteParams{}
		err = yaml.Unmarshal(contents, executeParams)
		if err != nil {
			return nil, err
		}
		yFileParams = GetExecuteParams(executeParams, yFileParams, flagsChecker)
	}
	return yFileParams, nil
}
