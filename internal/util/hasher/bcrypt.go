package hasher

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type bcryptHash struct {
	cost int
}

func NewBcrypt(cost int) Hasher {
	return &bcryptHash{
		cost: cost,
	}
}

func (h bcryptHash) Hash(value string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(value), h.cost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (h bcryptHash) IsEqual(hashed string, rawValue string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(rawValue))

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
