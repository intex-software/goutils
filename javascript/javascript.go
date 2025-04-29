package javascript

import (
	"fmt"
	"sync"

	"github.com/dop251/goja"

	"github.com/intex-software/goutils/javascript/modules"
)

type interp struct {
	vm *goja.Runtime
	m  sync.Mutex
}

// New creates a new JavaScript interpreter using goja.
func New() (*interp, error) {
	t := &interp{vm: goja.New()}

	t.vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	if err := modules.FmtModule(t.vm); err != nil {
		return nil, err
	}

	return t, nil
}

// Compile compiles the given JavaScript source code and returns the result.
func (t *interp) Compile(source string) (r goja.Value, err error) {
	r, err = t.vm.RunString(source)
	if err != nil {
		return
	}
	return
}

type Functor func(...any) (any, error)

// GetCall retrieves a function by name from the JavaScript interpreter and returns it as a Functor.
func (t *interp) GetCall(name string) (adapter Functor, err error) {
	var functor func(...any) any

	if value := t.vm.Get(name); value == nil {
		err = fmt.Errorf("function not found: '%s'", name)
		return
	} else if err = t.vm.ExportTo(value, &functor); err != nil {
		return
	} else {
		adapter = func(args ...any) (output any, err error) {
			t.m.Lock()
			defer t.m.Unlock()

			if ex := t.vm.Try(func() {
				output = functor(args...)
			}); ex != nil {
				err = ex
			}
			return
		}
	}

	return
}
