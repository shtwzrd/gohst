package main

import (
	gohst "github.com/warreq/gohstd/src/gohstd/common"
	"time"
)

type IndexEntry struct {
	Directory string
	Shell     string
	User      string
	Host      string
	Command   string
	Status    int8
	Tags      []string
	Timestamp time.Time
	HasStatus bool
	IsSynced  bool
}

func (e IndexEntry) ToInvocation() (inv gohst.Invocation) {
	inv.Timestamp = e.Timestamp
	inv.Tags = e.Tags
	inv.ExitCode = e.Status
	inv.Host = e.Host
	inv.User = e.User
	inv.Shell = e.Shell
	inv.Command = e.Command
	inv.Directory = e.Directory
	return
}
