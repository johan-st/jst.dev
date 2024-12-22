#!/bin/bash

set -e

templ generate
# ./tailwindcss -i src/web/assets/css/input.css -o src/web/assets/css/output.css
go build -o main src/main.go
./main