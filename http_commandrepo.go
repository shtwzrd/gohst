package main

import (
	"encoding/json"
	"errors"
	"fmt"
	g "github.com/warreq/gohstd/common"
)

type HttpCommandRepo struct {
	http Service
}

func (h HttpCommandRepo) InsertInvocations(user string, invs g.Invocations) error {
	return h.http.SendJson(fmt.Sprintf("/api/users/%s/commands", user), invs)
}

func (h HttpCommandRepo) GetInvocations(user string, n int) (g.Invocations, error) {
	route := fmt.Sprintf("/api/users/%s/commands?verbose=true&count=%d", user, n)
	content, err := h.http.Receive(route)
	if err != nil {
		return nil, errors.New((fmt.Sprintf("[gohst] %s: %s\n", "Connection Error: ", err)))
	}

	var result g.Invocations
	err = json.Unmarshal(content, &result)
	if err != nil {
		err = errors.New((fmt.Sprintf("[gohst] %s: %s\n", "Malformed Response Error: ", err)))
	}
	return result, err
}

func (h HttpCommandRepo) GetCommands(user string, n int) (g.Commands, error) {
	route := fmt.Sprintf("/api/users/%s/commands?verbose=false&count=%d", user, n)
	content, err := h.http.Receive(route)

	if err != nil {
		return nil, errors.New((fmt.Sprintf("[gohst] %s: %s\n", "Connection Error: ", err)))
	}

	fmt.Println(content)
	var result g.Commands
	err = json.Unmarshal(content, &result)
	if err != nil {
		err = errors.New((fmt.Sprintf("[gohst] %s: %s\n", "Malformed Response Error: ", err)))
	}
	return result, err
}
