package main

import (
	"fmt"
	"os"
)

const d rune = '\x01'

type localRepo struct {
	FilePath string
}

func (r *localRepo) write(e entry) error {
	file, err := os.OpenFile(r.FilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	var record string
	if e.Invocation != "" {
		record = fmt.Sprintf("%c%v%c%s%c%s%c%s%c%s%c", d, e.Timestamp, d, e.User, d, e.Host, d, e.Invocation, d, e.Directory, d)
	}

	if e.HasStatus {
		record = fmt.Sprintf("%s%d\n", record, e.Status)
	}

	file.WriteString(record)
	return nil
}
