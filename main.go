package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/mvndaai/go_wasm_tools/internal/bytesconvert"
	"github.com/mvndaai/go_wasm_tools/internal/htmltools"
	"github.com/mvndaai/go_wasm_tools/internal/jsontools"
	"github.com/mvndaai/go_wasm_tools/internal/php"
	"github.com/mvndaai/go_wasm_tools/internal/uuid"
)

type JSWrappable func(string) (string, error)

func main() {

	jsFuncs := []struct {
		name string
		f    JSWrappable
	}{
		{"bytesToString", bytesconvert.ToString},
		{"stringToBytes", bytesconvert.FromString},
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
		{"genUUIDv7", uuid.GenerateUUIDv7},
		{"timestampUUIDv7", uuid.TimestampUUIDv7},
		{"phpSerializeEncode", php.Encode},
		{"phpSerializeDecode", php.Decode},
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
