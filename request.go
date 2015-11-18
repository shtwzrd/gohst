package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetRequest(user string, url string, verbose bool) (result []string) {

	if verbose {
		result = receive(url+"/api/users/"+user+"/commands"+"?verbose=true", true)
	} else {
		result = receive(url+"/api/users/"+user+"/commands", false)
	}

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
			panic(fmt.Sprintf("[gohst] %s: %s", "Malformed Response Error: ", err))
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
