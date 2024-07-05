package xmlx

import (
	"bytes"
	"fmt"

	mxj "github.com/clbanning/mxj/v2"
	"github.com/kesonan/xlib/pkg/converter/anyx"
	"github.com/kesonan/xlib/pkg/converter/constx"
)

func Convert(s, toType string) (string, error) {
	mv, err := mxj.NewMapXmlReader(bytes.NewBufferString(s), true)
	if err != nil {
		return "", err
	}
	var v map[string]any
	err = mv.Struct(&v)
	if err != nil {
		return "", err
	}

	switch toType {
	case constx.TypeJSON:
		return anyx.ConvertToJson(v)
	case constx.TypeTOML:
		return anyx.ConvertToToml(v)
	case constx.TypeGoStruct:
		return anyx.ConvertToGoStruct(v, true)
	case constx.TypeYAML:
		return anyx.ConvertToYaml(v)
	case constx.TypeSQL:
		var value any
		for _, item := range v {
			value = item
			break
		}
		_, ok := value.(map[string]any)
		if !ok {
			return "", fmt.Errorf("invalid xml data")
		}
		return anyx.ConvertToSQL(value)
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
