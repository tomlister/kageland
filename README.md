# Kageland

## Develop
`vercel dev`

### Building viewer.wasm (in lieu of build system)
`env GOOS=js GOARCH=wasm go build -o viewer.wasm .`