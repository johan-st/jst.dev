package main

import "github.com/johan-st/jst.dev/auth"

func main() {
	_ = auth.NewSingleUserStore(&auth.User{Email: "jst@jst.dev", PasswordHash: []byte("password")})
}
