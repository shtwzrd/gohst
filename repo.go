package main

type Repo interface {
	Write(e Invocation) error
	Query() []Invocation
}
