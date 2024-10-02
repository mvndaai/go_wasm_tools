package pickle

import (
	"fmt"
	"strings"

	"github.com/MacIt/pickle"
)

func Encode(s string) (string, error) {
	var sb strings.Builder
	err := pickle.NewEncoder(&sb).Encode(s)
	return sb.String(), err
}

func Decode(s string) (string, error) {
	obj, err := pickle.NewDecoder(strings.NewReader(s)).Decode()
	return fmt.Sprint(obj), err
}
