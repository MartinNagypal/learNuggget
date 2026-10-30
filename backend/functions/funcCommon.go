package functions

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"path/filepath"
	"runtime"
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
