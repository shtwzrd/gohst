package main

import (
	"testing"
	"time"
)

func TestToHistLine_NoExitStatus(t *testing.T) {
	now := time.Now()

	e := Invocation{}
	e.Command = "git log --graph --abbrev-commit --decorate --date=relative --all"
	e.Tags = []string{"git", "log", "graph"}
	e.Directory = "/home/soren/src/project"
	e.HasStatus = false
	e.Host = "laptop"
	e.Shell = "bash"
	e.User = "soren"
	e.Timestamp = now

	expected := "" + now.Format(time.UnixDate) +
		"sorenlaptop/home/soren/src/project" +
		"git log --graph --abbrev-commit --decorate --date=relative --all" +
		"[git log graph]"
	result := toHistLine(e)

	if result != expected {
		t.Fail()
	}
}

func TestToHistLine_WithExitStatus(t *testing.T) {
	now := time.Now()

	e := Invocation{}
	e.Command = "git log --graph --abbrev-commit --decorate --date=relative --all"
	e.Tags = []string{"git", "log", "graph"}
	e.Directory = "/home/soren/src/project"
	e.HasStatus = true
	e.Status = 0
	e.Host = "laptop"
	e.Shell = "bash"
	e.User = "soren"
	e.Timestamp = now

	expected := "" + now.Format(time.UnixDate) +
		"sorenlaptop/home/soren/src/project" +
		"git log --graph --abbrev-commit --decorate --date=relative --all" +
		"[git log graph]0\n"

	result := toHistLine(e)

	if result != expected {
		t.Fail()
	}
}
