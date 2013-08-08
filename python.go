package main

import (
	"fmt"
	"os/exec"
	"regexp"
)

// Names of python executables to try, in preference order.
var (
	candidatePythons = []string{"python", "python2", "python2.7", "python2.6", "python2.5", "python2.4"}
	pyVersionRegex = regexp.MustCompile(`^Python (\d)\.(\d)\.(\d)`)
)

// findPython locates a workable version of python in the user's path.
// Although Pygments docs claim Pygments works with python 3:
//   http://pygments.org/faq/#python3
// pygmentize isn't actually python3 compatible.
func findPython() (string, error) {
	for _, c := range candidatePythons {
		path, err := exec.LookPath(c)
		if err != nil {
			continue
		}
		out, err := exec.Command(path, "--version").CombinedOutput()
		if err != nil {
			continue
		}
		parts := pyVersionRegex.FindAllSubmatch(out, -1)
		if len(parts) != 1 || len(parts[0]) != 4 {
			continue
		}
		match := parts[0]
		if match[1][0] == '2' && match[2][0] >= '4' {
			return path, nil
		}
	}
	return "", fmt.Errorf("No suitable Python version was found.")
}
