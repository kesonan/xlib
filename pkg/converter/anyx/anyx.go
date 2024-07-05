package anyx

import (
	"bytes"
	"encoding/json"
	"fmt"
	goformat "go/format"
	"strings"

	mxj "github.com/clbanning/mxj/v2"
	"github.com/iancoleman/strcase"
	"github.com/kesonan/xlib/pkg/converter/constx"
	"github.com/kesonan/xlib/pkg/converter/types"
	"github.com/kesonan/xlib/pkg/converter/vars"
	"github.com/kesonan/xlib/pkg/parser/api/format"
	"github.com/kesonan/xlib/pkg/sortmap"
	"github.com/kesonan/xlib/pkg/writer"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

func ConvertToJson(v any) (string, error) {
	w := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", constx.Indent)
	encoder.SetEscapeHTML(true)
	err := encoder.Encode(v)
	if err != nil {
		return "", err
	}
	return w.String(), nil
}

func ConvertToToml(v any) (string, error) {
	w := bytes.NewBuffer(nil)
	w.WriteString(fmt.Sprintf("# %s\n\n", constx.HeaderText))
	encoder := toml.NewEncoder(w)
	encoder.SetIndentSymbol(constx.Indent)
	encoder.SetIndentTables(true)
	encoder.SetMarshalJsonNumbers(true)
	err := encoder.Encode(v)
	if err != nil {
		return "", err
	}
	return w.String(), nil
}

func ConvertToYaml(v any) (string, error) {
	w := bytes.NewBuffer(nil)
	w.WriteString(fmt.Sprintf("# %s\n\n", constx.HeaderText))
	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)
	err := encoder.Encode(v)
	if err != nil {
		return "", err
	}
	return w.String(), nil
}

