package hash

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	strCost = 10
)

// Encrypt hashes a string
func Encrypt(str string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), strCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// CompareHashAndRaw compares hashed and raw string and
// returns error if they are not equal
func CompareHashAndRaw(rawStr, hashedStr string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedStr), []byte(rawStr)); err != nil {
		return err
	}

	return nil
}
