package goquery

import (
	"errors"
	"fmt"
	"reflect"

	sq "github.com/Masterminds/squirrel"
)

type Filter struct {
	From       string
	Where      map[string]interface{}
	Attributes []string
	Include    []map[string]interface{}
	Order      []map[string]interface{}
	Offset     uint
	Limit      uint
}

type Op string

const (
	OpAnd   = "$and"
	OpOr    = "$or"
	OpNot   = "$not"
	OpEq    = "$eq"
	OpNotEq = "$notEq"
	OpGt    = "$gt"
	OpGte   = "$gte"
	OpLt    = "$lt"
	OpLte   = "$lte"
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

type Builder struct {
	// sqBuilder    sq.SelectBuilder
	operators    map[Op]string
	revOperators map[string]Op
}

type BuilderConfig struct {
	operatorMapping map[string]string
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
	}
	return builder, nil
}

// Build build new query from filter
func (b *Builder) Build(filter Filter) (sq.Sqlizer, error) {
	wheres, err := b.parseWhere(filter.Where)
	if err != nil {
		return nil, err
	}
	bs := sq.Select("*").From(filter.From)
	for _, w := range wheres {
		bs = bs.Where(w)
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

func (b *Builder) toOperator(str string) (Op, error) {
	var op Op
	var ok bool
	if op, ok = b.revOperators[str]; !ok {
		return "", fmt.Errorf("not operator: %s", str)
	}
	return op, nil
}

func (b *Builder) parseWhere(where map[string]interface{}) ([]sq.Sqlizer, error) {
	conds := sq.And{}
	for k, v := range where {
		rv := reflect.ValueOf(v)
		t := rv.Kind()
		switch t {
		case reflect.Map:
			{
				var cond []sq.Sqlizer
				if op, err := b.toOperator(k); err == nil {
					cond, err = b.parseOp(op, v)
					if err != nil {
						return nil, err
					}
				} else {
					cond, err = b.parseVal(k, v)
					if err != nil {
						return nil, err
					}
				}
				conds = append(conds, cond...)
				break
			}
		case reflect.Struct:
			{
				break
			}
		default:
			{
				cond := sq.Eq{k: v}
				conds = append(conds, cond)
			}
		}
	}
	return conds, nil
}

func (b *Builder) parseOp(op Op, operand interface{}) ([]sq.Sqlizer, error) {
	switch op {
	case OpAnd:
		{
			conds, err := b.parseMultiple(operand)
			if err != nil {
				return nil, err
			}
			return sq.And(conds), nil
		}
	case OpOr:
		{
			conds, err := b.parseMultiple(operand)
			if err != nil {
				return nil, err
			}
			return sq.Or(conds), nil
		}
	case OpNot:
		{
			return nil, errors.New("unimplementd")

		}
	default:
		{
			conds, err := b.parseKeyValuePair(op, operand)
			if err != nil {
				return nil, err
			}
			return sq.And(conds), nil
		}
	}
}

func (b *Builder) parseMultiple(operand interface{}) ([]sq.Sqlizer, error) {
	rv := reflect.ValueOf(operand)
	kind := rv.Kind()
	switch kind {
	case reflect.Array:
		{
			conds := []sq.Sqlizer{}
			length := rv.Len()
			for i := 0; i < length; i++ {
				elem := rv.Index(i)
				cond, err := b.parseElem(elem)
				if err != nil {
					return nil, err
				}
				conds = append(conds, cond...)
			}
			return conds, nil
		}
	case reflect.Slice:
		{
			conds := []sq.Sqlizer{}
			length := rv.Len()
			for i := 0; i < length; i++ {
				elem := rv.Index(i)
				cond, err := b.parseElem(elem)
				if err != nil {
					return nil, err
				}
				conds = append(conds, cond...)
			}
			return conds, nil
		}
	default:
		{
			return nil, errors.New("invalid operand")
		}
	}
}

func (b *Builder) parseElem(elem reflect.Value) ([]sq.Sqlizer, error) {
	kind := elem.Kind()
	switch kind {
	case reflect.Map:
		{
			m := map[string]interface{}{}
			iter := elem.MapRange()
			for iter.Next() {
				key := iter.Key()
				value := iter.Value()
				keyKind := key.Kind()
				if keyKind != reflect.String {
					return nil, errors.New("key must be string")
				}
				m[key.String()] = value.Interface()
			}
			return b.parseWhere(m)
		}
	default:
		{
			return nil, errors.New("unimplemented")
		}
	}

}

func (b *Builder) parseKeyValuePair(op Op, elem interface{}) ([]sq.Sqlizer, error) {
	rv := reflect.ValueOf(elem)
	kind := rv.Kind()
	switch kind {
	case reflect.Map:
		{
			m := map[string]interface{}{}
			iter := rv.MapRange()
			for iter.Next() {
				key := iter.Key()
				value := iter.Value()
				keyKind := key.Kind()
				if keyKind != reflect.String {
					return nil, errors.New("key must be string")
				}
				m[key.String()] = value.Interface()
			}
			return wrapOp(op, m)
		}
	default:
		{
			return nil, errors.New("invalid operand")
		}
	}
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

func (b *Builder) parseVal(key string, where interface{}) ([]sq.Sqlizer, error) {
	rv := reflect.ValueOf(where)
	switch rv.Kind() {
	case reflect.Map:
		{
			conds := []sq.Sqlizer{}
			iter := rv.MapRange()
			for iter.Next() {
				k := iter.Key()
				switch k.Kind() {
				case reflect.String:
					{
						op, err := b.toOperator(k.String())
						if err != nil {
							return nil, err
						}
						operand := map[string]interface{}{
							key: iter.Value().Interface(),
						}
						parts, err := b.parseOp(op, operand)
						if err != nil {
							return nil, err
						}
						conds = append(conds, parts...)
					}
				default:
					{
						return nil, errors.New("invalid key type")
					}
				}
			}
			return conds, nil
		}
	default:
		{
			return nil, errors.New("invalid operand")
		}
	}
}
