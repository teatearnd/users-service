package validations

import (
	"fmt"
	"net/mail"
	"strings"
	"unicode"
)

func ValidateEmail(email string, allowedEmails []string) error {
	e := strings.TrimSpace(strings.ToLower(email))
	if e == "" {
		return fmt.Errorf("email cant be blank")
	}

	parsed, err := mail.ParseAddress(e)
	if err != nil {
		return fmt.Errorf("invalid mail format")
	}
	addr := strings.ToLower(parsed.Address)

	for _, domain := range allowedEmails {
		domain = strings.TrimSpace(strings.ToLower(domain))
		if domain == "" {
			continue
		}
		if strings.HasSuffix(addr, "@"+domain) {
			return nil
		}
	}
	return fmt.Errorf("email domain is not allowed")
}

func ValidatePassword(pass string) error {
	if len(pass) < 8 || len(pass) > 72 {
		return fmt.Errorf("password must be 8-72 characters")
	}
	//	const allowedSpecial = "!@#$%^&*()-_=+[]{};:,.?/\\|~"
	for _, r := range pass {
		if r > unicode.MaxASCII {
			return fmt.Errorf("password must contain only latin/ascii characters")
		}
	}
	return nil
}
