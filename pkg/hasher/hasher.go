package hasher

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type (
	Hasher struct{}

	Interactor interface {
		HashPassword(string) (string, error)
		CompareAndHash(string, string) bool
	}
)

func (h *Hasher) HashPassword(password string) (string, error) {
	var passwordBytes = []byte(password)

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)

	return string(hashedPasswordBytes), err
}

func (h *Hasher) CompareAndHash(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currPassword))
	fmt.Println(err)
	return err == nil
}

var _ Interactor = (*Hasher)(nil)

func NewHasher() *Hasher {
	return &Hasher{}
}
