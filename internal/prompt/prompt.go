package prompt

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// Reader is an interface for reading input, used for testing.
type Reader interface {
	ReadPassword() (string, error)
	ReadLine() (string, error)
}

// TermReader reads from the actual terminal.
type TermReader struct{}

// ReadPassword reads a password from the terminal without echoing.
func (t *TermReader) ReadPassword() (string, error) {
	bytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}
	return string(bytes), nil
}

// ReadLine reads a line of input from the terminal.
func (t *TermReader) ReadLine() (string, error) {
	var line string
	_, err := fmt.Scanln(&line)
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}
	return strings.TrimSpace(line), nil
}

// Prompter handles user prompts for the CLI.
type Prompter struct {
	reader Reader
}

// New creates a new Prompter with the given reader.
func New(reader Reader) *Prompter {
	return &Prompter{reader: reader}
}

// NewDefault creates a Prompter that reads from the real terminal.
func NewDefault() *Prompter {
	return &Prompter{reader: &TermReader{}}
}

// AskPassword prompts the user for a password and confirms it.
func (p *Prompter) AskPassword(confirm bool) (string, error) {
	fmt.Print("Enter passphrase: ")
	password, err := p.reader.ReadPassword()
	fmt.Println()
	if err != nil {
		return "", err
	}
	if password == "" {
		return "", fmt.Errorf("passphrase cannot be empty")
	}
	if !confirm {
		return password, nil
	}
	fmt.Print("Confirm passphrase: ")
	confirmPassword, err := p.reader.ReadPassword()
	fmt.Println()
	if err != nil {
		return "", err
	}
	if password != confirmPassword {
		return "", fmt.Errorf("passphrases do not match")
	}
	return password, nil
}

// AskInput prompts the user for a generic string value.
func (p *Prompter) AskInput(label string) (string, error) {
	fmt.Printf("%s: ", label)
	value, err := p.reader.ReadLine()
	if err != nil {
		return "", err
	}
	return value, nil
}
