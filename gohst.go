package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"os"
	"strings"
)

func main() {
	usage := `gohst.

Usage: gohst [(-u USER) (-p PASS) (-d URL)] <command> [<args>...]
       gohst --version
       gohst -h | --help

Options:
   -d --domain=URL  Domain of the web service, e.g. gohst.herokuapp.com
   -u --user=USER
   -p --password=PASS
   -h --help
   --version

The supported gohst commands are:
   get        Search the command history
   log        Record a command to the history index
   flush      Push new entries in the history index to the remote
   forget     Remove an entry from the history
   stat       Display usage statistics and current session information
   tags       List all tags referenced in the command history

See 'gohst <command> --help' for more information on a specific command.
`
	args, _ := docopt.Parse(usage, nil, true, "gohst version 0.1", true)

	cmd := args["<command>"].(string)

	user := args["--user"].(string)
	domain := args["--domain"].(string)
	password := args["--password"].(string)

	var url string
	if strings.HasPrefix(domain, "localhost") ||
		strings.HasPrefix(domain, "127.0.0.1") {
		url = "http://" + domain
	} else {
		url = "https://" + domain
	}

	service := NewService(user, password, url)
	cmdArgs := args["<args>"].([]string)
	err := runCommand(cmd, cmdArgs, user, service)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runCommand(cmd string, args []string, user string, serv Service) error {
	cmdRepo := HttpCommandRepo{serv}
	argv := make([]string, 1)
	argv[0] = cmd
	argv = append(argv, args...)
	switch cmd {
	case "get":
		return getCommand(argv, user, cmdRepo)
	case "flush":
		return flushCommand(argv, user, cmdRepo)
	case "log":
		return logCommand(argv, user, cmdRepo)
	default:
		fmt.Println("Not yet implemented")
	}

	return fmt.Errorf("%s is not a gohst command. See 'gohst help'", cmd)
}
