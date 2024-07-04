package yamlx

import (
	"fmt"

	"github.com/kesonan/xlib/pkg/converter/anyx"
	"github.com/kesonan/xlib/pkg/converter/constx"
	"gopkg.in/yaml.v3"
)

func Convert(s string, toType string) (string, error) {
	var v any
	err := yaml.Unmarshal([]byte(s), &v)
	if err != nil {
		return "", err
	}

	switch toType {
	case constx.TypeJSON:
		return anyx.ConvertToJson(v)
	case constx.TypeGoStruct:
		return anyx.ConvertToGoStruct(v, true)
	case constx.TypeTOML:
		return anyx.ConvertToToml(v)
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
