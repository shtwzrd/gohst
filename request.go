package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetRequest(user string, url string, verbose bool, count int) (result []string) {

	query := url + "/api/users/" + user + "/commands"
	if verbose {
		query += "?verbose=true"
	} else {
		query += "?verbose=false"
	}
	if count > 0 {
		query += fmt.Sprintf("%s%d", "&count=", count)
	}
	result = receive(query, verbose)
	return
}

func receive(url string, isJson bool) (commands []string) {
	resp, err := http.Get(url)

	if err != nil {
		panic(fmt.Sprintf("[gohst] %s: %s", "Connection Error: ", err))
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("[gohst] %s: %s", "Response Reading Error: ", err))
	}

	if isJson {
		jsonArray, err := marshalInvocations(contents)
		if err != nil {
			panic(fmt.Sprintf("[gohst] %s: %s", "Malformed Response Error: ", err))
		}
		for _, v := range jsonArray {
			commands = append(commands, fmt.Sprint(v))
		}
	} else {
		err = json.Unmarshal(contents, &commands)
		if err != nil {
			panic(fmt.Sprintf("[gohst] %s: %v", "Malformed Response Error: ", err))
		}
	}
	return
}

func marshalInvocations(content []byte) (result Invocations, err error) {
	err = json.Unmarshal(content, &result)
	if err != nil {
		panic(fmt.Sprintf("[gohst] %s: %s", "Malformed Response Error: ", err))
	}
	return
}
