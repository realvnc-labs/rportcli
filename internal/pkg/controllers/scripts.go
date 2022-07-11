package controllers

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
)

var fileExtInterpreterMap = map[string]string{
	".ps1": "powershell",
	".bat": "cmd",
}

type ScriptsController struct {
	*ExecutionHelper
}

func (cc *ScriptsController) Start(ctx context.Context, params *options.ParameterBag) (err error) {
	var scriptContent []byte
	var interpreter string

	scriptsFilePath := params.ReadString(config.Script, "")
	if scriptsFilePath != "" {
		scriptContent, err = cc.ReadScriptContent(scriptsFilePath)
		if err != nil {
			return err
		}

		interpreter = cc.resolveInterpreterByFileName(scriptsFilePath, params.ReadString(config.Interpreter, ""))
	} else {
		embeddedScriptContent := params.ReadString(config.EmbeddedScript, "")
		if embeddedScriptContent != "" {
			scriptContent = []byte(embeddedScriptContent)
		} else {
			scriptContent = []byte("")
		}
	}

	scriptContentBase64 := base64.StdEncoding.EncodeToString(scriptContent)

	return cc.execute(ctx, params, scriptContentBase64, interpreter)
}

func (cc *ScriptsController) ReadScriptContent(scriptsFilePath string) (scriptContent []byte, err error) {
	info, err := os.Stat(scriptsFilePath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("script file doesn't exist: %s", scriptsFilePath)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("script file %s is a directory", scriptsFilePath)
	}

	scriptFile, err := os.Open(scriptsFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", scriptsFilePath, err)
	}

	scriptContent, err = ioutil.ReadAll(scriptFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", scriptsFilePath, err)
	}

	return scriptContent, nil
}

func (cc *ScriptsController) resolveInterpreterByFileName(scriptFilePath, scriptsFilePathFromArgs string) string {
	if scriptsFilePathFromArgs != "" {
		return scriptsFilePathFromArgs
	}

	extension := filepath.Ext(scriptFilePath)
	if extension == "" {
		return ""
	}

	interpreter, ok := fileExtInterpreterMap[extension]
	if !ok {
		return ""
	}

	return interpreter
}
