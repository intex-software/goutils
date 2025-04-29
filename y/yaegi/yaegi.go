package yaegi

import (
	"fmt"
	"reflect"

	"github.com/traefik/yaegi/interp"
)

type yaegi struct {
	vm *interp.Interpreter
}

// Symbols is a map of symbols to be used in the yaegi interpreter.
func New() *yaegi {
	t := &yaegi{vm: interp.New(interp.Options{Unrestricted: false})}
	t.vm.Use(Symbols)
	t.vm.ImportUsed()
	return t
}

// NewUnrestricted creates a new yaegi interpreter with unrestricted mode enabled.
func NewUnrestricted() *yaegi {
	t := &yaegi{vm: interp.New(interp.Options{Unrestricted: true})}
	t.vm.Use(Symbols)
	t.vm.ImportUsed()
	return t
}

// Compile compiles the given Go source code and returns the result.
func (y *yaegi) Compile(source string) (reflect.Value, error) {
	return y.vm.Eval(source)
}

type Functor func(...any) (any, error)

// name is the name of the function to be called, maybe containing a package name (e.g. "main.Adapter").
func (y *yaegi) GetCall(name string) (adapter Functor, err error) {
	var ok bool
	var v reflect.Value
	if v, err = y.vm.Eval(name); err != nil {
		return
	} else if adapter, ok = v.Interface().(Functor); !ok {
		err = fmt.Errorf("oci adapter has the wrong interface: '%T' != 'Adapter'", v.Interface())
		return
	}
	return
}
