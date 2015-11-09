package main

import (
	"fmt"
	"os"
)

const d rune = '\x01'

type localRepo struct {
	FilePath string
}

func (r *localRepo) write(e invocation) (out string, err error) {
	file, err := os.OpenFile(r.FilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		return "", err
	}

	var record string
	if e.Command != "" {
		record = fmt.Sprintf("%c%v%c%s%c%s%c%s%c%s%c%s%c",
			d, e.Timestamp,
			d, e.User,
			d, e.Host,
			d, e.Directory,
			d, e.Command,
			d, e.Tags, d)
	}

	if e.HasStatus {
		record = fmt.Sprintf("%s%d\n", record, e.Status)
	}

	file.WriteString(record)
	return record, nil
}
