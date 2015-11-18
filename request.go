package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Request(user string, url string, verbose bool) (result []string, err error) {

	//gohst.herokuapp.com
	//?verbose=true

	var message = make([]string, 10)

	if verbose {
		requestCommand(url + "/api/users/" + user + "/commands" + "?verbose=true")
	} else {
		message, _ = requestCommand(url + "/api/users/" + user + "/commands")
	}

	for _, s := range message {
		fmt.Println(s)
	}

	return
}

func requestCommand(url string) (commands []string, err error) {
	resp, err := http.Get(url)

	if err != nil {
		panic("Could not reach web service with given URL. Check the connection")
	}

	var message = make([]string, 10)

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("Could not read from response body")
	}

	err = json.Unmarshal(contents, &message)

	return message, err
}
