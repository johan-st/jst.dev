FROM golang:1.23-alpine AS build
RUN apk add --no-cache curl alpine-sdk

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate
    # templ generate && \
    # curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 -o tailwindcss && \
    # chmod +x tailwindcss && \
    # ./tailwindcss -i cmd/web/assets/css/input.css -o cmd/web/assets/css/output.css

RUN go test ./... -v

RUN CGO_ENABLED=0 GOOS=linux go build -o main src/main.go

FROM alpine:3.20.1 AS prod
WORKDIR /app
COPY --from=build /app/main /app/main
EXPOSE ${PORT}
CMD ["./main"]


