package interp

import (
	"fmt"
	"path/filepath"

	"github.com/jhump/protoreflect/desc/protoparse"
	"go.starlark.net/starlark"
)

func loader() func(*starlark.Thread, string) (starlark.StringDict, error) {
	type entry struct {
		globals starlark.StringDict
		err     error
	}

	var cache = make(map[string]*entry)
	return func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
		e, ok := cache[module]

		if e == nil {
			if ok {
				return nil, fmt.Errorf("cycle in load graph")
			}

			// Add a placeholder to indicate "load in progress".
			cache[module] = nil

			var globals starlark.StringDict
			var err error
			switch ext := filepath.Ext(module); ext {
			case ".proto":
				globals, err = loadProto(thread, module)
			case ".star":
				globals, err = loadStar(thread, module)
			default:
				return nil, fmt.Errorf("unknown file type: %s", ext)
			}

			e = &entry{globals: globals, err: err}
			cache[module] = e
		}
		return e.globals, e.err
	}
}

func loadStar(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	t := &starlark.Thread{Name: module, Load: thread.Load}
	return starlark.ExecFile(t, module, nil, nil)
}

func loadProto(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	parser := protoparse.Parser{}
	fds, err := parser.ParseFiles(module)
	if err != nil {
		return nil, err
	}

	fd := fds[0]
	globals := make(starlark.StringDict)
	for _, mdesc := range fd.GetMessageTypes() {
		globals[mdesc.GetName()] = NewProtoMessageType(mdesc, module)
	}

	for _, edesc := range fd.GetEnumTypes() {
		globals[edesc.GetName()] = NewProtoEnumType(edesc)
	}

	return globals, err
}
