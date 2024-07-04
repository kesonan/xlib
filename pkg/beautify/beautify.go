package beautify

import (
	"bytes"
	"encoding/json"

	"github.com/kesonan/xlib/pkg/converter/constx"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

func JSON(s string) string {
	var v any
	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		return s
	}
	writer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", constx.Indent)
	encoder.SetEscapeHTML(true)
	err = encoder.Encode(v)
	if err != nil {
		return s
	}
	return writer.String()
}

func YAML(s string) string {
	var v any
	err := yaml.Unmarshal([]byte(s), &v)
	if err != nil {
		return s
	}
	writer := bytes.NewBuffer(nil)
	encoder := yaml.NewEncoder(writer)
	defer encoder.Close()
	encoder.SetIndent(2)
	err = encoder.Encode(v)
	if err != nil {
		return s
	}
	return writer.String()
}

func TOML(s string) string {
	var v any
	err := toml.Unmarshal([]byte(s), &v)
	if err != nil {
		return s
	}
	writer := bytes.NewBuffer(nil)
	encoder := toml.NewEncoder(writer)
	encoder.SetIndentSymbol(constx.Indent)
	encoder.SetIndentTables(true)
	encoder.SetMarshalJsonNumbers(true)
	err = encoder.Encode(v)
	if err != nil {
		return ""
	}
	return writer.String()
}
