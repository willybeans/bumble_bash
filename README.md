## A 2D arcade style game in Golang

This is a bullet hell type game of a bee trying to return pollen to its hive

### Compiling yourself

- https://ebitengine.org/en/documents/webassembly.html
- (compile) `env GOOS=js GOARCH=wasm go build -o bumble_bash.wasm github.com/willybeans/bumble_bash`
- (copy) `cp $(go env GOROOT)/misc/wasm/wasm_exec.js .`
- create an HTML file that will run the wasm file
- using iframe is recommended `<iframe src="main.html" width="640" height="480" allow="autoplay"></iframe>`
