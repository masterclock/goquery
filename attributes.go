package goquery

import (
	"errors"
	"fmt"
	"reflect"
)

func (b *builderContext) attrBuild(attributes []interface{}, tableName string) ([]string, error) {
	attrs := []string{}
	for _, v := range attributes {
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.String:
			{
				str := rv.String()
				colName := b.toFullName(tableName, str)
				aliased := fmt.Sprintf("%s AS %s", colName, b.quote(str))
				attrs = append(attrs, aliased)
				break
			}
		default:
			{
				return nil, errors.New("invalid syntax")
			}
		}
	}
	if len(attrs) == 0 {
		attrs = []string{"*"}
	}
	return attrs, nil
}
