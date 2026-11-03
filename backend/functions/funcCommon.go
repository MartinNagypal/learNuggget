package functions

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"path/filepath"
	"runtime"

	"github.com/go-passwd/validator"
)

const (
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars     = "0123456789"
)

func GetFile() string {
	_, file, _, _ := runtime.Caller(1)
	return filepath.Base(file)
}

func LogError(file string, err error) {
	log.Printf("Error in %s :", file)
	log.Println(err)
}

func GenerateSessionToken() (string, error) {
	token := make([]byte, 32) // 256 Bit

	if _, err := rand.Read(token); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(token), nil
}

func GenerateHash(data string) string {
	sum := sha256.Sum256([]byte(data))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func ValidatePw(password string) error {
	var passwordValidator = validator.New(
		validator.MinLength(8, errors.New("password must be at least 8 characters long")),
		validator.MaxLength(32, errors.New("password must be at most 32 characters long")),
		validator.ContainsAtLeast(uppercaseChars, 1, errors.New("password must contain at least one uppercase letter")),
		validator.ContainsAtLeast(lowercaseChars, 1, errors.New("password must contain at least one lowercase letter")),
		validator.ContainsAtLeast(digitChars, 1, errors.New("password must contain at least one number")),
		validator.CommonPassword(errors.New("password is too common")),
	)
	return passwordValidator.Validate(password)
}
