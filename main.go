package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/mvndaai/go_wasm_tools/internal/jsontools"
)

type JSWrappable func(string) (string, error)

func main() {

	jsFuncs := []struct {
		name string
		f    JSWrappable
	}{
		{"bytesToString", bytesToString},
		{"escapeJSON", jsontools.Escape},
		{"unescapeJSON", jsontools.Unescape},
		{"compressJSON", jsontools.Compress},
		{"prettyJSON", jsontools.Pretty},
	}

	for _, jsF := range jsFuncs {
		js.Global().Set(jsF.name, JSWrapper(jsF.f))
		fmt.Printf("Function '%s' loaded from Go\n", jsF.name)
	}

	// Keep the program alive so functions can be run over and over
	<-make(chan bool)
}

type output struct {
	Error    interface{} `json:"error"`
	Response interface{} `json:"response"`
}

func JSWrapper(f JSWrappable) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		input := args[0].String()
		if input == "" {
			return "No input provided"
		}
		resp, err := f(input)
		out := output{Response: resp}
		if err != nil {
			out.Error = err.Error()
		}

		pretty, err := json.Marshal(out)
		if err != nil {
			return "Could not marshal json"
		}
		return string(pretty)
	})
}

func bytesToString(in string) (string, error) {
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
