package customvalidator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

func PasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	fmt.Println(password)
	// Define password criteria
	var (
		minLength    = 8
		upperRegex   = `[A-Z]`
		lowerRegex   = `[a-z]`
		digitRegex   = `[0-9]`
		specialRegex = `[!@#$%^&*()]`
	)

	// Check for minimum length
	if len(password) < minLength {
		return false
	}

	// Check for at least one uppercase letter
	if matched, err := regexp.MatchString(upperRegex, password); err != nil || !matched {
		return false
	}

	// Check for at least one lowercase letter
	if matched, err := regexp.MatchString(lowerRegex, password); err != nil || !matched {
		return false
	}

	// Check for at least one digit
	if matched, err := regexp.MatchString(digitRegex, password); err != nil || !matched {
		return false
	}

	// Check for at least one special character
	if matched, err := regexp.MatchString(specialRegex, password); err != nil || !matched {
		return false
	}

	return true
}
