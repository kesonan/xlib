package vars

import (
	"time"

	"github.com/kesonan/xlib/pkg/converter/constx"
	"github.com/zeromicro/ddl-parser/parser"
)

var (
	MayTimeColumn = []string{"CREATEDAT", "CREATEAT", "UPDATEDAT", "UPDATEdAT", "DELETEDAT", "DELETEAT", "DATE", "TIME"}
	SqlDefault    = map[string]string{
		"string":  `''`,
		"int":     "0",
		"float64": "0.00",
		"bool":    "0",
	}
	SqlTypeFromGo = map[string]string{
		"string":    "varchar(255)",
		"int":       "int(11)",
		"float64":   "decimal",
		"bool":      "tinyint(1)",
		"time.Time": "timestamp",
	}

	SqlTypeToGo = map[int]any{
		parser.Bit:       0,
		parser.Timestamp: time.Now(),
		parser.DateTime:  time.Now(),
		parser.Date:      time.Now(),
		parser.Year:      "",
		parser.Decimal:   0.0,
		parser.Dec:       0.0,
		parser.Numeric:   0.0,
		parser.Float:     0.0,
		parser.Float4:    0.0,
		parser.Float8:    0.0,
		parser.Double:    0.0,
		parser.TinyInt:   0,
		parser.SmallInt:  0,
		parser.MediumInt: 0,
		parser.Int:       0,
		parser.Integer:   0,
		parser.BigInt:    0,
		parser.Int1:      0,
		parser.Int2:      0,
		parser.Int3:      0,
		parser.Int4:      0,
	}

	ValidMapper = map[string][]string{
		constx.TypeJSON: {constx.TypeGoStruct, constx.TypeProtoBuf, constx.TypeGoctlAPI, constx.TypeTOML, constx.TypeYAML, constx.TypeSQL, constx.TypeXml},
		constx.TypeSQL:  {constx.TypeJSON, constx.TypeGoStruct, constx.TypeProtoBuf, constx.TypeGoctlAPI, constx.TypeTOML, constx.TypeYAML, constx.TypeXml},
		constx.TypeTOML: {constx.TypeJSON, constx.TypeGoStruct, constx.TypeProtoBuf, constx.TypeGoctlAPI, constx.TypeYAML, constx.TypeSQL, constx.TypeXml},
		constx.TypeYAML: {constx.TypeJSON, constx.TypeGoStruct, constx.TypeProtoBuf, constx.TypeGoctlAPI, constx.TypeTOML, constx.TypeSQL, constx.TypeXml},
		constx.TypeXml:  {constx.TypeJSON, constx.TypeGoStruct, constx.TypeProtoBuf, constx.TypeGoctlAPI, constx.TypeYAML, constx.TypeTOML, constx.TypeSQL},
	}
)
