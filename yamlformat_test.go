package main

import (
	"testing"

	"bytes"

	. "github.com/bborbe/assert"
)

func TestDoEmpty(t *testing.T) {
	writer := bytes.NewBufferString("")
	err := do(writer, "", false)
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(writer.String(), NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(writer.String()) > 0, Is(true))
	if err != nil {
		t.Fatal(err)
	}
}
