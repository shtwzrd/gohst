package main

type repo interface {
	write(e invocation) error
}
