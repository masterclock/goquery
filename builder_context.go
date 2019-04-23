package goquery

import (
	"errors"
	"fmt"
	"reflect"

	sq "github.com/Masterminds/squirrel"
)

type builderContext struct {
	builder   *Builder
	rel       Op
	tableName string
}

func (b *builderContext) parseWhere(where interface{}) (sq.Sqlizer, error) {
	rv := reflect.ValueOf(where)
	switch rv.Kind() {
	case reflect.Map:
		{
			conds := []sq.Sqlizer{}
			iter := rv.MapRange()
			for iter.Next() {
				key := iter.Key()
				if key.Kind() != reflect.String {
					return nil, errors.New("invalid syntax, key must be string")
				}
				cond, err := b.parseWhereEntry(key.String(), iter.Value().Interface())
				if err != nil {
					return nil, err
				}
				conds = append(conds, cond)
			}
			switch b.rel {
			case OpAnd:
				{
					return sq.And(conds), nil
				}
			case OpOr:
				{
					return sq.Or(conds), nil
				}
			default:
				{
					return nil, errors.New("invalid syntax, expect relation op")
				}
			}
		}
	default:
		{
			return nil, errors.New("invalid syntax")
		}
	}
}

func (b *builderContext) parseWhereEntry(key string, value interface{}) (sq.Sqlizer, error) {
	var cond []sq.Sqlizer
	if op, err := b.toOperator(key); err == nil {
		cond, err = b.parseOp(op, value)
		if err != nil {
			return nil, err
		}
	} else {
		cond, err = b.parseVal(key, value)
		if err != nil {
			return nil, err
		}
	}
	// if len(cond) == 1 {
	// 	return cond[0], nil
	// }
	return sq.And(cond), nil
}

func (b *builderContext) parseOp(op Op, operand interface{}) ([]sq.Sqlizer, error) {
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

func (b *builderContext) parseMultiple(operand interface{}) ([]sq.Sqlizer, error) {
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
				conds = append(conds, cond)
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
				conds = append(conds, cond)
			}
			return conds, nil
		}
	default:
		{
			return nil, errors.New("invalid operand")
		}
	}
}

func (b *builderContext) parseElem(elem reflect.Value) (sq.Sqlizer, error) {
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
			nb := b.inherit()
			nb.rel = OpAnd
			fmt.Printf("%+v, %+v\n", nb, *b)
			return nb.parseWhere(m)
		}
	default:
		{
			return nil, errors.New("unimplemented")
		}
	}

}

func (b *builderContext) parseKeyValuePair(op Op, elem interface{}) ([]sq.Sqlizer, error) {
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
				attr := key.String()
				alias := b.toFullName(b.tableName, attr)
				m[alias] = value.Interface()
			}
			return wrapOp(op, m)
		}
	default:
		{
			return nil, errors.New("invalid operand")
		}
	}
}

func (b *builderContext) parseVal(key string, where interface{}) ([]sq.Sqlizer, error) {
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
			alias := b.toFullName(b.tableName, key)
			return []sq.Sqlizer{sq.Eq{alias: where}}, nil
		}
	}
}

func (b *builderContext) quote(name string) string {
	return fmt.Sprintf("%s%s%s", b.builder.config.Quote, name, b.builder.config.Quote)
}

func (b *builderContext) toOperator(str string) (Op, error) {
	var op Op
	var ok bool
	if op, ok = b.builder.revOperators[str]; !ok {
		return "", fmt.Errorf("not operator: %s", str)
	}
	return op, nil
}

func (b *builderContext) inherit() builderContext {
	return *b
}
