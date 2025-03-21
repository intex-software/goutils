package javascript

import (
	"github.com/dop251/goja"

	"github.com/intex-software/goutils/javascript/modules"
)

func New() (vm *goja.Runtime, err error) {
	vm = goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	if err = modules.FmtModule(vm); err != nil {
		return
	}

	return
}
