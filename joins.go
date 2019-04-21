package goquery

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func (b *builderContext) addJoins(bs sq.SelectBuilder, tableName string, includes []Include) (sq.SelectBuilder, error) {
	for _, include := range includes {
		if len(include.Where) == 0 {
			if include.Through == nil {
				src := b.toFullName(tableName, include.SourceKey)
				dst := b.toFullName(include.Table, include.ForeignKey)
				clause := fmt.Sprintf("%s ON %s = %s", include.Table, src, dst)
				bs = bs.LeftJoin(clause)
			} else {
				thSrc := b.toFullName(tableName, include.SourceKey)
				thDst := b.toFullName(include.Through.TableName, include.Through.SourceKey)
				throughClause := fmt.Sprintf("%s ON %s = %s", include.Through.TableName, thSrc, thDst)

				src := b.toFullName(include.Through.TableName, include.Through.ForeignKey)
				dst := b.toFullName(include.Table, include.ForeignKey)
				clause := fmt.Sprintf("%s ON %s = %s", include.Table, src, dst)
				bs = bs.LeftJoin(throughClause).LeftJoin(clause)
			}
		} else {
			if include.Through == nil {
				src := b.toFullName(tableName, include.SourceKey)
				dst := b.toFullName(include.Table, include.ForeignKey)
				clause := fmt.Sprintf("LEFT INNER JOIN %s ON %s = %s", include.Table, src, dst)
				bs = bs.JoinClause(clause)
				nb := b.inherit()
				nb.tableName = include.Table
				where, err := nb.parseWhere(include.Where)
				if err != nil {
					return bs, err
				}
				bs = bs.Where(where)
			}
		}
	}
	return bs, nil
}