func ConvertToSQL(v any) (string, error) {
	kv, ok := v.(map[string]any)
	if !ok {
		return "", fmt.Errorf("input must be object, got %T", v)
	}

	sm := sortmap.From(kv)
	err := sm.Range(func(_ int, key string, value any) error {
		if !types.IsBasic(value) {
			return fmt.Errorf("value must be basic type, got %T", v)
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	w := writer.New(constx.Indent)
	w.Writef("-- %s\n\n", constx.HeaderText)
	w.WriteStringln("CREATE TABLE IF NOT EXISTS `table_name` (")
	if types.MayContainsPrimary(sm) {
		w.WriteWithIndentStringln("`id` int unsigned NOT NULL AUTO_INCREMENT,")
		sm.Del(constx.MayIdColumn)
	}

	sm.Range(func(_ int, key string, value any) error {
		tp, defaultValue, ok := types.MaybeTimeType(key, value)
		if !ok {
			if types.IsInteger(value) {
				tp = "int"
			} else if types.IsFloat(value) {
				tp = "float64"
			} else if types.IsString(value) {
				tp = "string"
			} else if types.IsBool(value) {
				tp = "bool"
			} else if types.IsTime(value) {
				tp = "time.Time"
			} else {
				return fmt.Errorf("value must be basic type, got %T", v)
			}
			defaultValue = vars.SqlDefault[tp]
			tp = vars.SqlTypeFromGo[tp]
		}
		w.WriteWithIndentStringf("`%s` %s NOT NULL DEFAULT %s,", key, tp, defaultValue)
		w.NewLine()
		return nil
	})

	if types.MayContainsPrimary(sm) {
		w.WriteWithIndentStringln("PRIMARY KEY (`id`)")
	} else {
		w.UndoNewLine()
		w.Undo(",")
		w.NewLine()
	}
	w.WriteStringln(`) ENGINE=InnoDB DEFAULT CHARSET=utf8;`)
	return w.String(), nil
}

func ConvertToGoStruct(v any, root bool) (string, error) {
	kv, ok := v.(map[string]any)
	if !ok {
		return "", fmt.Errorf("input must be object, got %T", v)
	}
	w := writer.New(constx.Indent)
	if root {
		w.Writef("// %s\n\n", constx.HeaderText)
		w.WriteStringln("package types")
		w.NewLine()
		w.WriteStringln("type StructName struct{")
	} else {
		w.WriteStringln("struct{")
	}
	if len(kv) == 0 {
		w.UndoNewLine()
		w.Writef("}")
		return w.String(), nil
	}

	sm := sortmap.From(kv)
	sm.Range(func(_ int, key string, value any) error {
		memberType, err := convertGoStructMemberType(value)
		if err != nil {
			return err
		}
		w.Writef("%s %s `json:\"%s,omitempty\"`\n", strcase.ToCamel(key), memberType, key)
		return nil
	})
	w.Writef("}")
	if root {
		formated, err := goformat.Source(w.Bytes())
		if err != nil {
			return "", err
		}
		return string(formated), nil
	}
	return w.String(), nil
}

func convertGoStructMemberType(value any) (string, error) {
	switch {
	case types.IsInteger(value):
		return "int64", nil
	case types.IsFloat(value):
		return "float64", nil
	case types.IsBool(value):
		return "bool", nil
	case types.IsTime(value):
		return "time.Time", nil
	case types.IsString(value):
		return "string", nil
	default:
		_, ok := value.(map[string]any)
		if ok {
			return ConvertToGoStruct(value, false)
		}
		list, ok := value.([]any)
		if !ok {
			return "", fmt.Errorf("unsupport type, got %T", value)
		}
		if len(list) == 0 {
			return "[]any", nil
		}
		first := list[0]
		_, ok = first.(map[string]any)
		if ok {
			var memberSet = make(map[string]any)
			for _, v := range list {
				m, ok := v.(map[string]any)
				if !ok {
					continue
				}
				for k, v := range m {
					memberSet[k] = v
				}
			}
			tp, err := convertGoStructMemberType(memberSet)
			if err != nil {
				return "", err
			}
			return "[]" + tp, nil
		}
		tp, err := convertGoStructMemberType(first)
		if err != nil {
			return "", err
		}
		return "[]" + tp, nil
	}
}

func getIdent(c int) string {
	var list []string
	for i := 0; i < c; i++ {
		list = append(list, constx.Indent)
	}
	return strings.Join(list, "")
}

func ConvertToProtoBuf(indentCount int, key string, v any, root bool) (tp string, containsAny bool, err error) {
	name := strcase.ToCamel(key)
	kv, ok := v.(map[string]any)
	if !ok {
		return "", false, fmt.Errorf("input must be object, got %T", v)
	}
	w := writer.New(getIdent(indentCount))
	if root {
		w.Writef("// %s\n\n", constx.HeaderText)
		w.WriteStringln(`syntax = "proto3";`)
		w.NewLine()
		w.WriteStringln("package package_name;")
		w.NewLine()
		w.WriteStringln(`option go_package = ".;package_name";`)
		w.WriteStringln(constx.ProtobufTypeAny)
		w.WriteStringln(`message MessageRequest{}`)
		w.NewLine()
	}
	defer func() {
		if root && !containsAny {
			w.Remove(constx.ProtobufTypeAny)
			tp = w.String()
		}
	}()

	w.WriteWithIndentStringf("message %s{\n", name)
	if len(kv) == 0 {
		w.UndoNewLine()
		w.Writef("}")
		if root {
			w.NewLine()
			w.NewLine()
			w.Writef(`service ServiceName{
  rpc MethodName(MessageRequest) returns (MessageResponse);
}`)

		}
		return w.String(), false, nil
	}

	sm := sortmap.From(kv)
	var hasAnyType bool
	innerMessageWriter := writer.New(getIdent(indentCount + 1))
	memberWriter := writer.New(getIdent(indentCount + 1))
	sm.Range(func(idx int, key string, value any) error {
		result, err := convertProtobufMemberType(indentCount+1, key, value)
		if err != nil {
			return err
		}
		if result.ContainsAny {
			hasAnyType = true
		}
		if result.IsMessage {
			innerMessageWriter.WriteStringln(result.TypeExpr)
		}
		if result.IsArray {
			memberWriter.WriteWithIndentStringf("repeated %s %s = %d;\n", result.TypeName, strcase.ToSnake(key), idx+1)
		} else {
			memberWriter.WriteWithIndentStringf("optional %s %s = %d;\n", result.TypeName, strcase.ToSnake(key), idx+1)
		}

		return nil
	})
	w.Writef(memberWriter.String())
	w.Writef(innerMessageWriter.String())
	w.WriteWithIndentStringf("}")
	if root {
		w.NewLine()
		w.NewLine()
		w.Writef(`service ServiceName{
  rpc MethodName(MessageRequest) returns (MessageResponse);
}`)
	}
	return w.String(), hasAnyType, nil
}

type protobufMemberResult struct {
	TypeExpr    string
	TypeName    string
	IsMessage   bool
	ContainsAny bool
	IsArray     bool
}

func convertProtobufMemberType(indentCount int, key string, value any) (*protobufMemberResult, error) {
	resp := new(protobufMemberResult)
	switch {
	case types.IsInteger(value):
		resp.TypeExpr = "int64"
		resp.TypeName = "int64"
		return resp, nil
	case types.IsFloat(value):
		resp.TypeExpr = "double"
		resp.TypeName = "double"
		return resp, nil
	case types.IsBool(value):
		resp.TypeExpr = "bool"
		resp.TypeName = "bool"
		return resp, nil
	case types.IsTime(value):
		resp.TypeExpr = "string"
		resp.TypeName = "string"
		return resp, nil
	case types.IsString(value):
		resp.TypeExpr = "string"
		resp.TypeName = "string"
		return resp, nil
	default:
		_, ok := value.(map[string]any)
		if ok {
			tp, hasAnyType, err := ConvertToProtoBuf(indentCount, key, value, false)
			if err != nil {
				return nil, err
			}
			resp.TypeExpr = tp
			resp.TypeName = strcase.ToCamel(key)
			resp.IsMessage = true
			resp.ContainsAny = hasAnyType
			return resp, nil
		}
		list, ok := value.([]any)
		if !ok {
			return nil, fmt.Errorf("unsupport type, got %T", value)
		}
		if len(list) == 0 {
			resp.ContainsAny = true
			resp.TypeExpr = "google.protobuf.Any"
			resp.TypeName = "google.protobuf.Any"
			resp.IsArray = true
			return resp, nil
		}
		first := list[0]
		_, ok = first.(map[string]any)
		if ok {
			var memberSet = make(map[string]any)
			for _, v := range list {
				m, ok := v.(map[string]any)
				if !ok {
					continue
				}
				for k, v := range m {
					memberSet[k] = v
				}
			}
			tp, hasAny, err := ConvertToProtoBuf(indentCount, key, memberSet, false)
			if err != nil {
				return nil, err
			}
			resp.TypeExpr = tp
			resp.TypeName = strcase.ToCamel(key)
			resp.IsMessage = true
			resp.IsArray = true
			resp.ContainsAny = hasAny
			return resp, nil
		}
		result, err := convertProtobufMemberType(indentCount, key, first)
		if err != nil {
			return nil, err
		}
		resp.TypeExpr = result.TypeExpr
		resp.TypeName = result.TypeName
		resp.IsMessage = result.IsMessage
		resp.IsArray = true
		resp.ContainsAny = result.ContainsAny
		return resp, nil
	}
}

func ConvertToGoctlAPI(indentCount int, parent, key string, v any, root bool) (tp string, externalTypes []string, err error) {
	name := strcase.ToCamel(parent) + strcase.ToCamel(key)
	kv, ok := v.(map[string]any)
	if !ok {
		return "", nil, fmt.Errorf("input must be object, got %T", v)
	}

	w := writer.New(getIdent(indentCount))
	if root {
		w.Writef("// %s\n\n", constx.HeaderText)
		w.WriteStringln(`syntax = "v1"`)
		w.NewLine()
		w.WriteStringln(`type GeneratedRequest{}`)
		w.NewLine()
	}

	w.WriteWithIndentStringf("type %s {\n", name)
	if len(kv) == 0 {
		w.UndoNewLine()
		w.Writef("}")
		if root {
			w.NewLine()
			w.NewLine()
			w.Writef(constx.GoctlAPIService)

		}
		return w.String(), nil, nil
	}

	var externalTypeList []string
	sm := sortmap.From(kv)
	memberWriter := writer.New(getIdent(indentCount + 1))
	sm.Range(func(_ int, key string, value any) error {
		result, err := convertGoctlAPIMemberType(indentCount+1, name, key, value)
		if err != nil {
			return err
		}
		externalTypeList = append(externalTypeList, result.ExternalTypeExpr...)
		if result.IsStruct {
			externalTypeList = append(externalTypeList, result.TypeExpr)
		}
		if result.IsArray {
			memberWriter.WriteWithIndentStringf("%s []%s `json:\"%s\"`\n", strcase.ToCamel(key), result.TypeName, key)
		} else {
			memberWriter.WriteWithIndentStringf("%s %s `json:\"%s\"`\n", strcase.ToCamel(key), result.TypeName, key)
		}

		return nil
	})
	w.Writef(memberWriter.String())
	w.WriteWithIndentStringf("}")

	if root {
		w.NewLine()
		w.WriteStringln(strings.Join(externalTypeList, "\n\n"))
	}

	if root {
		w.NewLine()
		w.NewLine()
		w.Writef(constx.GoctlAPIService)
		outWriter := bytes.NewBuffer(nil)
		err = format.Source(w.Bytes(), outWriter)
		if err != nil {
			return "", nil, err
		}
		return outWriter.String(), nil, nil
	}
	return w.String(), externalTypeList, nil
}

type goctlAPIMemberResult struct {
	TypeExpr         string
	TypeName         string
	IsStruct         bool
	IsArray          bool
	ExternalTypeExpr []string
}

func convertGoctlAPIMemberType(indentCount int, parent, key string, value any) (*goctlAPIMemberResult, error) {
	resp := new(goctlAPIMemberResult)
	switch {
	case types.IsInteger(value):
		resp.TypeExpr = "int64"
		resp.TypeName = "int64"
		return resp, nil
	case types.IsFloat(value):
		resp.TypeExpr = "double"
		resp.TypeName = "double"
		return resp, nil
	case types.IsBool(value):
		resp.TypeExpr = "bool"
		resp.TypeName = "bool"
		return resp, nil
	case types.IsTime(value):
		resp.TypeExpr = "string"
		resp.TypeName = "string"
		return resp, nil
	case types.IsString(value):
		resp.TypeExpr = "string"
		resp.TypeName = "string"
		return resp, nil
	default:
		_, ok := value.(map[string]any)
		if ok {
			tp, externalTypes, err := ConvertToGoctlAPI(indentCount, parent, key, value, false)
			if err != nil {
				return nil, err
			}
			resp.TypeExpr = tp
			resp.TypeName = "*" + strcase.ToCamel(parent) + strcase.ToCamel(key)
			resp.IsStruct = true
			resp.ExternalTypeExpr = append(resp.ExternalTypeExpr, externalTypes...)
			return resp, nil
		}
		list, ok := value.([]any)
		if !ok {
			return nil, fmt.Errorf("unsupport type, got %T", value)
		}
		if len(list) == 0 {
			resp.TypeExpr = "interface{}"
			resp.TypeName = "interface{}"
			resp.IsArray = true
			return resp, nil
		}
		first := list[0]
		_, ok = first.(map[string]any)
		if ok {
			var memberSet = make(map[string]any)
			for _, v := range list {
				m, ok := v.(map[string]any)
				if !ok {
					continue
				}
				for k, v := range m {
					memberSet[k] = v
				}
			}
			tp, externalTypes, err := ConvertToGoctlAPI(indentCount, parent, key, memberSet, false)
			if err != nil {
				return nil, err
			}
			resp.TypeExpr = tp
			resp.TypeName = "*" + strcase.ToCamel(parent) + strcase.ToCamel(key)
			resp.IsStruct = true
			resp.IsArray = true
			resp.ExternalTypeExpr = append(resp.ExternalTypeExpr, externalTypes...)
			return resp, nil
		}
		result, err := convertGoctlAPIMemberType(indentCount, parent, key, first)
		if err != nil {
			return nil, err
		}
		resp.TypeExpr = result.TypeExpr
		resp.TypeName = result.TypeName
		resp.IsStruct = result.IsStruct
		resp.IsArray = true
		resp.ExternalTypeExpr = append(resp.ExternalTypeExpr, result.ExternalTypeExpr...)
		return resp, nil
	}
}

func ConvertToXml(v any) (string, error) {
	kv, ok := v.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("input must be object, got %T", v)
	}

	m := mxj.Map(kv)
	xmlObj, err := m.Xml()
	if err != nil {
		return "", err
	}

	mv, err := mxj.NewMapXmlReader(bytes.NewBuffer(xmlObj))
	if err != nil {
		return "", err
	}

	w := bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?>`)
	w.WriteByte('\n')
	err = mv.XmlIndentWriter(w, "", "  ")
	if err != nil {
		return "", err
	}
	return w.String(), nil
}
