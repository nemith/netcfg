package interp

import (
	"errors"
	"fmt"

	protodesc "github.com/jhump/protoreflect/desc"
	protodyn "github.com/jhump/protoreflect/dynamic"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

// interface to define protobuf type Value types (enum, messages)
type protoType interface {
	starlark.Value
	protoType()
}

var (
	errUnhashable        = errors.New("unhashable type")
	errUnexpectedPosArgs = errors.New("unexpected positional arguments")
)

// ProtoMessageType is a starlark.Value that represents a protobuf message type
//  (not an instance of that type).
type ProtoMessageType struct {
	desc        *protodesc.MessageDescriptor
	messageType string
}

var _ protoType = (*ProtoMessageType)(nil)
var _ starlark.HasAttrs = (*ProtoMessageType)(nil)
var _ starlark.Callable = (*ProtoMessageType)(nil)

func NewProtoMessageType(desc *protodesc.MessageDescriptor, parent string) *ProtoMessageType {
	return &ProtoMessageType{
		desc:        desc,
		messageType: parent + "." + desc.GetName(),
	}
}

func (ProtoMessageType) protoType()            {}
func (ProtoMessageType) Type() string          { return "proto.MessageType" }
func (ProtoMessageType) Freeze()               {}
func (ProtoMessageType) Hash() (uint32, error) { return 0, errUnhashable }
func (ProtoMessageType) Truth() starlark.Bool  { return starlark.True }
func (mt *ProtoMessageType) String() string    { return fmt.Sprintf("<ProtoMessage %q>", mt.Name()) }
func (mt *ProtoMessageType) Name() string      { return mt.desc.GetName() }

func (mt *ProtoMessageType) Attr(name string) (starlark.Value, error) {
	return nil, nil
}

func (mt *ProtoMessageType) AttrNames() []string {
	return []string{}
}

func (mt *ProtoMessageType) CallInternal(
	thread *starlark.Thread,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {
	if len(args) > 0 {
		return nil, errUnexpectedPosArgs
	}

	msg := protodyn.NewMessage(mt.desc)

	for _, kwarg := range kwargs {
		field := string(kwarg[0].(starlark.String))
		val := kwarg[1]
		if err := msg.TrySetFieldByName(field, val); err != nil {
			return nil, err
		}
	}

	return &ProtoMessage{typ: mt.messageType, msg: msg}, nil
}

// ProtoMessage is a starlark.Value that represents an instance of a protobuf message.
type ProtoMessage struct {
	msg *protodyn.Message
	typ string
}

var _ starlark.HasAttrs = (*ProtoMessage)(nil)
var _ starlark.HasSetField = (*ProtoMessage)(nil)
var _ starlark.Comparable = (*ProtoMessage)(nil)

func (m ProtoMessage) Type() string       { return m.typ }
func (ProtoMessage) Truth() starlark.Bool { return starlark.True }

func (m ProtoMessage) String() string { return "" }

func (m *ProtoMessage) Freeze() {
	//	for _, f := range m.fields {
	//		f.value.Freeze()
	//	}
}

func (m ProtoMessage) Hash() (uint32, error) { return 0, nil }

func (m ProtoMessage) Attr(name string) (starlark.Value, error)        { return nil, nil }
func (m ProtoMessage) AttrNames() []string                             { return nil }
func (m *ProtoMessage) SetField(name string, val starlark.Value) error { return nil }

func (enum ProtoMessage) CompareSameType(op syntax.Token, y starlark.Value, depth int) (bool, error) {
	return false, nil
}
