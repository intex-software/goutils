package main

import (
	"fmt"
	"os"

	"github.com/intex-software/goutils/git/version"
)

func main() {
	if len(os.Args) < 2 {
		panic("Usage: version [windows|version]")
	}

	switch os.Args[1] {
	case "windows":
		fmt.Println(version.Windows())
	case "version":
		fmt.Println(version.Git())
	default:
		fmt.Println(version.Git())
	}
}
