package main

import (
	"time"
)

type entry struct {
	Directory  string
	Shell      string
	User       string
	Host       string
	Invocation string
	Status     int8
	Commands   []string
	Tags       []string
	Timestamp  time.Time
	HasStatus  bool
	IsSynced   bool
}
