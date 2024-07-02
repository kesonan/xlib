package vars

import "github.com/kesonan/xlib/pkg/converter/internal/constx"

var (
	MayTimeColumn = []string{"CREATEDAT", "CREATEAT", "UPDATEDAT", "UPDATEdAT", "DELETEDAT", "DELETEAT", "DATE", "TIME"}
	SqlDefault    = map[string]string{
		"string":  `''`,
		"int":     "0",
		"float64": "0.00",
		"bool":    "0",
	}
	SqlType = map[string]string{
		"string":  "varchar(255)",
		"int":     "int(11)",
		"float64": "decimal",
		"bool":    "tinyint(1)",
	}

	ValidMapper = map[string][]string{
		constx.TypeJSON: {constx.TypeGoStruct, constx.TypeProtoBuf, constx.TypeGoctlAPI, constx.TypeTOML, constx.TypeYAML, constx.TypeSQL},
		constx.TypeSQL:  {constx.TypeJSON, constx.TypeGoStruct, constx.TypeProtoBuf, constx.TypeGoctlAPI, constx.TypeTOML, constx.TypeYAML},
		constx.TypeTOML: {constx.TypeJSON, constx.TypeGoStruct, constx.TypeProtoBuf, constx.TypeGoctlAPI, constx.TypeYAML},
		constx.TypeYAML: {constx.TypeJSON, constx.TypeGoStruct, constx.TypeProtoBuf, constx.TypeGoctlAPI, constx.TypeTOML},
	}
)
