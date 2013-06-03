package main

import (
	"testing"
)

type someType struct {
	S string
}

func TestStringStuff(t *testing.T) {
	st := &someType{}
	println(st.S)

	if st.S == nil {
		t.Error("...")
	}
}
