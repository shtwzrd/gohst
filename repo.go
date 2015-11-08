package main

type repo interface {
	write(e entry) error
}
