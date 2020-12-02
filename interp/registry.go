package interp

import (
	"fmt"

	"go.starlark.net/starlark"
)

func FnRegisterObject(
	thread *starlark.Thread,
	fn *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {

	var (
		protoMessage *ProtoMessageType
		srcFiles     starlark.Iterable
	)

	if err := starlark.UnpackArgs(fn.Name(), args, kwargs, "object", &protoMessage, "srcs", &srcFiles); err != nil {
		return nil, err
	}

	fmt.Println(protoMessage)
	fmt.Println(srcFiles)
	return starlark.None, nil
}
