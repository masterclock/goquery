package goquery

import "errors"

func (b *builderContext) buildFrom(from interface{}) (string, string, error) {
	switch t := from.(type) {
	case string:
		{
			return b.quote(t), t, nil
		}
	default:
		{
			return "", "", errors.New("invalid syntax from")
		}
	}
}
