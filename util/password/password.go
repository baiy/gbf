package password

import (
	"log/slog"

	"github.com/alexedwards/argon2id"
)

func Generate(value string) string {
	hash, err := argon2id.CreateHash(value, argon2id.DefaultParams)
	if err != nil {
		slog.Error(err.Error())
	}
	return hash
}

func Verify(value, password string) bool {
	is, err := argon2id.ComparePasswordAndHash(value, password)
	if err != nil {
		slog.Error(err.Error())
		return false
	}
	return is
}
