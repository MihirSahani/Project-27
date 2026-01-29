package internal

type Authenticator interface {
	GenerateToken(int64) (string, error)
	ValidateToken(string) (int64, error)
}