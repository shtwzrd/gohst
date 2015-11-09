package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
)

func getCommand(argv []string) (err error) {
	usage := `gohst -- your history, remote and secure.

Usage:
  gohst get [options] [<searchterm>...]

options:
  -h, --help
  -t, --tag=(<tag>,..)         history matching given tags
  -s, --session=(<session>,..) restrict search to particular session ids
  --shell=(<shell>,..)         restrict search to entries from particular shells
  -d, --dir=(<dir),..)         filter for entries issued from certain directories
  -b, --before=<time>          restrict search to history before a given timestamp
  -a, --after=<time>           restrict search to history after a given timestamp
  -n, --count=<num>            return first n results [default: 100]
  -v, --verbose                show metadata for history
  -x, --exclude-fail           filter out entries with non-0 exit statuses
  -X, --exclude-success        filter out entries with a 0 exit status
  -a, --all                    return all matching history results
  -A, --ALL                    return EVERYTHING
  --local-only                 search only the locally stored history index
  --refresh                    force a call for latest from remote
`

	args, _ := docopt.Parse(usage, nil, true, "", false)
	fmt.Println(args)
	return
}
