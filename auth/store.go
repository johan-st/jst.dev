package auth

type Store interface {
	StoreUser(user *User) error
	FindUserByEmail(email string) (*User, error)
}

func NewSingleUserStore(user *User) Store {
	return &singleUserStore{user}
}

type singleUserStore struct {
	user *User
}

func (s *singleUserStore) StoreUser(user *User) error {
	return nil
}

func (s *singleUserStore) FindUserByEmail(email string) (*User, error) {
	return s.user, nil
}
