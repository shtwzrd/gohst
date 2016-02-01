package main

import (
	"encoding/json"
	"errors"
	"fmt"
	g "github.com/warreq/gohstd/common"
)

type HttpCommandRepo struct {
	http Service
	key  []byte
}

func NewHttpCommandRepo(http Service, key []byte) HttpCommandRepo {
	repo := HttpCommandRepo{}
	repo.http = http
	repo.key = key
	return repo
}

func (h HttpCommandRepo) InsertInvocations(user string, invs g.Invocations) error {
	// Send the invocations unencrypted if no key was provided
	if len(h.key) > 0 {
		encrypted, err := EncryptInvocations(invs, h.key)
		if err != nil {
			return errors.New(fmt.Sprintf("Error: Could not encrypt commands: %s", err))
		}
		encoded := InvocationsToBase64(encrypted)
		fmt.Println(encoded)
		return h.http.SendJson(user, fmt.Sprintf("/api/users/%s/commands", user), encoded)
	}
	return h.http.SendJson(user, fmt.Sprintf("/api/users/%s/commands", user), invs)
}

func (h HttpCommandRepo) GetInvocations(user string, n int) (g.Invocations, error) {
	route := fmt.Sprintf("/api/users/%s/commands?verbose=true&count=%d", user, n)
	content, err := h.http.Receive(user, route)
	if err != nil {
		return nil, errors.New((fmt.Sprintf("[gohst] %s: %s\n", "Connection Error: ", err)))
	}

	var result g.Invocations
	err = json.Unmarshal(content, &result)
	if err != nil {
		err = errors.New((fmt.Sprintf("[gohst] %s: %s\n", "Malformed Response Error: ", err)))
	}

	if len(h.key) > 0 {
		decoded, err := InvocationsFromBase64(result)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error: Could not decode commands: %s", err))
		}
		decrypted, err := DecryptInvocations(decoded, h.key)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error: Could not decrypt commands: %s", err))
		}
		return decrypted, err
	}
	// Return invocations without decryption if no key was provided
	return result, err
}

func (h HttpCommandRepo) GetCommands(user string, n int) (g.Commands, error) {
	route := fmt.Sprintf("/api/users/%s/commands?verbose=false&count=%d", user, n)
	content, err := h.http.Receive(user, route)

	if err != nil {
		return nil, errors.New((fmt.Sprintf("[gohst] %s: %s\n", "Connection Error: ", err)))
	}

	var result g.Commands
	err = json.Unmarshal(content, &result)
	if err != nil {
		err = errors.New((fmt.Sprintf("[gohst] %s: %s\n", "Malformed Response Error: ", err)))
	}

	if len(h.key) > 0 {
		decoded, err := CommandsFromBase64(result)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error: Could not decode commands: %s", err))
		}
		fmt.Println(decoded)
		decrypted, err := DecryptCommands(decoded, h.key)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error: Could not decrypt commands: %s", err))
		}
		fmt.Println(decrypted)
		return decrypted, err
	}
	// Return commands without decryption if no key was provided
	return result, err
}
