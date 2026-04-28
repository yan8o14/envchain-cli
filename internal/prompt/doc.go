// Package prompt provides utilities for interacting with the user via the
// terminal. It abstracts password input (with optional confirmation) and
// plain-text line input behind a Reader interface so that the behaviour can
// be replaced with a mock during testing.
//
// Typical usage:
//
//	p := prompt.NewDefault()
//	passphrase, err := p.AskPassword(true)  // ask + confirm
//	if err != nil {
//		log.Fatal(err)
//	}
package prompt
