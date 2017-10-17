package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func dlog(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func main() {
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := branchCmd.Output()
	if err != nil {
		dlog("Error: %v\n", err)
		return
	}

	outStr := string(out)
	outStr = strings.TrimSpace(outStr)
	comps := strings.Split(outStr, "/")

	ticketPrefix := ""

	if len(comps) >= 2 {
		ticketPrefix = fmt.Sprintf("%s ", comps[len(comps)-1])
	}

	message := strings.Join(os.Args[1:], " ")

	var commitParams []string

	commitParams = append(commitParams, "commit")
	commitParams = append(commitParams, "-a")
	commitParams = append(commitParams, "-m")

	commitMessage := fmt.Sprintf("\"%s%s\"", ticketPrefix, message)
	commitParams = append(commitParams, commitMessage)

	commitCmd := exec.Command("git", commitParams...)
	out, err = commitCmd.Output()
	if err != nil {
		dlog("Error: %v\n", err)
		dlog("Output: %v", string(out))
		return
	}
}
