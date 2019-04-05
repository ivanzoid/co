package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"
	"regexp"
	"strings"
)

func dlog(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintf(os.Stderr, "\n")
}

func vlog(format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stderr, format, args...)
		fmt.Fprintf(os.Stderr, "\n")
	}
}

var (
	verbose bool
)

func main() {

	noAdd := flag.Bool("na", false, "Do not pass -a to git.")
	flag.BoolVar(&verbose, "v", false, "Verbose")
	flag.Parse()

	args := flag.Args()

	branchName, err := runProgram1("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		dlog("Error getting current branch: %v\n", err)
	}

	ticketPrefix := ticketNumber(branchName)

	message := strings.Join(args, " ")

	var commitParams []string

	commitParams = append(commitParams, "commit")
	if !*noAdd {
		commitParams = append(commitParams, "-a")
	}
	commitParams = append(commitParams, "-m")

	commitMessage := fmt.Sprintf("%s  %s", ticketPrefix, message)
	commitParams = append(commitParams, commitMessage)

	fmt.Println(commitMessage)

	out, err := runProgram1("git", commitParams...)
	if err != nil {
		dlog("Error: %v\n", err)
		dlog("Output: %v", string(out))
		return
	}
}

func ticketNumber(branchName string) string {

	regexps := regexps()

	vlog("Regexps: %v", regexps)
	vlog("Branch name: %v", branchName)

	for _, r := range regexps {
		re, err := regexp.Compile(r)
		if err != nil {
			continue
		}
		results := re.FindStringSubmatch(branchName)

		vlog("Results: %v", results)

		if len(results) == 2 {
			return results[1]
		} else {
			return ""
		}
	}

	return ""
}

func readLinesFromFile(path string) (lines []string, err error) {

	vlog("Reading %v", path)

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return
}

func regexps() []string {
	user, _ := user.Current()
	filePath := path.Join(user.HomeDir, ".co", "config")
	regexps, err := readLinesFromFile(filePath)
	if err != nil {
		dlog("Can't read config %v: %v", filePath, err)
	}
	return regexps
}

func runProgram1(program string, args ...string) (string, error) {
	outStrings, err := runProgram(program, args...)
	if err != nil {
		return "", err
	}
	if len(outStrings) == 0 {
		return "", nil
	}

	return outStrings[0], nil
}

func runProgram(program string, args ...string) ([]string, error) {

	vlog("Running %v", cmdString(program, args))

	cmd := exec.Command(program, args...)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	outString := string(out)
	if len(outString) == 0 {
		return nil, nil
	}

	outStrings := strings.Split(outString, "\n")
	return outStrings, nil
}

// Utils

func cmdString(program string, args []string) string {
	comps := make([]string, 0)

	comps = append(comps, program)

	for _, arg := range args {
		if strings.Contains(arg, " ") {
			comps = append(comps, fmt.Sprintf("'%v'", arg))
		} else {
			comps = append(comps, arg)
		}
	}

	result := strings.Join(comps, " ")
	return result
}
