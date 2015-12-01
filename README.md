#gohst

A command-line client for logging and querying your command line history across all of your machines.

## Building from source

1. Install Go for your platform
2. Set the GOPATH environment variable
3. Fetch gohst
Install it via `go get` like so:

    `go get github.com/warreq/gohst`

_or_

    git clone https://github.com/warreq/gohst
    cd gohst 
    go build -o gohst

## Installation

1. Drop the gohst binary somewhere on your PATH
2. Edit the relevant script for your shell, e.g. `bash_example_config.sh`
   * 'gohst -u test -d gohst.herokuapp.com' should be changed to match your username and service location
   * if you want to write history to another file than the standard ~/.gohstry, specify it with the --FILE= switch
   * configure any other preferences you may have; see `gohst -h` for details
3. Source your script or place its contents in your shrc

## Usage

Once installed, commands issued from your shell will be written to the gohst history file and periodically synced to the remote service

You can apply tags to your commands by commenting them like so: 

    git log --graph --abbrev-commit --pretty=oneline origin..mybranch #commit-history git-log

which will make your commands searchable by the tags you've defined, in addition to standard keyword search

You can also omit the logging of individual commands by placing #shh after them

    export DATABASE_CONNECTIONSTRING=mysql://secret:secret@url:port/db #shh

Use `gohst get` to query your history by keyword, tag, host, timestamp and more

or just type `gohst get` with no args to fetch your last 100 commands from any machine

See `gohst get --help` for comprehensive query information 
