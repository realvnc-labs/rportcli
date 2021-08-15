package controllers

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"

	options "github.com/breathbath/go_utils/v2/pkg/config"
)

var fileExtInterpreterMap = map[string]string{
	".ps1": "powershell",
	".bat": "cmd",
}

type ScriptsController struct {
	*ExecutionHelper
}

func (cc *ScriptsController) Start(ctx context.Context, params *options.ParameterBag) error {
	scriptsFilePath, err := params.ReadRequiredString(Script)
	if err != nil {
		return err
	}

	info, err := os.Stat(scriptsFilePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("script file doesn't exist: %s", scriptsFilePath)
	}
	if info.IsDir() {
		return fmt.Errorf("script file %s is a directory", scriptsFilePath)
	}

	file, err := os.Open(scriptsFilePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", scriptsFilePath, err)
	}

	buf := &bytes.Buffer{}
	enc := base64.NewEncoder(base64.StdEncoding, buf)
	_, err = io.Copy(enc, file)
	if err != nil {
		return fmt.Errorf("failed to encode file %s to base64: %w", scriptsFilePath, err)
	}

	interpreter := cc.resolveInterpreterByFileName(scriptsFilePath, params.ReadString(Interpreter, ""))

	return cc.execute(ctx, params, buf.String(), interpreter)
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
