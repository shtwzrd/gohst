package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
)

func getCommand(argv []string) (err error) {
	usage := `gohst -- your history, remote and secure.

Usage: gohst get [options] [--] [<searchterm>...]

options:
	-h, --help
	-t, --tag            history matching given tags
	-c, --cmd            history containing the given commands
  -s, --session        restrict search to particular session ids
  -d, --directory      filter for entries issued from particular directories
  -b, --before         restrict search to history before a given timestamp
  -a, --after          restrict search to history after a given timestamp
  -n, --count          return first n results [default=20]
	-v, --verbose        show metadata for history
	-A, --all            return all matching history results
  --local-only         search only the locally stored history index
	--refresh            force a call for latest from remote
`

	args, _ := docopt.Parse(usage, nil, true, "", false)
	fmt.Println(args)
	return
}
