package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func dlog(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func main() {

	noAdd := flag.Bool("na", false, "Do not pass -a to git.")
	flag.Parse()

	args := flag.Args()

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

	if len(comps) >= 3 {
		ticket := comps[len(comps)-1]
		comps2 := strings.Split(ticket, "-")
		if len(comps2) == 2 {
			ticketPrefix = fmt.Sprintf("%s  ")
		}
	}

	message := strings.Join(args, " ")

	var commitParams []string

	commitParams = append(commitParams, "commit")
	if !*noAdd {
		commitParams = append(commitParams, "-a")
	}
	commitParams = append(commitParams, "-m")

	commitMessage := fmt.Sprintf("%s%s", ticketPrefix, message)
	commitParams = append(commitParams, commitMessage)

	fmt.Println(commitMessage)

	commitCmd := exec.Command("git", commitParams...)
	out, err = commitCmd.Output()
	if err != nil {
		dlog("Error: %v\n", err)
		dlog("Output: %v", string(out))
		return
	}
}
