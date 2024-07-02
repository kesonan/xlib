package jsonx

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kesonan/xlib/pkg/converter/anyx"
	"github.com/kesonan/xlib/pkg/converter/internal/constx"
)

func Convert(s string, toType string) (string, error) {
	var v any
	err := json.Unmarshal([]byte(s), &v)
	var jsonSyntaxErr *json.SyntaxError
	if errors.As(err, &jsonSyntaxErr) {
		return "", fmt.Errorf("invalid json, offset: %d, msg: %s", jsonSyntaxErr.Offset, jsonSyntaxErr.Error())
	}

	switch toType {
	case constx.TypeTOML:
		return anyx.ConvertToToml(v)
	case constx.TypeGoStruct:
		return anyx.ConvertToGoStruct(v, true)
	case constx.TypeYAML:
		return anyx.ConvertToYaml(v)
	case constx.TypeSQL:
		return anyx.ConvertToSQL(v)
	case constx.TypeGoctlAPI:
		tp, _, err := anyx.ConvertToGoctlAPI(0, "", "GeneratedResponse", v, true)
		return tp, err
	case constx.TypeProtoBuf:
		val, _, err := anyx.ConvertToProtoBuf(0, "MessageResponse", v, true)
		return val, err
	default:
		return "", fmt.Errorf("invalid to type: %s", toType)
	}
}
