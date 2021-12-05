package str

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
	"unicode/utf8"
)

func ToStrStruct(value interface{}) []string {
	rf := reflect.TypeOf(value)
	rv := reflect.ValueOf(value)

	num := rf.NumField()
	row := make([]string, num)
	for i := 0; i < num; i++ {
		row[i] = ToStr(rv.Field(i))
	}
	return row
}

func ToStr(v interface{}) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case []byte:
		if ok := utf8.Valid(t); ok {
			return string(t)
		}
	case int:
		return strconv.Itoa(t)
	case int32:
		return strconv.FormatInt(int64(t), 10)
	case int64:
		return strconv.FormatInt(t, 10)
	case time.Time:
		return t.Format(time.RFC3339)
	}
	return fmt.Sprint(v)
}
