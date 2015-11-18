package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetRequest(user string, url string, verbose bool) (result []string) {

	//gohst.herokuapp.com
	//?verbose=true
	if verbose {
		result = receive(url + "/api/users/" + user + "/commands" + "?verbose=true")
	} else {
		result = receive(url + "/api/users/" + user + "/commands")
	}

	return
}

func receive(url string) (commands []string) {
	resp, err := http.Get(url)

	if err != nil {
		panic(fmt.Sprintf("[gohst] %s: %s", "Connection Error: ", err))
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("[gohst] %s: %s", "Response Reading Error: ", err))
	}

	err = json.Unmarshal(contents, &commands)
	if err != nil {
		panic(fmt.Sprintf("[gohst] %s: %s", "Malformed Response Error: ", err))
	}

	return
}
