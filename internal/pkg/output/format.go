package output

import (
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

const (
	FormatHuman      = "human"
	FormatJSON       = "json"
	FormatJSONPretty = "json-pretty"
	FormatYAML       = "yaml"
)

func RenderByFormat(format string, w io.Writer, source interface{}, renderCallback func() error) error {
	if format == "" {
		format = FormatHuman
	}

	switch format {
	case FormatHuman:
		return renderCallback()
	case FormatJSON:
		return RenderJSON(w, source, false)
	case FormatJSONPretty:
		return RenderJSON(w, source, true)
	case FormatYAML:
		yamlEncoder := yaml.NewEncoder(w)
		yamlEncoder.SetIndent(2)
		return yamlEncoder.Encode(source)
	}

	return fmt.Errorf("unknown rendering format: %s", format)
}

func RenderJSON(w io.Writer, source interface{}, isPretty bool) error {
	jsonEncoder := json.NewEncoder(w)
	if isPretty {
		jsonEncoder.SetIndent("", "  ")
	}

	return jsonEncoder.Encode(source)
}
