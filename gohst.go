package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"os"
)

func main() {
	usage := `gohst.

Usage: gohst (-u USER) (-d URL) <command> [<args>...]
       gohst --version
       gohst -h | --help

Options:
   -d --domain=URL  Domain of the web service, e.g. gohst.herokuapp.com
   -u --user=USER
   -h --help
   --version

The supported gohst commands are:
   get        Search the command history
   log        Record a command to the history index
   flush      Push new entries in the history index to the remote
   forget     Remove an entry from the history
   stat       Display usage statistics and current session information
   tags       List all tags referenced in the command history
   commands   List all commands gohst has seen you use

See 'gohst help <command>' for more information on a specific command.
`
	args, _ := docopt.Parse(usage, nil, true, "gohst version 0.1", true)

	cmd := args["<command>"].(string)

	user := args["--user"].(string)
	domain := args["--domain"].(string)

	var url string
	if domain == "localhost" || domain == "127.0.0.1" {
		url = "http://" + domain
	} else {
		url = "https://" + domain
	}

	cmdArgs := args["<args>"].([]string)
	err := runCommand(cmd, cmdArgs, user, url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runCommand(cmd string, args []string, user string, url string) (err error) {
	argv := make([]string, 1)
	argv[0] = cmd
	argv = append(argv, args...)
	switch cmd {
	case "get":
		return getCommand(argv, user, url)
	case "flush":
		return flushCommand(argv, user, url)
	case "log":
		return logCommand(argv, user, url)
	case "help", "":
		fmt.Println("Not yet implemented")
	default:
		fmt.Println("Not yet implemented")
	}

	return fmt.Errorf("%s is not a gohst command. See 'gohst help'", cmd)
}
