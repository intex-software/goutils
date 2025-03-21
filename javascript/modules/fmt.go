package modules

import (
	"fmt"

	"github.com/dop251/goja"
)

func args(args []goja.Value) []any {
	v := make([]any, len(args))
	for i, arg := range args {
		v[i] = arg.Export()
	}
	return v
}

func FmtModule(vm *goja.Runtime) error {
	module := vm.NewObject()
	module.Set("sprintf", func(call goja.FunctionCall) goja.Value {
		format := call.Argument(0).String()
		return vm.ToValue(fmt.Sprintf(format, args(call.Arguments[1:])...))
	})
	exception := func(call goja.FunctionCall) goja.Value {
		format := call.Argument(0).String()
		return vm.ToValue(fmt.Errorf(format, args(call.Arguments[1:])...))
	}
	module.Set("exception", exception)
	module.Set("errorf", exception)
	return vm.Set("fmt", module)
}
