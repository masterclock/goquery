package goquery

import (
	"fmt"
)

func (b *builderContext) toFullName(tableName, colName string) string {
	return fmt.Sprintf("%s.%s", b.quote(tableName), b.quote(colName))
}
