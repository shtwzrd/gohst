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

func logCommand(argv []string, user string, url string) (err error) {
	usage := `gohst log; write commands to history
Usage:
	gohst log basic [options] <cmd> <exitcode>
	gohst log result [options] <exitcode>
	gohst log context [options] <user> <host> <shell> <dir> <cmd>

options:
	-h, --help
	<user>               the user issuing the command
	<host>               the hostname identifying the machine
	<shell>              the name of the shell from which the command is invoked
	<dir>                the directory from which the command was invoked
	<exitcode>           the exit code of the command
	<cmd>                the executed command
	--FILE=<file>        alternate hist file, relative to home [default: .gohstry]
	-f, --force          write entry immediately to the remote
`

	arguments, _ := docopt.Parse(usage, argv, true, "", false)

	path := fmt.Sprintf("%s/%s", os.Getenv("HOME"), arguments["--FILE"].(string))
	index := Index{path}

	if arguments["basic"].(bool) {
		err = logBasic(arguments, index)
		if err == nil && arguments["--force"].(bool) {
			FlushRequest(user, url, index)
			index.MarkSynced()
		}
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	if arguments["context"].(bool) {
		return logContext(arguments, index)
	}

	if arguments["result"].(bool) {
		err = logResult(arguments, index)
		if err == nil && arguments["--force"].(bool) {
			FlushRequest(user, url, index)
			index.MarkSynced()
		}
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	return
}

func logBasic(args map[string]interface{}, index Index) (err error) {
	cmd, tags := parseOutTags(args["<cmd>"].(string))

	if hasSilentTag(tags) {
		return nil
	}

	index.Flush()
	e := IndexEntry{}
	e.Timestamp = time.Now().UTC()
	e.Command = cmd
	e.Tags = tags
	e.Status = getResult(args)
	e.HasStatus = true

	err = index.Write(e)
	return
}

func logContext(args map[string]interface{}, index Index) (err error) {
	if index.canWriteContext() {
		index.Flush()
		cmd, tags := parseOutTags(args["<cmd>"].(string))

		if hasSilentTag(tags) {
			index.Write(dummyContext())
			return
		}

		e := IndexEntry{}
		e.Timestamp = time.Now().UTC()
		e.Command = cmd
		e.Tags = tags
		e.Directory = args["<dir>"].(string)
		e.Host = args["<host>"].(string)
		e.Shell = args["<shell>"].(string)
		e.User = args["<user>"].(string)

		return index.Write(e)
	}
	return
}

func logResult(args map[string]interface{}, index Index) (err error) {
	if index.canWriteResult() {

		e := IndexEntry{}
		e.Status = getResult(args)
		e.HasStatus = true

		err = index.Write(e)
	}
	return
}

func getResult(args map[string]interface{}) (result int8) {
	if val, exists := args["<exitcode>"]; exists {
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
	cmdSplit := strings.Split(s, "#")
	if len(cmdSplit) > 1 {
		tags = strings.Fields(cmdSplit[1])
		command = input[0:strings.LastIndex(input, "#")]
	} else {
		command = input
	}
	return
}

func hasSilentTag(input []string) bool {
	for _, t := range input {
		if strings.HasPrefix(t, "shh") {
			return true
		}
	}
	return false
}

func dummyContext() IndexEntry {
	e := IndexEntry{}
	e.Timestamp = time.Now().UTC()
	e.Command = "null"
	e.Tags = nil
	e.Directory = "null"
	e.Host = "null"
	e.Shell = "null"
	e.User = "null"
	return e
}
