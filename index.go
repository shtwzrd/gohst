package main

import (
	"bufio"
	"errors"
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

func (r Index) Write(e IndexEntry) (err error) {
	file, err := os.OpenFile(r.FilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	defer file.Close()

	if err != nil {
		return err
	}

	record := toHistLine(e)

	file.WriteString(record)
	return nil
}

func (r Index) MarkSynced() (err error) {
	file, err := os.OpenFile(r.FilePath, os.O_RDWR, 0644)
	if err != nil {
		return
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		return err
	}

	file.Close()
	os.Remove(r.FilePath)
	file, err = os.Create(r.FilePath)
	if err != nil {
		return
	}

	w := bufio.NewWriter(file)
	for _, line := range lines {
		for index, runeval := range line {
			if index > 0 {
				break
			}
			if runeval != Syncd {
				fmt.Fprintf(w, "%c%s\n", Syncd, line[1:len(line)-1])
			} else {
				fmt.Fprintf(w, "%s\n", line)
			}
		}
	}
	err = w.Flush()
	file.Close()
	return
}

func (r Index) GetUnsynced() (result []IndexEntry, err error) {
	file, err := os.OpenFile(r.FilePath, os.O_RDWR, 0644)
	if err != nil {
		return
	}

	state, err := file.Stat()
	if err != nil {
		return
	}

	// exit if file is effectively empty
	if state.Size() < 8 {
		return nil, nil
	}

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		for index, runeval := range scanner.Text() {
			if index > 0 {
				break
			}
			if runeval != Syncd {
				line := scanner.Text()
				invocation, parseErr := parseToEntry(line)
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

func (r Index) Flush() {
	file, err := os.OpenFile(r.FilePath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return
	}

	defer file.Close()

	if _, err = file.WriteString("\n"); err != nil {
		panic(err)
	}
}

func (r Index) getLastLine() (last string) {
	file, err := os.OpenFile(r.FilePath, os.O_RDONLY, 0644)
	if err != nil {
		return
	}

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		last = scanner.Text()
	}
	return last
}

func (r Index) canWriteResult() (valid bool) {
	line := r.getLastLine()
	for index, runeval := range line {
		if index > 0 {
			break
		}
		if runeval == Syncd || runeval == 'U' {
			valid = true
		} else {
			valid = false
		}
	}

	if !valid {
		return
	}

	return !r.canWriteContext()
}

func (r Index) canWriteContext() (valid bool) {
	line := r.getLastLine()
	if line == "" {
		return true
	}

	for _, runeval := range line {
		if runeval != D {
			valid = true
		} else {
			valid = false
		}
	}
	return valid
}

func parseToEntry(line string) (e IndexEntry, err error) {
	tokens := strings.FieldsFunc(line, func(c rune) bool {
		return c == D
	})

	e = IndexEntry{}
	if tokens[0] == strconv.QuoteRune(Syncd) {
		e.IsSynced = true
	} else if tokens[0] == "U" {
		e.IsSynced = false
	} else {
		errstr := "Your gohst index file appears malformed with line:\n" + line
		return IndexEntry{}, errors.New(errstr)
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
	// if we have a command with no exit code, it means the shell probably
	// exited between logging context and logging exit code. This could happen
	// if the command was, for example, 'exit'. We assume it returned 0
	var exitcode int
	if len(tokens) < 9 {
		exitcode = 0
	} else {
		exitcode, err = strconv.Atoi(tokens[8])
		if err != nil {
			return
		}
	}
	e.Status = int8(exitcode)
	e.HasStatus = true
	return
}

func toHistLine(e IndexEntry) (record string) {
	if e.Command != "" {
		if e.Host == "" {
			e.Host = "null"
		}
		if e.User == "" {
			e.User = "null"
		}
		if e.Shell == "" {
			e.Shell = "null"
		}
		if e.Directory == "" {
			e.Directory = "null"
		}
		record = fmt.Sprintf("%c%c%v%c%s%c%s%c%s%c%s%c%s%c%s%c", 'U',
			D, e.Timestamp.Format(time.UnixDate),
			D, e.User,
			D, e.Host,
			D, e.Shell,
			D, e.Directory,
			D, e.Command,
			D, e.Tags, D)
	}

	if e.HasStatus {
		record = fmt.Sprintf("%s%d%c", record, e.Status, D)
	}
	return
}
