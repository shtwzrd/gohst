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

	client := &http.Client{}
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

	client := &http.Client{}
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
