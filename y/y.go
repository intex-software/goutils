package main

import (
	"flag"
	"fmt"
	"go/build"
	"os"
	"strconv"
	"strings"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/syscall"
	"github.com/traefik/yaegi/stdlib/unrestricted"
	"github.com/traefik/yaegi/stdlib/unsafe"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		showError(err)
		os.Exit(1)
	}
}

func run(arg []string) error {
	var install bool
	var transform bool
	var noAutoImport bool
	var tags string
	var err error

	// The following flags are initialized from environment.
	useSyscall, _ := strconv.ParseBool(os.Getenv("YAEGI_SYSCALL"))
	useUnrestricted, _ := strconv.ParseBool(os.Getenv("YAEGI_UNRESTRICTED"))
	useUnsafe, _ := strconv.ParseBool(os.Getenv("YAEGI_UNSAFE"))

	rflag := flag.NewFlagSet("y", flag.ContinueOnError)
	rflag.BoolVar(&install, "install", false, "install under GOPATH")
	rflag.BoolVar(&transform, "transform", false, "install the current path")
	rflag.BoolVar(&useSyscall, "syscall", useSyscall, "include syscall symbols")
	rflag.BoolVar(&useUnrestricted, "unrestricted", useUnrestricted, "include unrestricted symbols")
	rflag.StringVar(&tags, "tags", "", "set a list of build tags")
	rflag.BoolVar(&useUnsafe, "unsafe", useUnsafe, "include unsafe symbols")
	rflag.BoolVar(&noAutoImport, "noautoimport", false, "do not auto import pre-compiled packages. Import names that would result in collisions are automatically renamed\n(e.g. rand from \"crypto/rand\" and from \"math/rand\" to \"crypto_rand\" and \"math_rand\")")
	rflag.Usage = func() {
		fmt.Println("Usage: y [options] [path] [args]")
		fmt.Println("Options:")
		rflag.PrintDefaults()
	}
	if err = rflag.Parse(arg); err != nil {
		return err
	}
	args := rflag.Args()
	installArgs := ""

	i := interp.New(interp.Options{
		GoPath:       build.Default.GOPATH,
		BuildTags:    strings.Split(tags, ","),
		Env:          os.Environ(),
		Unrestricted: useUnrestricted,
	})
	if err := i.Use(stdlib.Symbols); err != nil {
		return err
	}
	if err := i.Use(interp.Symbols); err != nil {
		return err
	}
	if useSyscall {
		installArgs += " -syscall"
		if err := i.Use(syscall.Symbols); err != nil {
			return err
		}
		// Using a environment var allows a nested interpreter to import the syscall package.
		if err := os.Setenv("YAEGI_SYSCALL", "1"); err != nil {
			return err
		}
	}
	if useUnsafe {
		installArgs += " -unsafe"
		if err := i.Use(unsafe.Symbols); err != nil {
			return err
		}
		if err := os.Setenv("YAEGI_UNSAFE", "1"); err != nil {
			return err
		}
	}
	if useUnrestricted {
		installArgs += " -unrestricted"
		// Use of unrestricted symbols should always follow stdlib and syscall symbols, to update them.
		if err := i.Use(unrestricted.Symbols); err != nil {
			return err
		}
		if err := os.Setenv("YAEGI_UNRESTRICTED", "1"); err != nil {
			return err
		}
	}

	if len(args) == 0 {
		rflag.Usage()
		return nil
	}

	// Skip first os arg to set command line as expected by interpreted main.
	path := args[0]
	os.Args = arg
	flag.CommandLine = flag.NewFlagSet(path, flag.ExitOnError)

	if install || transform {
		if noAutoImport {
			installArgs += " -noautoimport"
		}
		if tags != "" {
			installArgs += " -tags " + tags
		}
		if install {
			return installScript(path, installArgs)
		} else if transform {
			return transformScript(path, installArgs)
		}
	} else if isFile(path) {
		err = runFile(i, path, noAutoImport)
	} else {
		_, err = i.EvalPath(path)
	}

	if err != nil {
		return err
	}

	return err
}

func isFile(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.Mode().IsRegular()
}

func runFile(i *interp.Interpreter, path string, noAutoImport bool) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var hasSheBang bool
	var script string
	if script = string(b); strings.HasPrefix(script, "#!") {
		script = strings.Replace(script, "#!", "//", 1)
		hasSheBang = true
	} else if strings.HasPrefix(script, "@y") {
		script = strings.Replace(script, "@", "//", 2)
		hasSheBang = true
	}

	if hasSheBang {
		if !noAutoImport {
			i.ImportUsed()
		}
		_, err = i.Eval(script)
		return err
	}

	// Files not starting with "#!" are supposed to be pure Go, directly Evaled.
	_, err = i.EvalPath(path)
	return err
}

func showError(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	if p, ok := err.(interp.Panic); ok {
		fmt.Fprintln(os.Stderr, string(p.Stack))
	}
}
