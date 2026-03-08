package service

//go:generate mockery --name PasswordHasher --output ./mocks --outpkg mocks
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}
