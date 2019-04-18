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

func contains(lst []sq.Sqlizer, v sq.Sqlizer) bool {
	for _, item := range lst {
		if reflect.DeepEqual(item, v) {
			return true
		}
	}
	return false
}
