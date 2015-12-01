package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"strconv"
)

func getCommand(argv []string, user string, url string) (err error) {
	usage := `gohst.

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
  -A, --ALL                    return everything
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	verbose := args["--verbose"].(bool)

	count := 0
	if !args["--ALL"].(bool) {
		count, err = strconv.Atoi(args["--count"].(string))
		if err != nil {
			panic(fmt.Sprintf("%s is not an integer", args["num"].(string)))
		}
	}

	results := GetRequest(user, url, verbose, count)

	for i := len(results) - 1; i >= 0; i-- {
		fmt.Println(results[i])
	}

	return
}
