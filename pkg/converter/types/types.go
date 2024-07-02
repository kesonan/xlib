package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kesonan/xlib/pkg/converter/constx"
	"github.com/kesonan/xlib/pkg/converter/vars"
	"github.com/kesonan/xlib/pkg/list"
	"github.com/kesonan/xlib/pkg/sortmap"
	"github.com/kesonan/xlib/pkg/stringx"
)

func MayContainsPrimary(m *sortmap.SortMap) bool {
	v, ok := m.Get(constx.MayIdColumn)
	if !ok {
		return false
	}
	return IsNumeric(v)
}

func MaybeTimeType(key string, v any) (tp string, defaultValue string, isTime bool) {
	key = strings.ReplaceAll(key, "_", "")
	key = strings.ToUpper(key)
	if !stringx.ContainsAny(key, vars.MayTimeColumn...) {
		return "", "", false
	}
	if MaybeTimestamp(v) {
		return "bigint", "0", true
	}
	if MaybeDateTime(v) {
		return "timestamp", "CURRENT_TIMESTAMP", true
	}
	if MaybeDate(v) {
		return "date", "'1970-01-01'", true
	}
	if MaybeTime(v) {
		return "time", "'00:00:01'", true
	}
	if MaybeYear(v) {
		return "year", "'1970'", true
	}
	return "", "", false
}

func MaybeYear(v any) bool {
	val, ok := v.(string)
	if !ok {
		return false
	}
	_, err := time.Parse("2006", val)
	if err != nil {
		return false
	}
	return true
}

func MaybeTime(v any) bool {
	val, ok := v.(string)
	if !ok {
		return false
	}
	_, err := time.Parse("15:04:05", val)
	if err != nil {
		return false
	}
	return true
}

func MaybeDate(v any) bool {
	val, ok := v.(string)
	if !ok {
		return false
	}
	_, err := time.Parse("2006-01-02", val)
	if err != nil {
		return false
	}
	return true
}

func MaybeDateTime(v any) bool {
	val, ok := v.(string)
	if !ok {
		return false
	}
	_, err := time.Parse("2006-01-02 15:04:05", val)
	if err != nil {
		return false
	}
	return true
}

func MaybeTimestamp(v any) bool {
	if !IsInteger(v) {
		return false
	}
	valInt, err := json.Number(fmt.Sprint(v)).Int64()
	if err != nil {
		return false
	}

	if valInt > 0 && valInt < constx.TimeStampEnd {
		return true
	}
	return false
}

func IsBasic(v any) bool {
	switch v.(type) {
	case uint8, uint16, uint32, uint64, int8, int16, int32, int64, int, uint, float32, float64, string, bool:
		return true
	default:
		return false
	}
}

func IsInteger(v any) bool {
	switch v.(type) {
	case uint8, uint16, uint32, uint64, int8, int16, int32, int64, int, uint:
		return true
	case float32, float64:
		s := fmt.Sprintf("%0.1f", v)
		return strings.HasSuffix(s, ".0")
	default:
		return false
	}
}

func IsBool(v any) bool {
	_, ok := v.(bool)
	return ok
}

func IsString(v any) bool {
	_, ok := v.(string)
	return ok
}

func IsFloat(v any) bool {
	switch v.(type) {
	case float32, float64:
		return strings.Contains(fmt.Sprint(v), ".")
	default:
		return false
	}
}

func IsNumeric(v any) bool {
	switch v.(type) {
	case uint8, uint16, uint32, uint64, int8, int16, int32, int64, int, uint, float32, float64:
		return true
	default:
		return false
	}
}

func IsValidFromType(tp string) bool {
	tp = strings.ToUpper(tp)
	_, ok := vars.ValidMapper[tp]
	return ok
}

func IsValidToType(fromType, toType string) bool {
	tp := strings.ToUpper(fromType)
	l, ok := vars.ValidMapper[tp]
	if !ok {
		return false
	}
	return list.Contains[string](l, strings.ToUpper(toType))
}
