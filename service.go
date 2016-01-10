package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// HttpSuccess represents the highest HTTP status code that can still indicate
// a success; '226 IM Used' (RFC 3229)
const HttpSuccess = 226

type Service struct {
	Username string
	Password string
	Domain   string
}

func NewService(user string, pass string, domain string) Service {
	service := Service{}
	service.Username = user
	service.Password = pass
	service.Domain = domain
	return service
}

func (s Service) SendJson(route string, data interface{}) error {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		panic(fmt.Sprintf("[gohst] %s: %s\n", "JSON Encoding Error: ", err))
	}

	req, err := http.NewRequest("POST", s.Domain+route, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(s.Username, s.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("[gohst] %s: %s\n", "Error sending data: ", err))
	}

	fmt.Println(resp.StatusCode)
	if resp.StatusCode > HttpSuccess {
		return errors.New(fmt.Sprintf("Server responded with HTTP status code %d", resp.StatusCode))
	} else {
		err = resp.Body.Close()
		return err
	}
}

func (s Service) Receive(route string) (content []byte, err error) {
	req, err := http.NewRequest("GET", s.Domain+route, nil)
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(s.Username, s.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	fmt.Println(resp.StatusCode)
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
