package securities_impl

import (
	"os"
	"strconv"

	"github.com/danzbraham/cats-social/internal/applications/securities"
	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordHasher struct{}

func NewBcryptPasswordHasher() securities.PasswordHasher {
	return &BcryptPasswordHasher{}
}

func (b *BcryptPasswordHasher) HashPassword(password string) (string, error) {
	salt, err := strconv.Atoi(os.Getenv("BCRYPT_SALT"))
	if err != nil {
		return "", err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), salt)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (b *BcryptPasswordHasher) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
