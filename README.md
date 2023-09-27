# Internal web

## Dev

**anko** has not been working as well as I hoped. I'm still doing a lot of `ctrl+c` and `go run .`.

```bash	
# generate templ files
templ -watch -path pages -cmd mv "pages/*_templ.go pages/generated"

# run server
go run .
```

## Prequisites

- open-ai api key at `OPENAI_API_KEY`


## Philosophy

- I am embedding all files needed for operation into a single binary. (This is done with the `go:embed` directive.) 
- I'm trying to keep the code fairly 
