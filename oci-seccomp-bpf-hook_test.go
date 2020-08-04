package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	types "github.com/seccomp/containers-golang"
	"github.com/stretchr/testify/assert"
)

func TestParseAnnotation(t *testing.T) {
	testProfile := types.Seccomp{}
	testProfile.DefaultAction = types.ActErrno

	tmpFile, err := ioutil.TempFile(os.TempDir(), "input-*.json")
	if err != nil {
		t.Fatalf("cannot create temporary file")
	}
	defer os.Remove(tmpFile.Name())
	testProfileByte, err := json.Marshal(testProfile)
	if err != nil {
		t.Fatalf("cannot marshal json")
	}

	if _, err := tmpFile.Write(testProfileByte); err != nil {
		t.Fatalf("cannot write to the temporary file")
	}

	for _, c := range []struct {
		annotation, input, output string
	}{
		{"if:" + tmpFile.Name() + ";of:/home/test/output.json", tmpFile.Name(), "/home/test/output.json"},
		{"of:/home/test/output.json", "", "/home/test/output.json"},
		{"of:/home/test/output.json;if:" + tmpFile.Name(), tmpFile.Name(), "/home/test/output.json"},
	} {
		output, input, err := parseAnnotation(c.annotation)
		assert.Nil(t, err)
		assert.Equal(t, c.input, input)
		assert.Equal(t, c.output, output)
	}

	// test malformed annotations
	for _, c := range []string{
		"if:/home/test/input1.json;if:/home/test/input2.json;of:/home/test/output.json",
		"if:" + tmpFile.Name(),
		"if:input;of:/home/test/output.json",
		"if:" + tmpFile.Name() + ";of:output",
	} {
		_, _, err := parseAnnotation(c)
		assert.NotNil(t, err)
	}
}
