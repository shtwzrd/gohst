package main

import (
	"fmt"
	"os"
	"time"
)

const d rune = '\x01'

type LocalRepo struct {
	FilePath string
}

func (r LocalRepo) Query() (result []Invocation) {
	return nil
}

func (r LocalRepo) Write(e Invocation) (err error) {
	file, err := os.OpenFile(r.FilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	record := toHistLine(e)

	file.WriteString(record)
	return nil
}

func toHistLine(e Invocation) (record string) {
	if e.Command != "" {
		record = fmt.Sprintf("%c%v%c%s%c%s%c%s%c%s%c%s%c",
			d, e.Timestamp.Format(time.UnixDate),
			d, e.User,
			d, e.Host,
			d, e.Directory,
			d, e.Command,
			d, e.Tags, d)
	}

	if e.HasStatus {
		record = fmt.Sprintf("%s%d\n", record, e.Status)
	}
	return
}
