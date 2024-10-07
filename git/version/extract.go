package version

import (
	"os/exec"
	"strings"
	"sync"
)

func Git() (string, error) {
	return gitVersionOnce()
}

func Windows() (string, error) {
	return gitWindowsVersionOnce()
}

var gitWindowsVersionOnce = sync.OnceValues[string, error](func() (version string, err error) {
	_gitVersion, err := gitVersionOnce()
	if err != nil {
		return
	}

	version, _, _ = strings.Cut(_gitVersion, "-g")
	if strings.Index(version, "+") > 0 {
		version = strings.Replace(version, "+", ".", 1)
	}
	for strings.Count(version, ".") < 3 {
		version += ".0"
	}

	for strings.Count(version, ".") > 3 {
		version = version[:strings.LastIndexByte(version, '.')]
	}

	return
})

var gitVersionOnce = sync.OnceValues[string, error](func() (version string, err error) {
	gvBytes, err := git("describe", "--tags", "--always", "--match=v*")
	if err != nil {
		return
	}

	output, _, _ := strings.Cut(strings.TrimSpace(string(gvBytes))[1:], "-g")
	if strings.ContainsRune(output, '+') {
		version = strings.Replace(output, "-", ".", 1)
	} else {
		version = strings.Replace(output, "-", "+", 1)
	}

	return
})

func git(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	return cmd.CombinedOutput()
}
