package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"syscall/js"
	"time"

	uuid "github.com/gofrs/uuid/v5"
	"github.com/mvndaai/go_wasm_tools/internal/htmltools"
	"github.com/mvndaai/go_wasm_tools/internal/jsontools"
	"github.com/mvndaai/go_wasm_tools/internal/pickle"
)

type JSWrappable func(string) (string, error)

func main() {

	jsFuncs := []struct {
		name string
		f    JSWrappable
	}{
		{"bytesToString", bytesToString},
		{"stringToBytes", stringToBytes},
		{"escapeJSON", jsontools.Escape},
		{"unescapeJSON", jsontools.Unescape},
		{"compressJSON", jsontools.Compress},
		{"prettyJSON", jsontools.Pretty},
		{"escapeHTML", htmltools.Escape},
		{"unescapeHTML", htmltools.Unescape},
		{"bEncode", htmltools.B64Encode},
		{"bDecode", htmltools.B64Decode},
		{"urlEncode", htmltools.URLEncode},
		{"urlDecode", htmltools.URLDecode},
		{"genUUIDv7", GenerateUUIDv7},
		{"pickleEncode", pickle.Encode},
		{"pickleDecode", pickle.Decode},
	}

	var loadedFuncs []string
	var m = make(map[string]any)
	for _, jsF := range jsFuncs {
		m[jsF.name] = JSWrapper(jsF.f)
		loadedFuncs = append(loadedFuncs, jsF.name)
	}
	js.Global().Set("golang", m)
	fmt.Println("Functions loaded into 'golang' object", loadedFuncs)

	// Keep the program alive so functions can be run over and over
	<-make(chan bool)
}

type output struct {
	Error    interface{} `json:"error"`
	Response interface{} `json:"response"`
}

func JSWrapper(f JSWrappable) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		input := ""
		if len(args) > 1 {
			return "Invalid no of arguments passed"
		}
		if len(args) == 1 {
			input = args[0].String()
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

func stringToBytes(s string) (string, error) {
	return fmt.Sprint([]byte(s)), nil
}

func GenerateUUIDv7(s string) (string, error) {
	if s == "" {
		guid, err := uuid.NewV7()
		if err != nil {
			return "", err
		}
		return guid.String(), nil
	}

	atTime, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return "", fmt.Errorf("could not parse time: %w", err)
	}
	guid, err := uuid.NewV7AtTime(atTime)
	if err != nil {
		return "", err
	}
	return guid.String(), nil
}
