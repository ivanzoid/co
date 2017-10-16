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
	comps := strings.Split(outStr, "/")

	ticketPrefix := ""

	if len(comps) >= 2 {
		ticketPrefix = fmt.Sprintf("%s ", comps[len(comps)-1])
	}

	message := strings.Join(os.Args[1:], " ")

	commitCmdParamsString := fmt.Sprintf("commit -a \"%s%s\"", ticketPrefix, message)

	fmt.Printf("%s", commitCmdParamsString)

	//commitCmd := exec.Command(commitCmdString)
	//if err != nil {
	//	dlog("Error: %v\n", err)
	//	return
	//}

}
