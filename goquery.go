package goquery

import (
	"errors"

	sq "github.com/Masterminds/squirrel"
)

// IncludeThrough include with a joint table
type IncludeThrough struct {
	TableName  string
	SourceKey  string
	ForeignKey string
}

// Include join definition
type Include struct {
	Table      string
	SourceKey  string
	ForeignKey string
	Through    *IncludeThrough
	Where      map[string]interface{}
}

// Filter filter structure
type Filter struct {
	From       string
	Where      map[string]interface{}
	Attributes []interface{}
	Include    []Include
	Order      []string
	Offset     *uint64
	Limit      *uint64
}

// Op operator type
type Op string

const (
	// OpAnd `and` relation
	OpAnd = "$and"
	// OpOr `or` relation
	OpOr = "$or"
	// OpNot negate
	OpNot = "$not"
	// OpEq equal
	OpEq = "$eq"
	// OpNotEq not equal
	OpNotEq = "$notEq"
	// OpGt greater than
	OpGt = "$gt"
	// OpGte greater than or equal
	OpGte = "$gte"
	// OpLt less than
	OpLt = "$lt"
	// OpLte less than or equal
	OpLte = "$lte"
	// OpLike like
	OpLike = "$like"
)

var defaultOpMapping = map[Op]string{
	OpAnd:   "$and",
	OpOr:    "$or",
	OpNot:   "$not",
	OpEq:    "$eq",
	OpNotEq: "$notEq",
	OpGt:    "$gt",
	OpGte:   "$gte",
	OpLt:    "$lt",
	OpLte:   "$lte",
}

// Builder goquery builder struct
type Builder struct {
	// sqBuilder    sq.SelectBuilder
	operators    map[Op]string
	revOperators map[string]Op
	config       BuilderConfig
}

// BuilderConfig goquery builder config
type BuilderConfig struct {
	operatorMapping map[string]string
	quote           string
}

// New create new builder
func New(config BuilderConfig) (*Builder, error) {
	ops := map[Op]string{}
	for k, v := range defaultOpMapping {
		ops[k] = v
	}
	for k, v := range config.operatorMapping {
		ops[Op(k)] = v
	}
	revOps := map[string]Op{}
	for k, v := range ops {
		revOps[v] = k
	}
	builder := &Builder{
		operators:    ops,
		revOperators: revOps,
		config:       config,
	}
	return builder, nil
}

// Build build new query from filter
func (b *Builder) Build(filter Filter) (sq.Sqlizer, error) {
	ctx := builderContext{
		builder: b,
		rel:     OpAnd,
	}
	// build main table and alias
	from, tableAlias, err := ctx.buildFrom(filter.From)
	if err != nil {
		return nil, err
	}
	ctx.tableName = tableAlias
	// build fully qualified attributes
	attributes, err := ctx.attrBuild(filter.Attributes, tableAlias)
	if err != nil {
		return nil, err
	}
	bs := sq.Select(attributes...).From(from)

	// build includes
	if len(filter.Include) > 0 {
		bs, err = ctx.addJoins(bs, tableAlias, filter.Include)
		if err != nil {
			return nil, err
		}
	}

	// build wheres
	wheres, err := ctx.parseWhere(filter.Where)
	if err != nil {
		return nil, err
	}
	bs = bs.Where(wheres)

	// add order
	if len(filter.Order) > 0 {
		bs = bs.OrderBy(filter.Order...)
	}

	// add limit
	if filter.Limit != nil && *(filter.Limit) != 0 {
		bs = bs.Limit(*(filter.Limit))
	}

	// add offset
	if filter.Offset != nil && *(filter.Offset) != 0 {
		bs = bs.Offset(*(filter.Limit))
	}
	return bs, nil
}

func (b *Builder) isOperator(op string) bool {
	for _, v := range b.operators {
		if v == op {
			return true
		}
	}
	return false
}

func wrapOp(op Op, m map[string]interface{}) ([]sq.Sqlizer, error) {
	switch op {
	case OpEq:
		{
			conds := []sq.Sqlizer{}
			for k, v := range m {
				cond := sq.Eq{k: v}
				conds = append(conds, cond)
			}
			return conds, nil
		}
	case OpNotEq:
		{
			conds := []sq.Sqlizer{}
			for k, v := range m {
				cond := sq.NotEq{k: v}
				conds = append(conds, cond)
			}
			return conds, nil
		}
	case OpGt:
		{
			conds := []sq.Sqlizer{}
			for k, v := range m {
				cond := sq.Gt{k: v}
				conds = append(conds, cond)
			}
			return conds, nil
		}
	case OpGte:
		{
			conds := []sq.Sqlizer{}
			for k, v := range m {
				cond := sq.GtOrEq{k: v}
				conds = append(conds, cond)
			}
			return conds, nil
		}
	case OpLt:
		{
			conds := []sq.Sqlizer{}
			for k, v := range m {
				cond := sq.Lt{k: v}
				conds = append(conds, cond)
			}
			return conds, nil
		}
	case OpLte:
		{
			conds := []sq.Sqlizer{}
			for k, v := range m {
				cond := sq.LtOrEq{k: v}
				conds = append(conds, cond)
			}
			return conds, nil
		}
	default:
		{
			return nil, errors.New("invalid op")
		}
	}
}
