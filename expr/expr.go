package expr

import (
	"errors"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/intex-software/goutils/expr/modules"
)

type interp struct {
	vm   *vm.VM
	prog *vm.Program
	ops  []expr.Option
}

var (
	ErrNoProgram = errors.New("no program")
)

func New(ops ...expr.Option) (t *interp, err error) {
	fmt := modules.Fmt()

	t = &interp{
		vm:  &vm.VM{},
		ops: make([]expr.Option, 0, len(ops)+len(fmt)),
	}

	t.ops = append(t.ops, ops...)
	t.ops = append(t.ops, fmt...)

	return
}

func (t *interp) Compile(input string) (err error) {
	t.prog, err = expr.Compile(input, t.ops...)
	return
}

func (t *interp) Run(env any) (result any, err error) {
	if t.prog == nil {
		return nil, ErrNoProgram
	}

	result, err = t.vm.Run(t.prog, t.ops)
	if err != nil {
		return nil, err
	}

	return result, nil
}
