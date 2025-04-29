package modules

import (
	"fmt"

	"github.com/expr-lang/expr"
)

func Fmt() []expr.Option {
	return []expr.Option{
		print(),
		println(),
		printf(),
		sprintf(),
		errorf(),
	}
}

func sprintf() expr.Option {
	return expr.Function("sprintf", func(args ...any) (any, error) {
		return fmt.Sprintf(args[0].(string), args[1:]...), nil
	})
}

func errorf() expr.Option {
	return expr.Function("errorf", func(args ...any) (any, error) {
		return fmt.Sprintf(args[0].(string), args[1:]...), nil
	})
}

func print() expr.Option {
	return expr.Function("print", func(args ...any) (any, error) {
		fmt.Print(args...)
		return nil, nil
	})
}

func println() expr.Option {
	return expr.Function("println", func(args ...any) (any, error) {
		fmt.Println(args...)
		return nil, nil
	})
}

func printf() expr.Option {
	return expr.Function("printf", func(args ...any) (any, error) {
		fmt.Printf(args[0].(string), args...)
		return nil, nil
	})
}
