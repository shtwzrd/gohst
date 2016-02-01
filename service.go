package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	g "github.com/warreq/gohstd/common"
	"io/ioutil"
	"net/http"
)

// HttpSuccess represents the highest HTTP status code that can still indicate
// a success; '226 IM Used' (RFC 3229)
const HttpSuccess = 226

type Service struct {
	Password string
	Url      string
}

func NewService(url string, password string) Service {
	service := Service{}
	service.Password = password
	service.Url = url
	return service
}

func (s Service) SendJson(user, route string, data interface{}) error {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		panic(fmt.Sprintf("[gohst] %s: %s\n", "JSON Encoding Error: ", err))
	}

	req, err := http.NewRequest("POST", s.Url+route, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(user, s.Password)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("[gohst] %s: %s\n", "Error sending data: ", err))
	}

	if resp.StatusCode > HttpSuccess {
		return errors.New(fmt.Sprintf("Server responded with HTTP status code %d", resp.StatusCode))
	} else {
		err = resp.Body.Close()
		return err
	}
}

func (s Service) Receive(user, route string) (content []byte, err error) {
	req, err := http.NewRequest("GET", s.Url+route, nil)
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(user, s.Password)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > HttpSuccess {
		return nil, errors.New(fmt.Sprintf("Server responded with HTTP status code %d", resp.StatusCode))
	}

	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	return
}

func InvocationToBase64(invoc g.Invocation) g.Invocation {
	enc := base64.StdEncoding
	invoc.Command = enc.EncodeToString([]byte(invoc.Command))
	invoc.User = enc.EncodeToString([]byte(invoc.User))
	invoc.Host = enc.EncodeToString([]byte(invoc.Host))
	invoc.Shell = enc.EncodeToString([]byte(invoc.Shell))
	invoc.Directory = enc.EncodeToString([]byte(invoc.Directory))
	for i, t := range invoc.Tags {
		invoc.Tags[i] = enc.EncodeToString([]byte(t))
	}
	return invoc
}

func InvocationsToBase64(invocs g.Invocations) g.Invocations {
	for i, v := range invocs {
		invocs[i] = InvocationToBase64(v)
	}
	return invocs
}

func InvocationFromBase64(invoc g.Invocation) (g.Invocation, error) {
	enc := base64.StdEncoding
	var err error
	cmd, err := enc.DecodeString(invoc.Command)
	user, err := enc.DecodeString(invoc.User)
	host, err := enc.DecodeString(invoc.Host)
	shell, err := enc.DecodeString(invoc.Shell)
	dir, err := enc.DecodeString(invoc.Directory)
	tags := make([][]byte, len(invoc.Tags))
	for i, t := range invoc.Tags {
		tags[i], err = enc.DecodeString(t)
	}
	if err != nil {
		return invoc, err
	}
	invoc.Command = string(cmd)
	invoc.User = string(user)
	invoc.Host = string(host)
	invoc.Shell = string(shell)
	invoc.Directory = string(dir)
	for i, t := range tags {
		invoc.Tags[i] = string(t)
	}
	return invoc, nil
}

func InvocationsFromBase64(invocs g.Invocations) (g.Invocations, error) {
	var err error
	for i, v := range invocs {
		invocs[i], err = InvocationFromBase64(v)
	}
	return invocs, err
}

func CommandFromBase64(cmd g.Command) (g.Command, error) {
	enc := base64.StdEncoding
	c, err := enc.DecodeString(string(cmd))
	if err != nil {
		return cmd, nil
	}
	return g.Command(c), nil
}

func CommandsFromBase64(cmds g.Commands) (g.Commands, error) {
	var err error
	for i, v := range cmds {
		cmds[i], err = CommandFromBase64(v)
	}
	return cmds, err
}
