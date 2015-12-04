package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"os"
)

func flushCommand(argv []string, user string, url string) (err error) {
	usage := `gohst flush; sync history with the remote
Usage:
	gohst flush [options]

options:
	-h, --help
	--FILE=<file>        alternate hist file, relative to home [default: .gohstry]
`

	arguments, _ := docopt.Parse(usage, argv, false, "", false)

	if arguments["-h"].(bool) || arguments["--help"].(bool) {
		fmt.Println(usage)
		os.Exit(0)
	}

	path := fmt.Sprintf("%s/%s", os.Getenv("HOME"), arguments["--FILE"].(string))
	index := Index{path}

	FlushRequest(user, url, index)
	return index.MarkSynced()
}
