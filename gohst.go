package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"os"
	"strings"
)

func main() {
	usage := `gohst.

Usage: gohst [-u USER -p PASS -d URL --dir=CFG_DIR] <command> [<args>...]
       gohst --version
       gohst -h | --help

Options:
   -d --domain=URL     Domain of the web service, e.g. gohst.herokuapp.com
   -u --user=USER
   -p --password=PASS
   --dir=CFG_DIR       Alternative directory for fetching configuration
   -h --help
   --version

The supported gohst commands are:
   get        Search the command history
   log        Record a command to the history index
   flush      Push new entries in the history index to the remote
   forget     Remove an entry from the history
   setup      Interactively generate a configuration file
   stat       Display usage statistics and current session information
   tags       List all tags referenced in the command history

See 'gohst <command> --help' for more information on a specific command.
`
	args, _ := docopt.Parse(usage, nil, true, "gohst version 0.1", true)

	cmd := args["<command>"].(string)
	cfg := MergeConfig(args)

	url := DetectProtocol(cfg.Domain)

	service := NewService(url, cfg.Password)
	cmdArgs := args["<args>"].([]string)
	err := RunCommand(cmd, cmdArgs, cfg, service)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func MergeConfig(args map[string]interface{}) Config {
	cfg := Config{}

	if dir := args["--dir"]; dir != nil {
		cfg = LoadConfig(dir.(string))
	} else {
		cfg = LoadConfig(DefaultDir())
	}

	if user := args["--user"]; user != nil {
		cfg.Username = user.(string)
	}
	if password := args["--password"]; password != nil {
		cfg.Password = password.(string)
	}
	if domain := args["--domain"]; domain != nil {
		cfg.Domain = domain.(string)
	}
	return cfg
}

func DetectProtocol(domain string) (url string) {
	if strings.HasPrefix(domain, "localhost") ||
		strings.HasPrefix(domain, "127.0.0.1") {
		url = "http://" + domain
	} else {
		url = "https://" + domain
	}
	return
}

func RunCommand(cmd string, args []string, cfg Config, serv Service) error {
	cmdRepo := NewHttpCommandRepo(serv, cfg.Key)
	argv := make([]string, 1)
	argv[0] = cmd
	argv = append(argv, args...)
	switch cmd {
	case "get":
		return getCommand(argv, cfg, cmdRepo)
	case "flush":
		return flushCommand(argv, cfg, cmdRepo)
	case "log":
		return logCommand(argv, cfg, cmdRepo)
	case "setup":
		return setup(argv, cfg)
	default:
		fmt.Println("Not yet implemented")
	}

	return fmt.Errorf("%s is not a gohst command. See 'gohst help'", cmd)
}
