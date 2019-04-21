package goquery

import (
	"reflect"

	sq "github.com/Masterminds/squirrel"
)

func setCompare(a []sq.Sqlizer, b []sq.Sqlizer) bool {
	if len(a) != len(b) {
		return false
	}
	for _, va := range a {
		if !contains(b, va) {
			return false
		}
	}
	for _, vb := range b {
		if !contains(a, vb) {
			return false
		}
	}
	return true
}

func andCompare(a sq.And, b sq.And) bool {
	return setCompare([]sq.Sqlizer(a), []sq.Sqlizer(b))
}

func orCompare(a sq.Or, b sq.Or) bool {
	return setCompare([]sq.Sqlizer(a), []sq.Sqlizer(b))
}

func contains(lst []sq.Sqlizer, v sq.Sqlizer) bool {
	for _, item := range lst {
		if sqCompare(item, v) {
			return true
		}
	}
	return false
}

func sqCompare(a sq.Sqlizer, b sq.Sqlizer) bool {
	switch ta := a.(type) {
	case sq.And:
		{
			return sqCompareAnd2Other(ta, b)
		}
	case sq.Or:
		{
			return sqCompareOr2Other(ta, b)
		}
	default:
		return reflect.DeepEqual(a, b)
	}
}

func sqCompareAnd2Other(ta sq.And, b sq.Sqlizer) bool {
	switch tb := b.(type) {
	case sq.And:
		{
			return andCompare(ta, tb)
		}
	default:
		return false
	}
}

func sqCompareOr2Other(ta sq.Or, b sq.Sqlizer) bool {
	switch tb := b.(type) {
	case sq.Or:
		{
			return orCompare(ta, tb)
		}
	default:
		return false
	}
}
