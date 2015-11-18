package main

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestToHistLine_NoExitStatus(t *testing.T) {
	now := time.Now()

	e := IndexEntry{}
	e.Command = "git log --graph --abbrev-commit --decorate --date=relative --all"
	e.Tags = []string{"git", "log", "graph"}
	e.Directory = "/home/soren/src/project"
	e.HasStatus = false
	e.Host = "laptop"
	e.Shell = "bash"
	e.User = "soren"
	e.Timestamp = now

	expected := "" + now.Format(time.UnixDate) +
		"sorenlaptopbash/home/soren/src/project" +
		"git log --graph --abbrev-commit --decorate --date=relative --all" +
		"[git log graph]"
	result := toHistLine(e)

	if result != expected {
		t.Fail()
	}
}

func TestToHistLine_WithExitStatus(t *testing.T) {
	now := time.Now()

	e := IndexEntry{}
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
		"sorenlaptopbash/home/soren/src/project" +
		"git log --graph --abbrev-commit --decorate --date=relative --all" +
		"[git log graph]0\n"

	result := toHistLine(e)

	if result != expected {
		t.Fail()
	}
}

func TestParseToInvocation(t *testing.T) {
	ti := time.Time{}
	sync := strconv.QuoteRune(Syncd)

	sample := sync + "" + ti.Format(time.UnixDate) +
		"sorenlaptopbash/home/soren/src/project" +
		"git log --graph --abbrev-commit --decorate --date=relative --all" +
		"[git log graph]0\n"

	expect := IndexEntry{}
	expect.User = "soren"
	expect.Host = "laptop"
	expect.Shell = "bash"
	expect.Timestamp = ti
	expect.Directory = "/home/soren/src/project"
	expect.Command = "git log --graph --abbrev-commit --decorate --date=relative --all"
	expect.Tags = []string{"git", "log", "graph"}
	expect.Status = 0
	expect.HasStatus = true
	expect.IsSynced = true

	res, err := parseToEntry(sample)
	if err != nil {
		t.Fail()
	}

	expectedStr := fmt.Sprint(expect)
	resultStr := fmt.Sprint(res)

	if expectedStr != resultStr {
		t.Fail()
	}

}
