package bytesconvert

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

func ToString(in string) (string, error) {
	if strings.Contains(in, "{0x") {
		return bytesFromGolangFormat(in)
	}

	// Remove characters that are not digits
	if strings.Contains(in, "[") {
		in = strings.Split(in, "[")[1]
	}
	if strings.Contains(in, "]") {
		in = strings.Split(in, "]")[0]
	}
	in = strings.ReplaceAll(in, "\n", "")
	in = strings.TrimSpace(in)

	// Split the string into an array of digits
	in = strings.ReplaceAll(in, "0x", "")
	in = strings.ReplaceAll(in, " ", ",")
	in = strings.ReplaceAll(in, ",,", ",")
	parts := strings.Split(in, ",")

	var out []byte
	for _, p := range parts {
		i, err := strconv.Atoi(p)
		if err != nil {
			return "", fmt.Errorf("could not convert into int to use as byte: %w", err)
		}
		out = append(out, byte(i))
	}

	return string(out), nil
}

func bytesFromGolangFormat(in string) (string, error) {
	// ex: []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64}

	if strings.Contains(in, "{") {
		in = strings.Split(in, "{")[1]
	}
	if strings.Contains(in, "}") {
		in = strings.Split(in, "}")[0]
	}
	in = strings.ReplaceAll(in, "0x", "")
	in = strings.ReplaceAll(in, ",", "")
	in = strings.ReplaceAll(in, " ", "")

	decoded, err := hex.DecodeString(in)
	if err != nil {
		return "", fmt.Errorf("could not decode hex string: %w", err)
	}
	return string(decoded), nil
}

func FromString(s string) (string, error) {
	return fmt.Sprint([]byte(s)), nil
}
