package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func logCommand(argv []string) (err error) {
	usage := `gohst log; write commands to history
Usage:
	gohst log basic <cmd> <exitcode>
	gohst log result <exitcode>
	gohst log context <user> <host> <shell> <dir> <cmd>
	gohst log -h | --help | --version | -f | --force | --FILE=<file>

options:
	-h, --help
	<user>               the user issuing the command
	<host>               the hostname identifying the machine
	<shell>              the name of the shell from which the command is invoked
	<dir>                the directory from which the command was invoked
	<exitcode>           the exit code of the command
	<cmd>                the executed command
	--FILE=<file>        alternate history file [default: ~/.gohstry]
	-f, --force          write entry immediately to the remote
`

	arguments, _ := docopt.Parse(usage, nil, true, "", false)
	fmt.Println(arguments)

	path := fmt.Sprintf("%s/%s", os.Getenv("HOME"), ".gohstry")
	index := LocalRepo{path}

	if arguments["basic"].(bool) {
		return logBasic(arguments, index)
	}

	if arguments["context"].(bool) {
		return logContext(arguments, index)
	}

	if arguments["result"].(bool) {
		return logResult(arguments, index)
	}

	return
}

func logBasic(args map[string]interface{}, index Repo) (err error) {
	cmd, tags := parseOutTags(args["<cmd>"].(string))

	e := Invocation{}
	e.Timestamp = time.Now().UTC()
	e.Command = cmd
	e.Tags = tags
	e.Status = getResult(args)
	e.HasStatus = true

	return index.Write(e)
}

func logContext(args map[string]interface{}, index Repo) (err error) {
	cmd, tags := parseOutTags(args["<cmd>"].(string))

	e := Invocation{}
	e.Timestamp = time.Now().UTC()
	e.Command = cmd
	e.Tags = tags
	e.Directory = args["<dir>"].(string)
	e.Host = args["<host>"].(string)
	e.Shell = args["<shell>"].(string)
	e.User = args["<user>"].(string)

	return index.Write(e)
}

func logResult(args map[string]interface{}, index Repo) (err error) {
	e := Invocation{}
	e.Timestamp = time.Now().UTC()
	e.Status = getResult(args)
	e.HasStatus = true

	return index.Write(e)
}

func getResult(args map[string]interface{}) (result int8) {
	if val, exists := args["<status>"]; exists {
		status, err := strconv.ParseInt(val.(string), 10, 8)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		result = int8(status)
	}
	return
}

func parseOutTags(input string) (command string, tags []string) {
	// Regex to remove all single and double quoted substrings
	re := regexp.MustCompile("['][^']*[']|[\"][^\"]*[\"]")
	s := re.ReplaceAllString(input, "")
	tags = strings.Fields(strings.Split(s, "#")[1])
	command = input[0:strings.LastIndex(input, "#")]
	return
}
