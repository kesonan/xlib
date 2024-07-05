package sqlx

import (
	"fmt"
	"os"

	"github.com/kesonan/xlib/pkg/converter/anyx"
	"github.com/kesonan/xlib/pkg/converter/constx"
	"github.com/kesonan/xlib/pkg/converter/vars"
	"github.com/zeromicro/ddl-parser/parser"
)

func Convert(s string, toType string) (string, error) {
	temp, err := os.CreateTemp("", "xlib")
	if err != nil {
		return "", err
	}
	defer func() {
		_ = os.Remove(temp.Name())
		_ = temp.Close()
	}()

	if _, err = temp.WriteString(s); err != nil {
		return "", err
	}

	p := parser.NewParser()
	tables, err := p.From(temp.Name())
	if err != nil {
		return "", err
	}

	if len(tables) == 0 {
		return "", nil
	}
	if len(tables) > 1 {
		return "", fmt.Errorf("only support one table, but got %d", len(tables))
	}
	table := tables[0]
	var v = make(map[string]any)
	for _, col := range table.Columns {
		name := col.Name
		tp, ok := vars.SqlTypeToGo[col.DataType.Type()]
		if !ok {
			tp = "example string"
		}
		v[name] = tp
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
	case constx.TypeGoctlAPI:
		tp, _, err := anyx.ConvertToGoctlAPI(0, "", "GeneratedResponse", v, true)
		return tp, err
	case constx.TypeProtoBuf:
		val, _, err := anyx.ConvertToProtoBuf(0, "MessageResponse", v, true)
		return val, err
	case constx.TypeXml:
		return anyx.ConvertToXml(v)
	default:
		return "", fmt.Errorf("invalid to type: %s", toType)
	}
}
