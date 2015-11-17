package main

import (
	"net/http"
	"errors"
)

func Request(string user, string url, bool verbose) (err error) (){

	//gohst.herokuapp.com
	//?verbose=true

	var resp *http.Response
	var err error

	if verbose {
		resp, err = http.Get("gohst.herokuapp.com/user/" + user + "?verbose=true")	
	} else {
		resp, err = http.Get("gohst.herokuapp.com/user/" + user)
	}

	return
}