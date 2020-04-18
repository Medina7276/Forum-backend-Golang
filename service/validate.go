package service

import (
	"errors"
	"strings"
	"unicode"

	"git.01.alem.school/qjawko/forum/model"
)

func withinCharset(s, charset string) bool {
	for _, r := range s {
		if !strings.ContainsRune(charset, r) {
			return false
		}
	}

	return true
}

func validateUsername(s string) error {
	const legalcharset = ".0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	if len(s) < 3 {
		return errors.New("Username should be at least 3 characters long")
	}

	if !withinCharset(s, legalcharset) {
		return errors.New("Username shold only contain english letters, numbers, or dots ('.')")
	}

	return nil
}

func validatePassword(s string) error {
	if len(s) < 6 {
		return errors.New("Password should be at least 6 characters long")
	}

	legal := ""
	for i := '!'; i <= '~'; i++ {
		legal += string(i)
	}

	for _, r := range s {
		if !strings.ContainsRune(legal, r) {
			return errors.New("Password shold only contain english letters, numbers, or special characters ('.')")
		}
	}

	return nil
}

func validateName(s string) error {

	for _, r := range s {
		if r == ' ' {
			continue
		}

		if !unicode.IsLetter(r) {
			return errors.New("Name should contain only letters and spaces")
		}
	}

	return nil
}

func ValidateUser(u *model.User) error {

	if err := validateName(u.Name); err != nil {
		return err
	}
	if err := validateUsername(u.Username); err != nil {
		return err
	}

	return validatePassword(u.Password)
}

func validateLogin(u model.User) error {

	if err := validateUsername(u.Username); err != nil {
		return err
	}

	return validatePassword(u.Password)

}
