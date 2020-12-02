package interp

import (
	"fmt"

	"go.starlark.net/starlark"
)

func ExecFile(filename string) error {
	predefined := starlark.StringDict{
		"glob":            starlark.NewBuiltin("glob", FnGlob),
		"register_object": starlark.NewBuiltin("register_object", FnRegisterObject),
	}
	thread := &starlark.Thread{Name: filename, Print: printer, Load: loader()}
	_, err := starlark.ExecFile(thread, filename, nil, predefined)
	return err
}

func printer(thread *starlark.Thread, msg string) {
	fmt.Printf("[%s] %s", thread.Name, msg)
}
