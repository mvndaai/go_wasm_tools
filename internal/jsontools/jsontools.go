package jsontools

import (
	"bytes"
	"encoding/json"
)

func Escape(s string) (string, error) {
	b, err := json.Marshal(s)
	return string(b), err
}

func Unescape(s string) (string, error) {
	var u string
	err := json.Unmarshal([]byte(s), &u)
	return u, err
}

func Pretty(s string) (string, error) {
	var b bytes.Buffer
	err := json.Indent(&b, []byte(s), "", "\t")
	return b.String(), err
}

func Compress(s string) (string, error) {
	var b bytes.Buffer
	err := json.Compact(&b, []byte(s))
	return b.String(), err
}
