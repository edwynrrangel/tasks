package password

import (
	"crypto/rand"
	"math/big"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/edwynrrangel/tasks/logger"
)

type Password interface {
	IsValid(userPassword, commingPassword string) bool
	Hash(password string) (string, error)
	GenerateRamdom(length uint8) (string, error)
}

type password struct {
}

func New() Password {
	return &password{}
}

func (p *password) IsValid(userPassword, commingPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(commingPassword)); err != nil {
		return false
	}
	return true
}

func (p *password) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(
			"error hashing password",
			"func", "hashPassword - bcrypt.GenerateFromPassword",
			"error", err,
		)
		return "", err
	}
	return string(hashedPassword), nil
}

func (p *password) GenerateRamdom(length uint8) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789" +
		"!@#$%^&*-_+=<>?"

	var password strings.Builder
	for i := uint8(0); i < length; i++ {
		randomInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password.WriteByte(charset[randomInt.Int64()])
	}
	return password.String(), nil
}
