//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/mvndaai/go_wasm_tools/internal/bytesconvert"
	"github.com/mvndaai/go_wasm_tools/internal/dadjoke"
	"github.com/mvndaai/go_wasm_tools/internal/htmltools"
	"github.com/mvndaai/go_wasm_tools/internal/jsontools"
	"github.com/mvndaai/go_wasm_tools/internal/php"
	"github.com/mvndaai/go_wasm_tools/internal/uuid"
)

type JSWrappable func(string) (string, error)

func main() {

	jsFuncs := []struct {
		name  string
		f     JSWrappable
		async bool
	}{
		{"bytesToString", bytesconvert.ToString, false},
		{"stringToBytes", bytesconvert.FromString, false},
		{"escapeJSON", jsontools.Escape, false},
		{"unescapeJSON", jsontools.Unescape, false},
		{"compressJSON", jsontools.Compress, false},
		{"prettyJSON", jsontools.Pretty, false},
		{"escapeHTML", htmltools.Escape, false},
		{"unescapeHTML", htmltools.Unescape, false},
		{"bEncode", htmltools.B64Encode, false},
		{"bDecode", htmltools.B64Decode, false},
		{"urlEncode", htmltools.URLEncode, false},
		{"urlDecode", htmltools.URLDecode, false},
		{"genUUIDv7", uuid.GenerateUUIDv7, false},
		{"timestampUUIDv7", uuid.TimestampUUIDv7, false},
		{"phpSerializeEncode", php.Encode, false},
		{"phpSerializeDecode", php.Decode, false},
		{"getDadJoke", dadjoke.GetJoke, true},
	}

	var loadedFuncs []string
	var m = make(map[string]any)
	for _, jsF := range jsFuncs {
		if jsF.async {
			m[jsF.name] = JSWrapperAsync(jsF.f)
		} else {
			m[jsF.name] = JSWrapper(jsF.f)
		}
		loadedFuncs = append(loadedFuncs, jsF.name)
	}
	js.Global().Set("golang", m)
	fmt.Println("Functions loaded into 'golang' object", loadedFuncs)

	// Keep the program alive so functions can be run over and over
	<-make(chan bool)
}

type output struct {
	Error    any `json:"error"`
	Response any `json:"response"`
}

func JSWrapper(f JSWrappable) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		input := ""
		if len(args) > 1 {
			return "Invalid number of arguments passed"
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

func JSWrapperAsync(f JSWrappable) js.Func { // If you have an http request it needs to return a promise
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		input := ""
		if len(args) > 1 {
			return js.Global().Get("Promise").Call("reject", "Invalid number of arguments passed")
		}
		if len(args) == 1 {
			input = args[0].String()
		}

		handler := js.FuncOf(func(this js.Value, args []js.Value) any {
			resolve := args[0]
			reject := args[1]

			go func() {
				resp, err := f(input)
				out := output{Response: resp}
				if err != nil {
					out.Error = err.Error()
				}

				pretty, err := json.Marshal(out)
				if err != nil {
					reject.Invoke("Could not marshal json")
					return
				}
				resolve.Invoke(string(pretty))
			}()

			return nil
		})

		return js.Global().Get("Promise").New(handler)
	})
}
