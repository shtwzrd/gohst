package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"github.com/docopt/docopt-go"
	"io/ioutil"
	"os"
	"path/filepath"
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

	dir := DefaultDir()

	key := LoadKeyFile(dir)
	nokey := len(key) == 0

	reader := bufio.NewReader(os.Stdin)
	if nokey {
		fmt.Printf("No %s found in %s. Generate a new one? [Y/n] ", Key, dir)
	} else {
		fmt.Printf("Found %s in %s. Back up the existing file before attempting to generate a new one.\n", Key, dir)
		os.Exit(0)
	}
	yn, _, err := reader.ReadRune()
	fmt.Println("")
	if err != nil {
		os.Exit(1)
	}
	if yn == 'Y' || yn == 'y' {
		newKey, err := GenerateKey()
		if err != nil {
			fmt.Printf("Error: Failed to generate a key: %s\n", err)
			os.Exit(1)
		}
		err = ioutil.WriteFile(filepath.Join(dir, Key), newKey, os.ModePerm)
		if err != nil {
			fmt.Printf("Error: Could not write file: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully generated key file %s in %s\n", Key, dir)
		os.Exit(0)
	}

	return
}

func GenerateKey() ([]byte, error) {
	key := make([]byte, 64)

	_, err := rand.Read(key)
	return key, err
}
