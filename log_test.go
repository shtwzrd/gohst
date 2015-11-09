package main

import "testing"

func TestParseOutTags(t *testing.T) {
	s := "echo \"'hello \"world\", quotes \"\" are \"fun\"' #fun not"
	o, tags := parseOutTags(s)
	expected := "echo \"'hello \"world\", quotes \"\" are \"fun\"' "

	if o != expected {
		t.Fail()
	}
	if tags[0] != "fun" || tags[1] != "not" {
		t.Fail()
	}
}
