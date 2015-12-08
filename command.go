package main

import (
	"time"
)

type Command string
type Commands []Command

type Invocation struct {
	Id        int       `json:"id"`
	SessionId int       `json:"sessionid"`
	Command   string    `json:"command"`
	Directory string    `json:"directory"`
	User      string    `json:"user"`
	Host      string    `json:"host"`
	Shell     string    `json:"shell"`
	ExitCode  int8      `json:",exitcode"`
	Timestamp time.Time `json:"timestamp"`
	Tags      []string  `json:"tags"`
}

type Invocations []Invocation
