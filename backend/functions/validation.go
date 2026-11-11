package functions

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Size limits. Data above these limits is rejected.
const (
	MaxNameLength         = 100  // collection and group names
	MaxCardFieldLength    = 5000 // question and answer, each
	MaxCardsPerCollection = 2000
	MaxCardsPerImport     = 500 // cards per /addBulk request

	UsernameMinLength = 5
	UsernameMaxLength = 32
)

// Validation errors. Their messages are meant for the user, so handlers can
// send them to the client with status 400.
var (
	ErrUsernameLength = fmt.Errorf("Username must be %d-%d characters long", UsernameMinLength, UsernameMaxLength)
	ErrUsernameChars  = errors.New("Username may only contain letters (a-z, A-Z), digits and _ . -")
	ErrNameEmpty      = errors.New("Name is empty")
	ErrNameTooLong    = fmt.Errorf("Name must be at most %d characters long", MaxNameLength)
	ErrCardEmpty      = errors.New("Question and answer must not be empty")
	ErrCardTooLong    = fmt.Errorf("Question and answer must be at most %d characters long", MaxCardFieldLength)
	ErrTooManyCards   = fmt.Errorf("A collection can hold at most %d cards", MaxCardsPerCollection)
)

// Only ASCII letters, digits and "_.-": no whitespace, no invisible or
// look-alike Unicode characters (e.g. Cyrillic "а" instead of Latin "a").
var usernamePattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)

// NormalizeUsername trims the username and checks the registration rules.
func NormalizeUsername(username string) (string, error) {
	username = strings.TrimSpace(username)
	length := utf8.RuneCountInString(username)
	if length < UsernameMinLength || length > UsernameMaxLength {
		return "", ErrUsernameLength
	}
	if !usernamePattern.MatchString(username) {
		return "", ErrUsernameChars
	}
	return username, nil
}

// NormalizeName trims a collection or group name and checks its length.
func NormalizeName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", ErrNameEmpty
	}
	if utf8.RuneCountInString(name) > MaxNameLength {
		return "", ErrNameTooLong
	}
	return name, nil
}

// ValidateCard checks that question and answer are present and not too long.
func ValidateCard(question, answer string) error {
	if strings.TrimSpace(question) == "" || strings.TrimSpace(answer) == "" {
		return ErrCardEmpty
	}
	if utf8.RuneCountInString(question) > MaxCardFieldLength || utf8.RuneCountInString(answer) > MaxCardFieldLength {
		return ErrCardTooLong
	}
	return nil
}
