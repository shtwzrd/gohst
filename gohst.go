package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"os"
)

func main() {
	usage := `gohst -- your history, remote and secure.

Usage: gohst [--user=<USER>] [--url=<URL>] [--version] [-h|--help] <command> [<args>...]

options:
   -h, --help
  --user=<user>					user that stores or retrieves history
  --url=<url>					url for the web service	

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

	//user := args["<user>"].(string)
	//url := args["<url>"].(string)

	//Hard coded values, as we had issues with docopt parsing
	user := "test"
	url := "http://gohst.herokuapp.com"

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
	case "log":
		return logCommand(argv)
	case "help", "":
		fmt.Println("Not yet implemented")
	default:
		fmt.Println("Not yet implemented")
	}

	return fmt.Errorf("%s is not a gohst command. See 'gohst help'", cmd)
}
