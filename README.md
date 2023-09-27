# Internal web

## Dev

**anko** has not been working as well as I hoped. I'm still doing a lot of `ctrl+c` and `go run .`.

```bash	
# generate templ files (requires templ)
templ generate -watch -path pages

# live reload (requires anko)
anko

# run server
go run . -v -dev
```

## Build

```bash
# build
go build -o dpj-internal-web .
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
