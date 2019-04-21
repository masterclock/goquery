package goquery

import (
	"reflect"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-test/deep"
)

func TestBuilderContext_parseKeyValuePair(t *testing.T) {
	type args struct {
		op   Op
		elem interface{}
	}
	tests := []struct {
		name    string
		b       *builderContext
		args    args
		want    []sq.Sqlizer
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "$and key:value",
			b: &builderContext{
				builder:   &Builder{},
				tableName: "table1",
			},
			args: args{
				op: OpEq,
				elem: map[string]interface{}{
					"a": 1,
					"b": 2,
				},
			},
			want: []sq.Sqlizer{
				sq.Eq{`table1.a`: 1},
				sq.Eq{`table1.b`: 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.parseKeyValuePair(tt.args.op, tt.args.elem)
			if (err != nil) != tt.wantErr {
				t.Errorf("builderContext.parseKeyValuePair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("builderContext.parseKeyValuePair() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilderContext_parseElem(t *testing.T) {
	type args struct {
		elem reflect.Value
	}
	builder, _ := New(BuilderConfig{})
	tests := []struct {
		name    string
		b       *builderContext
		args    args
		want    sq.Sqlizer
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "parseElem",
			b: &builderContext{
				builder:   builder,
				tableName: "table1",
			},
			args: args{
				elem: reflect.ValueOf(map[string]interface{}{
					"a": 1,
					"b": 2,
				}),
			},
			want: sq.And{
				sq.And([]sq.Sqlizer{sq.Eq{"table1.a": 1}}),
				sq.And([]sq.Sqlizer{sq.Eq{"table1.b": 2}}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.parseElem(tt.args.elem)
			if (err != nil) != tt.wantErr {
				t.Errorf("builderContext.parseElem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !andCompare(got.(sq.And), tt.want.(sq.And)) {
				diff := deep.Equal(got, tt.want)
				t.Errorf("builderContext.parseElem() = %v, want %v, diff = %v", got, tt.want, diff)
			}
		})
	}
}

func TestBuilderContext_parseVal(t *testing.T) {
	config := BuilderConfig{}
	builder, _ := New(config)
	type args struct {
		key   string
		where interface{}
	}
	tests := []struct {
		name    string
		b       *builderContext
		args    args
		want    []sq.Sqlizer
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "trivial",
			b: &builderContext{
				builder:   builder,
				tableName: "table1",
			},
			args: args{
				key:   "a",
				where: 1,
			},
			want: []sq.Sqlizer{
				sq.Eq{"table1.a": 1},
			},
		},
		{
			name: "with op",
			b: &builderContext{
				builder:   builder,
				tableName: "table1",
			},
			args: args{
				key: "a",
				where: map[string]interface{}{
					"$lt": 1,
					"$gt": 2,
				},
			},
			want: []sq.Sqlizer{
				sq.Lt{"table1.a": 1},
				sq.Gt{"table1.a": 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.parseVal(tt.args.key, tt.args.where)
			if (err != nil) != tt.wantErr {
				t.Errorf("builderContext.parseVal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("builderContext.parseVal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilderContext_parseOp(t *testing.T) {
	builder, _ := New(BuilderConfig{})
	type args struct {
		op      Op
		operand interface{}
	}
	tests := []struct {
		name    string
		b       *builderContext
		args    args
		want    []sq.Sqlizer
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "trivial",
			b: &builderContext{
				builder:   builder,
				tableName: "table1",
			},
			args: args{
				op: OpEq,
				operand: map[string]interface{}{
					"a": 1,
					"b": 2,
				},
			},
			want: []sq.Sqlizer{
				sq.Eq{"table1.a": 1},
				sq.Eq{"table1.b": 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.parseOp(tt.args.op, tt.args.operand)
			if (err != nil) != tt.wantErr {
				t.Errorf("builderContext.parseOp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("builderContext.parseOp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_builderContext_parseWhere(t *testing.T) {
	builder, _ := New(BuilderConfig{})
	type args struct {
		where interface{}
	}
	tests := []struct {
		name    string
		b       *builderContext
		args    args
		want    sq.Sqlizer
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "trivial and",
			b: &builderContext{
				builder:   builder,
				tableName: "table1",
				rel:       OpAnd,
			},
			args: args{
				where: map[string]interface{}{
					"a": 1,
					"b": 2,
				},
			},
			want: sq.And{
				sq.And([]sq.Sqlizer{sq.Eq{"table1.a": 1}}),
				sq.And([]sq.Sqlizer{sq.Eq{"table1.b": 2}}),
			},
		},
		{
			name: "trivial or",
			b: &builderContext{
				builder:   builder,
				tableName: "table1",
				rel:       OpOr,
			},
			args: args{
				where: map[string]interface{}{
					"a": 1,
					"b": 2,
				},
			},
			want: sq.Or{
				sq.And([]sq.Sqlizer{sq.Eq{"table1.a": 1}}),
				sq.And([]sq.Sqlizer{sq.Eq{"table1.b": 2}}),
			},
		},
		{
			name: "op: {key: value} syntax",
			b: &builderContext{
				builder:   builder,
				tableName: "table1",
				rel:       OpAnd,
			},
			args: args{
				where: map[string]interface{}{
					"$eq": map[string]interface{}{
						"a": 1,
						"b": 2,
					},
					"b": 2,
				},
			},
			want: sq.And{
				sq.And([]sq.Sqlizer{sq.Eq{"table1.a": 1}, sq.Eq{"table1.b": 2}}),
				sq.And([]sq.Sqlizer{sq.Eq{"table1.b": 2}}),
			},
		},
		{
			name: "key: {op: value} syntax",
			b: &builderContext{
				builder:   builder,
				tableName: "table1",
				rel:       OpAnd,
			},
			args: args{
				where: map[string]interface{}{
					"a": map[string]interface{}{
						"$eq": 1,
						"$gt": 2,
					},
					"b": 2,
				},
			},
			want: sq.And{
				sq.And([]sq.Sqlizer{sq.Eq{"table1.a": 1}, sq.Gt{"table1.a": 2}}),
				sq.And([]sq.Sqlizer{sq.Eq{"table1.b": 2}}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.parseWhere(tt.args.where)
			if (err != nil) != tt.wantErr {
				t.Errorf("builderContext.parseWhere() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !sqCompare(got, tt.want) {
				t.Errorf("builderContext.parseWhere() = %v, want %v", got, tt.want)
			}
		})
	}
}
