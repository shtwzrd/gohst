package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	g "github.com/warreq/gohstd/common"
	"os"
	"path/filepath"
)

func flushCommand(argv []string, cfg Config, repo g.CommandRepo) (err error) {
	usage := `gohst flush; sync history with the remote
Usage:
	gohst flush [options]

options:
	-h, --help
	--FILE=<file>        alternate hist file, relative to home
`

	arguments, _ := docopt.Parse(usage, argv, false, "", false)

	if arguments["-h"].(bool) || arguments["--help"].(bool) {
		fmt.Println(usage)
		os.Exit(0)
	}

	var path string
	if arguments["--FILE"] != nil {
		path = filepath.Join(os.Getenv("HOME"), arguments["--FILE"].(string))
	} else {
		path = filepath.Join(DefaultDir(), Hist)
	}
	index := Index{path}

	err = flush(cfg.Username, index, repo)
	if err != nil {
		return
	}
	return index.MarkSynced()
}

func flush(user string, index Index, repo g.CommandRepo) error {
	unsynced, err := index.GetUnsynced()
	if err != nil {
		panic(fmt.Sprintf("[gohst] %s: %s\n", "Invalid Hist File Error: ", err))
	}

	payload := make(g.Invocations, len(unsynced))

	for i, v := range unsynced {
		payload[i] = v.ToInvocation()
	}
	return repo.InsertInvocations(user, payload)
}
