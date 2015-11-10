package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const D rune = '\x01'
const Syncd rune = '✓'
const UnixDate string = "Mon Jan _2 15:04:05 MST 2006"

type Index struct {
	FilePath string
}

func (r Index) Write(e Invocation) (err error) {
	file, err := os.OpenFile(r.FilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	record := toHistLine(e)

	file.WriteString(record)
	return nil
}

func (r Index) GetUnsynced() (result []Invocation, err error) {
	file, err := os.OpenFile(r.FilePath, os.O_RDWR, 0644)
	if err != nil {
		return
	}
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		for index, runeval := range scanner.Text() {
			if index > 0 {
				break
			}
			if runeval == Syncd {
				invocation, parseErr := parseToInvocation(scanner.Text())
				if parseErr != nil {
					err = parseErr
					return
				}
				result = append(result, invocation)
			}
		}
	}
	return
}

func parseToInvocation(line string) (e Invocation, err error) {
	tokens := strings.FieldsFunc(line, func(c rune) bool {
		return c == D
	})

	e = Invocation{}
	if tokens[0] == strconv.QuoteRune(Syncd) {
		e.IsSynced = true
	} else {
		e.IsSynced = false
	}

	e.Timestamp, err = time.Parse(UnixDate, tokens[1])
	if err != nil {
		return
	}
	e.User = tokens[2]
	e.Host = tokens[3]
	e.Shell = tokens[4]
	e.Directory = tokens[5]
	e.Command = tokens[6]
	e.Tags = strings.Fields(tokens[7][1 : len(tokens[7])-1])
	exitcode, err := strconv.Atoi(tokens[8])
	if err != nil {
		return
	}
	e.Status = int8(exitcode)
	e.HasStatus = true
	return
}

func (r Index) Sync() (result []Invocation, err error) {
	return nil, nil
}

func toHistLine(e Invocation) (record string) {
	if e.Command != "" {
		record = fmt.Sprintf("%c%v%c%s%c%s%c%s%c%s%c%s%c%s%c",
			D, e.Timestamp.Format(time.UnixDate),
			D, e.User,
			D, e.Host,
			D, e.Shell,
			D, e.Directory,
			D, e.Command,
			D, e.Tags, D)
	}

	if e.HasStatus {
		record = fmt.Sprintf("%s%d%c\n", record, e.Status, D)
	}
	return
}
