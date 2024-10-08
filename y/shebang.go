package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func init() {
	log.SetFlags(0)
}

const (
	sheBangBash    = "#! /usr/bin/env -S yaegi%s\n\n"
	windowsSheBang = "@y%s %%~f0 %%*\n@exit /b %%errorlevel%%\n\n"
)

func sheBang(args string) string {
	log.Println("Adding shebang")

	if runtime.GOOS == "windows" {
		return fmt.Sprintf(windowsSheBang,
			args,
		)
	}

	return fmt.Sprintf(sheBangBash,
		args,
	)
}

func target(source string) string {
	home := os.Getenv("GOPATH")
	target := filepath.Join(home, "bin", strings.TrimSuffix(filepath.Base(source), ".go"))

	if runtime.GOOS == "windows" {
		target += ".bat"
	}

	return target
}

func targetLocal(source string) string {
	home, _ := filepath.Abs(".")
	destination := filepath.Join(home, strings.TrimSuffix(filepath.Base(source), ".go"))

	if runtime.GOOS == "windows" {
		destination += ".bat"
	} else {
		destination += ".sh"
	}

	return destination
}

func installScript(source, args string) (err error) {
	destination := target(source)
	log.Println("Installing", source, "to", destination)
	return transfer(source, destination, args)
}

func transformScript(source, args string) error {
	destination := targetLocal(source)
	log.Println("Transform", source, "to", destination)
	return transfer(source, destination, args)
}

func transfer(source, destination string, args string) (err error) {
	log.Println("Reading source file")

	var script string
	if content, err := os.ReadFile(source); err != nil {
		return fmt.Errorf("error reading source file: %v", err)
	} else {
		script = string(content)
	}

	script = sheBang(args) + script

	if err := os.WriteFile(destination, []byte(script), 0o755); err != nil {
		return fmt.Errorf("error writing source file: %v", err)
	}

	log.Println("Created", destination)
	return
}
