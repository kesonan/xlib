package tomlx

import (
	"errors"
	"fmt"

	"github.com/kesonan/xlib/pkg/converter/anyx"
	"github.com/kesonan/xlib/pkg/converter/constx"
	"github.com/pelletier/go-toml/v2"
)

func Convert(s string, toType string) (string, error) {
	var v any
	err := toml.Unmarshal([]byte(s), &v)
	var syntaxErr *toml.DecodeError
	if err != nil && errors.As(err, &syntaxErr) {
		return "", errors.New(syntaxErr.String())
	}

	switch toType {
	case constx.TypeJSON:
		return anyx.ConvertToJson(v)
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
