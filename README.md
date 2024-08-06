# go_wasm_tools
This is a page for tools for go running in wasm

## Building

To build the wasm file run:

```bash
GOOS=js GOARCH=wasm go build -o main.wasm
```


## Notes

### wasm_exec.js

In order use wasm in Go you need Go's `wasm_exec.js` file.

This repo is uses [jsDeliver](https://www.jsdelivr.com/) as a CDN from the [Go repo](https://github.com/golang/go/blob/master/misc/wasm/wasm_exec.js) to get that file with this HTML header:
```html
<script src="https://cdn.jsdelivr.net/gh/golang/go@go1.20.2/misc/wasm/wasm_exec.js"></script>
```

To copy the file from your computer's Go install to your current directory use this command:
```bash
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
```

### Making Go functions avaliable from Javascript

The way to expose Go functions to the Javascript is by importing `syscall/js` and calling `js.Global().Set` with a `js.FuncOf` wrapper around your Go function. The name given in `.Set` will be what is available in the Javascript.

### Long Running

Go by default ends execution as soon as the main function completes. If you want the functions to be available you will have to keep the main function alive. I do this by having a channel at the end of `main`:
```go
<-make(chan bool)
```


## References

https://golangbot.com/webassembly-using-go/

https://codepo8.github.io/css-fork-on-github-ribbon/