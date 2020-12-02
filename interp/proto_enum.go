package interp

import (
	"fmt"

	protodesc "github.com/jhump/protoreflect/desc"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

// ProtoEnumType is a starlark.Value that repsents a protobuf enum
type ProtoEnumType struct {
	name   string
	values map[string]int32
}

var _ protoType = (*ProtoEnumType)(nil)
var _ starlark.HasAttrs = (*ProtoEnumType)(nil)

// TODO implement casting via callable
// var _ starlark.Callable = (*ProtoEnumType)(nil)

func NewProtoEnumType(desc *protodesc.EnumDescriptor) *ProtoEnumType {
	t := &ProtoEnumType{
		name:   desc.GetName(),
		values: make(map[string]int32),
	}

	for _, valdesc := range desc.GetValues() {
		t.values[valdesc.GetName()] = valdesc.GetNumber()
	}
	return t
}

func (ProtoEnumType) protoType()              {}
func (ProtoEnumType) Type() string            { return "ProtoEnum" }
func (ProtoEnumType) Truth() starlark.Bool    { return starlark.True }
func (ProtoEnumType) Freeze()                 {}
func (t ProtoEnumType) Hash() (uint32, error) { return 0, errUnhashable }
func (t ProtoEnumType) String() string        { return fmt.Sprintf("<proto.Enum %q>", t.Name()) }
func (t ProtoEnumType) Name() string          { return t.name }

func (t ProtoEnumType) Attr(name string) (starlark.Value, error) {
	if v, ok := t.values[name]; ok {
		return &ProtoEnumValue{
			typeName: t.name,
			name:     name,
			value:    v,
		}, nil
	}
	return nil, nil
}

func (t ProtoEnumType) AttrNames() []string { return nil }

type ProtoEnumValue struct {
	typeName string
	name     string
	value    int32
}

var _ starlark.Comparable = (*ProtoEnumValue)(nil)

func (v ProtoEnumValue) Type() string         { return v.typeName }
func (v ProtoEnumValue) Truth() starlark.Bool { return starlark.True }
func (v ProtoEnumValue) Freeze()              {}

func (v ProtoEnumValue) String() string {
	return fmt.Sprintf("<%s.%s %d>", v.typeName, v.name, v.value)
}

func (v *ProtoEnumValue) Hash() (uint32, error) {
	return starlark.MakeInt64(int64(v.value)).Hash()
}

func (v *ProtoEnumValue) CompareSameType(op syntax.Token, y starlark.Value, depth int) (bool, error) {
	other := y.(*ProtoEnumValue)
	switch op {
	case syntax.EQL:
		return v.value == other.value, nil
	case syntax.NEQ:
		return v.value != other.value, nil
	default:
		return false, fmt.Errorf("enums support only `==' and `!=' comparisons, got: %#v", op)
	}
}
