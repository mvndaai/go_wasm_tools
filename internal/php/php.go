package php

import (
	"github.com/elliotchance/phpserialize"
)

func Encode(s string) (string, error) {
	r, err := phpserialize.Marshal(s, nil)
	return string(r), err
}

func Decode(s string) (string, error) {
	var r string
	err := phpserialize.Unmarshal([]byte(s), &r)
	return r, err
}
