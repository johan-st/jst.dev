package auth

type User struct {
	Email        string
	PasswordHash []byte
}
