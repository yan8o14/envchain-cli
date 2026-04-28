package prompt_test

import (
	"fmt"
	"testing"

	"github.com/user/envchain-cli/internal/prompt"
)

// mockReader is a test double for prompt.Reader.
type mockReader struct {
	passwords []string
	lines     []string
	pwIndex   int
	lineIndex int
	pwErr     error
	lineErr   error
}

func (m *mockReader) ReadPassword() (string, error) {
	if m.pwErr != nil {
		return "", m.pwErr
	}
	if m.pwIndex >= len(m.passwords) {
		return "", fmt.Errorf("no more passwords")
	}
	pw := m.passwords[m.pwIndex]
	m.pwIndex++
	return pw, nil
}

func (m *mockReader) ReadLine() (string, error) {
	if m.lineErr != nil {
		return "", m.lineErr
	}
	if m.lineIndex >= len(m.lines) {
		return "", fmt.Errorf("no more lines")
	}
	line := m.lines[m.lineIndex]
	m.lineIndex++
	return line, nil
}

func TestAskPasswordNoConfirm(t *testing.T) {
	reader := &mockReader{passwords: []string{"secret123"}}
	p := prompt.New(reader)
	pw, err := p.AskPassword(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pw != "secret123" {
		t.Errorf("expected 'secret123', got '%s'", pw)
	}
}

func TestAskPasswordWithConfirmMatch(t *testing.T) {
	reader := &mockReader{passwords: []string{"secret123", "secret123"}}
	p := prompt.New(reader)
	pw, err := p.AskPassword(true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pw != "secret123" {
		t.Errorf("expected 'secret123', got '%s'", pw)
	}
}

func TestAskPasswordWithConfirmMismatch(t *testing.T) {
	reader := &mockReader{passwords: []string{"secret123", "different"}}
	p := prompt.New(reader)
	_, err := p.AskPassword(true)
	if err == nil {
		t.Fatal("expected error for mismatched passphrases, got nil")
	}
}

func TestAskPasswordEmpty(t *testing.T) {
	reader := &mockReader{passwords: []string{""}}
	p := prompt.New(reader)
	_, err := p.AskPassword(false)
	if err == nil {
		t.Fatal("expected error for empty passphrase, got nil")
	}
}

func TestAskInput(t *testing.T) {
	reader := &mockReader{lines: []string{"my-project"}}
	p := prompt.New(reader)
	val, err := p.AskInput("Project name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "my-project" {
		t.Errorf("expected 'my-project', got '%s'", val)
	}
}
