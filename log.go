package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	g "github.com/warreq/gohstd/common"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func logCommand(argv []string, user string, repo g.CommandRepo) (err error) {
	usage := `gohst log; write commands to history
Usage:
	gohst log basic <cmd> <exitcode> [--FILE=<file>] [-f | --force]
	gohst log result <exitcode> [--FILE=<file>] [-f | --force]
	gohst log context <user> <host> <shell> <dir> <cmd> [--FILE=<file>]
	gohst log -h | --help

options:
	--FILE=<file>        alternate hist file, relative to home [default: .gohstry]
	-f, --force          write entry immediately to the remote
`

	arguments, _ := docopt.Parse(usage, argv, false, "", false)

	path := fmt.Sprintf("%s/%s", os.Getenv("HOME"), arguments["--FILE"].(string))
	index := Index{path}

	if arguments["-h"].(bool) || arguments["--help"].(bool) {
		fmt.Println(usage)
		os.Exit(0)
	}

	if arguments["basic"].(bool) {
		err = logBasic(arguments, index)
		if err == nil && arguments["--force"].(bool) {
			flush(user, index, repo)
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
			err = flush(user, index, repo)
			if err == nil {
				index.MarkSynced()
			}
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
