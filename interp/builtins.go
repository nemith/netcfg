package interp

import (
	"fmt"

	"github.com/mattn/go-zglob"
	"go.starlark.net/starlark"
)

func FnExport(thread *starlark.Thread,
	fn *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {
	return nil, nil
}

type register struct {
	objects []string
}

func FnGlob(
	thread *starlark.Thread,
	fn *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {

	var (
		patterns *starlark.List
	)

	if err := starlark.UnpackArgs(fn.Name(), args, kwargs, "", &patterns); err != nil {
		return nil, err
	}

	iter := patterns.Iterate()
	defer iter.Done()

	var matchList starlark.List
	var val starlark.Value
	for iter.Next(&val) {
		pattern, ok := val.(starlark.String)
		if !ok {
			return nil, fmt.Errorf("value not a string")
		}

		matches, err := zglob.Glob(string(pattern))
		if err != nil {
			return nil, err
		}
		for _, match := range matches {
			if err := matchList.Append(starlark.String(match)); err != nil {
				return nil, err
			}
		}
	}
	return &matchList, nil
}
