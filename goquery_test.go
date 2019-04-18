package goquery

import (
	"fmt"
	"reflect"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-test/deep"
)

func TestBuilder_parseWhere(t *testing.T) {
	type args struct {
		where map[string]interface{}
	}
	cfg := BuilderConfig{}
	builder, _ := New(cfg)
	tests := []struct {
		name    string
		b       *Builder
		args    args
		want    []sq.Sqlizer
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "simple value equal",
			b:    builder,
			args: args{
				where: map[string]interface{}{
					"a": "str_a",
					"b": "str_b",
					"c": 5,
					"d": true,
					"e": 1.0,
				},
			},
			want: sq.And{
				sq.Eq{"a": "str_a"},
				sq.Eq{"b": "str_b"},
				sq.Eq{"c": 5},
				sq.Eq{"d": true},
				sq.Eq{"e": 1.0},
			},
			wantErr: false,
		},
		{
			name: "op: {key: value} syntax",
			b:    builder,
			args: args{
				where: map[string]interface{}{
					OpGt: map[string]interface{}{
						"a": 1,
						"b": 2,
					},
				},
			},
			want: sq.And{
				sq.Gt{"a": 1},
				sq.Gt{"b": 2},
			},
		},
		{
			name: "key: {op: value} syntax",
			b:    builder,
			args: args{
				where: map[string]interface{}{
					"a": map[Op]interface{}{
						OpGt: 1,
						OpLt: 2,
					},
				},
			},
			want: sq.And{
				sq.Gt{"a": 1},
				sq.Lt{"a": 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.parseWhere(tt.args.where)
			if tt.wantErr && err == nil {
				t.Errorf("Builder.parseWhere = %v, err = nil, wantErr = %v", got, tt.want)
				return
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Builder.parseWhere = %v, err = %v, want = %v", got, err, tt.want)
				return
			}
			if !setCompare(got, tt.want) {
				t.Errorf("Builder.parseWhere() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_parseElem(t *testing.T) {
	type args struct {
		elem reflect.Value
	}
	tests := []struct {
		name    string
		b       *Builder
		args    args
		want    []sq.Sqlizer
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "parseElem",
			b:    &Builder{},
			args: args{
				elem: reflect.ValueOf(map[string]interface{}{
					"a": "a_val",
					"b": 2,
				}),
			},
			want: []sq.Sqlizer{
				sq.Eq{"a": "a_val"},
				sq.Eq{"b": 2},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.parseElem(tt.args.elem)
			if (err != nil) != tt.wantErr {
				t.Errorf("Builder.parseElem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.parseElem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_parseOp(t *testing.T) {
	type args struct {
		op      Op
		operand interface{}
	}
	tests := []struct {
		name    string
		b       *Builder
		args    args
		want    []sq.Sqlizer
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "parse $and",
			b:    &Builder{},
			args: args{
				op: OpAnd,
				operand: []map[string]interface{}{
					{
						"a": 1,
						"b": 2,
					},
				},
			},
			want: sq.And{
				sq.Eq{"a": 1},
				sq.Eq{"b": 2},
			},
			wantErr: false,
		},
		{
			name: "parse $or",
			b:    &Builder{},
			args: args{
				op: OpOr,
				operand: []map[string]interface{}{
					{
						"a": 1,
						"b": 2,
					},
				},
			},
			want: sq.Or{
				sq.Eq{"a": 1},
				sq.Eq{"b": 2},
			},
			wantErr: false,
		},
		{
			name: "parse $eq",
			b:    &Builder{},
			args: args{
				op: OpEq,
				operand: map[string]interface{}{
					"a": 1,
					"b": 2,
				},
			},
			want: sq.And{
				sq.Eq{"a": 1},
				sq.Eq{"b": 2},
			},
			wantErr: false,
		},
		{
			name: "parse $notEq",
			b:    &Builder{},
			args: args{
				op: OpNotEq,
				operand: map[string]interface{}{
					"a": 1,
					"b": 2,
				},
			},
			want: sq.And{
				sq.NotEq{"a": 1},
				sq.NotEq{"b": 2},
			},
			wantErr: false,
		},
		{
			name: "parse $gt",
			b:    &Builder{},
			args: args{
				op: OpGt,
				operand: map[string]interface{}{
					"a": 1,
					"b": 2,
				},
			},
			want: sq.And{
				sq.Gt{"a": 1},
				sq.Gt{"b": 2},
			},
			wantErr: false,
		},
		{
			name: "parse $gte",
			b:    &Builder{},
			args: args{
				op: OpGte,
				operand: map[string]interface{}{
					"a": 1,
					"b": 2,
				},
			},
			want: sq.And{
				sq.GtOrEq{"a": 1},
				sq.GtOrEq{"b": 2},
			},
			wantErr: false,
		},
		{
			name: "parse $lt",
			b:    &Builder{},
			args: args{
				op: OpLt,
				operand: map[string]interface{}{
					"a": 1,
					"b": 2,
				},
			},
			want: sq.And{
				sq.Lt{"a": 1},
				sq.Lt{"b": 2},
			},
			wantErr: false,
		},
		{
			name: "parse $lte",
			b:    &Builder{},
			args: args{
				op: OpLte,
				operand: map[string]interface{}{
					"a": 1,
					"b": 2,
				},
			},
			want: sq.And{
				sq.LtOrEq{"a": 1},
				sq.LtOrEq{"b": 2},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.parseOp(tt.args.op, tt.args.operand)
			if (err != nil) != tt.wantErr {
				t.Errorf("Builder.parseOp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !setCompare(got, tt.want) {
				diff := deep.Equal(got, tt.want)
				t.Errorf("Builder.parseOp() = %v, want %v, diff = %v", got, tt.want, diff)
			}
		})
	}
}

func Test_wrapOp(t *testing.T) {
	type args struct {
		op Op
		m  map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []sq.Sqlizer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := wrapOp(tt.args.op, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("wrapOp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("wrapOp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Build(t *testing.T) {
	type args struct {
		filter Filter
	}
	cfg := BuilderConfig{}
	builder, _ := New(cfg)
	tests := []struct {
		name    string
		b       *Builder
		args    args
		want    sq.Sqlizer
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "trival",
			b:    builder,
			args: args{
				filter: Filter{
					From: "table1",
					Where: map[string]interface{}{
						"a": 1,
						"b": 2,
						"$gt": map[string]interface{}{
							"a_1": 10,
							"b_1": 11,
						},
					},
				},
			},
			want:    sq.Eq{"a": 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.Build(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Builder.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			sql, params, err := got.ToSql()
			fmt.Printf("sql = %s\nparams = %v\nerr = %v\n", sql, params, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}
