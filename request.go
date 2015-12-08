package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func FlushRequest(user string, url string, index Index) (success bool) {
	unsynced, err := index.GetUnsynced()
	if err != nil {
		panic(fmt.Sprintf("[gohst] %s: %s\n", "Invalid Hist File Error: ", err))
	}

	payload := make(Invocations, len(unsynced))

	for i, v := range unsynced {
		payload[i] = v.ToInvocation()
	}
	query := url + "/api/users/" + user + "/commands"
	success = send(query, payload)
	return
}

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

func send(url string, invs Invocations) (success bool) {
	jsonStr, err := json.Marshal(invs)
	if err != nil {
		panic(fmt.Sprintf("[gohst] %s: %s\n", "JSON Encoding Error: ", err))
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("[gohst] %s: %s\n", "Connection Error: ", err))
	}
	if resp.StatusCode != http.StatusCreated {
		success = false
	} else {
		success = true
	}

	defer resp.Body.Close()
	return
}

func receive(url string, isJson bool) (commands []string) {
	resp, err := http.Get(url)

	if err != nil {
		panic(fmt.Sprintf("[gohst] %s: %s\n", "Connection Error: ", err))
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("[gohst] %s: %s\n", "Response Reading Error: ", err))
	}

	if isJson {
		jsonArray, err := unmarshalInvocations(contents)
		if err != nil {
			panic(fmt.Sprintf("[gohst] %s: %s\n", "Malformed Response Error: ", err))
		}
		for _, v := range jsonArray {
			commands = append(commands, fmt.Sprint(v))
		}
	} else {
		err = json.Unmarshal(contents, &commands)
		if err != nil {
			panic(fmt.Sprintf("[gohst] %s: %v\n", "Malformed Response Error: ", err))
		}
	}
	return
}

func unmarshalInvocations(content []byte) (result Invocations, err error) {
	err = json.Unmarshal(content, &result)
	if err != nil {
		panic(fmt.Sprintf("[gohst] %s: %s\n", "Malformed Response Error: ", err))
	}
	return
}
