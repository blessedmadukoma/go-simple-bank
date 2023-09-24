package gvalidator

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
)

// ValidateString validates the given input between the given min and max length.
func ValidateString(input string, min, max int) error {
	if len(input) < min || len(input) > max {
		return fmt.Errorf("input must be between %d and %d characters", min, max)
	}
	return nil
}

// ValidateUsername validates the given username.
func ValidateUsername(username string) error {
	if err := ValidateString(username, 3, 16); err != nil {
		return err
	}

	if !isValidUsername(username) {
		return fmt.Errorf("username must be between 3 and 16 characters and must contain only lowercase letters, numbers and underscores")
	}

	return nil
}

// ValidateFullName validates the given name.
func ValidateFullName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidFullName(value) {
		return fmt.Errorf("must contain only letters or spaces")
	}
	return nil
}

// ValidatePassword validates the given password.
func ValidatePassword(password string) error {
	return ValidateString(password, 6, 32)
}

// ValidateEmail validates the given email.
func ValidateEmail(email string) error {
	if err := ValidateString(email, 3, 200); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("not a valid email address")
	}

	return nil
}
