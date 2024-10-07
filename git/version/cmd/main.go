package main

import (
	"fmt"
	"os"

	"github.com/intex-software/goutils/git/version"
)

func main() {
	windows := len(os.Args) > 1 && os.Args[1] == "windows"
	if ver, err := value(windows); err != nil {
		panic(err)
	} else {
		fmt.Println(ver)
	}
}

func value(windows bool) (string, error) {
	if windows {
		return version.Windows()
	} else {
		return version.Git()
	}
}
