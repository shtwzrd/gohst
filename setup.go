package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"os"
)

func setup(argv []string, cfg Config) (err error) {
	usage := `gohst setup; interactively generate a configuration file
Usage:
	gohst setup
  gohst setup -h

options:
	-h, --help
`
	arguments, _ := docopt.Parse(usage, argv, false, "", false)

	if arguments["--help"].(bool) {
		fmt.Println(usage)
		os.Exit(0)
	}
	fmt.Println(cfg)

	return
}
