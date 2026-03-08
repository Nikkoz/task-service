package service

//go:generate mockery --name TokenManager --output ./mocks --outpkg mocks
type TokenManager interface {
	Generate(userID uint64) (string, error)
}
