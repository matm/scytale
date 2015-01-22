package scytale

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"code.google.com/p/go.crypto/ssh/terminal"
)

// ReadPassword reads a password of minLen lenght on the CLI
// and returns it. If twice is true, asks for confirmation.
// An error is returned if the password is empty or passwords don't
// match, of length of password is lesser than minLen.
func ReadPassword(minLen int, twice bool) (string, error) {
	fmt.Printf("Password: ")
	pwd, err := terminal.ReadPassword(int(os.Stdout.Fd()))
	if err != nil {
		return "", nil
	}
	password := strings.Trim(string(pwd), " ")
	fmt.Println()
	if len(password) < minLen {
		return "", errors.New("Password too short")
	}
	if !twice {
		return password, nil
	}
	fmt.Printf("Repeat: ")
	pwd2, err := terminal.ReadPassword(int(os.Stdout.Fd()))
	if err != nil {
		return "", err
	}
	fmt.Println()
	confirm := string(pwd2)
	if password != confirm {
		return "", errors.New("Passwords mismatch")
	}
	return password, nil
}
