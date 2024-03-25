# Internal web

## Dev

**anko** has not been working as well as I hoped. I'm still doing a lot of `ctrl+c` and `go run .`.

```bash	
# generate templ files (requires templ)
templ generate -watch -path pages

# live reload (air)
air
# run for dev
go run . -debug -dev

# run as in production
go run .
```

## Build


```bash
# build for local environment
go build -o dpj-web . 
```
```bash
GOOS=linux GOARCH=amd64 go build -o build/dpj_linux && \
GOOS=darwin GOARCH=amd64 go build -o build/dpj_macos && \
GOOS=windows GOARCH=amd64 go build -o build/dpj_win.exe && echo "done"
```


## Prequisites

### Running
- open-ai api key at `OPENAI_API_KEY`

### Development
- open-ai api key at `OPENAI_API_KEY`
- `templ` for generating templates
- `anko` for live reloading (optional)


## Philosophy

- I am embedding all files needed for operation into a single binary. (This is done with the `go:embed` directive.) 
- I'm trying to keep the code fairly 
