//nolint:revive,wrapcheck // This file is taken from Prysm
package app

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/logrusorgru/aurora"

	"github.com/piplabs/story/lib/errors"
)

const (
	// Constants for passwords.
	minPasswordLength = 8

	// NewKeyPasswordPromptText for key creation.
	NewKeyPasswordPromptText = "New key password"
	// PasswordPromptText for wallet unlocking.
	PasswordPromptText = "Key password"
	// ConfirmPasswordPromptText for confirming a key password.
	ConfirmPasswordPromptText = "Confirm password"
)

var (
	au = aurora.NewAurora(true)

	errPasswordWeak = errors.New("password must have at least 8 characters")
)

// PasswordReaderFunc takes in a *file and returns a password using the terminal package.
func passwordReaderFunc(file *os.File) ([]byte, error) {
	pass, err := terminal.ReadPassword(int(file.Fd()))

	return pass, err
}

// PasswordReader has passwordReaderFunc as the default but can be changed for testing purposes.
var PasswordReader = passwordReaderFunc

// PasswordPrompt prompts the user for a password, that repeatedly requests the password until it qualifies the
// passed in validation function.
func PasswordPrompt(promptText string, validateFunc func(string) error) (string, error) {
	var responseValid bool
	var response string
	for !responseValid {
		fmt.Printf("%s: ", au.Bold(promptText))
		bytePassword, err := PasswordReader(os.Stdin)
		if err != nil {
			return "", err
		}
		response = strings.TrimRight(string(bytePassword), "\r\n")
		if err := validateFunc(response); err != nil {
			fmt.Printf("\nEntry not valid: %s\n", au.BrightRed(err))
		} else {
			fmt.Println("")
			responseValid = true
		}
	}

	return response, nil
}

// InputPassword with a custom validator along capabilities of confirming
// the password and reading it from disk if a specified flag is set.
func InputPassword(
	promptText, confirmText string,
	shouldConfirmPassword bool,
	passwordValidator func(input string) error,
) (string, error) {
	if strings.Contains(strings.ToLower(promptText), "new wallet") {
		fmt.Println("Password requirements: at least 8 characters")
	}
	var hasValidPassword bool
	var password string
	var err error
	for !hasValidPassword {
		password, err = PasswordPrompt(promptText, passwordValidator)
		if err != nil {
			return "", errors.Wrap(err, "could not read password")
		}
		if shouldConfirmPassword {
			passwordConfirmation, err := PasswordPrompt(confirmText, passwordValidator)
			if err != nil {
				return "", errors.Wrap(err, "could not read password confirmation")
			}
			if password != passwordConfirmation {
				fmt.Println(au.BrightRed("Passwords do not match"))
				continue
			}
			hasValidPassword = true
		} else {
			return password, nil
		}
	}

	return password, nil
}

// ValidatePasswordInput validates a strong password input for new accounts,
// including a min length.
func ValidatePasswordInput(input string) error {
	if len(input) < minPasswordLength {
		return errPasswordWeak
	}

	return nil
}
