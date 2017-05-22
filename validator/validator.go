package validator

import "regexp"

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

var passwordRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]{8,}$")

const (
	MinNum = 1
	MaxNum = 2147483647
)

func ValidateEmail(email string) bool {
	if emailRegexp.MatchString(email) {
		return true
	}
	return false
}

func ValidatePassword(password string) bool {
	// The rules are:
	// - Length >= 8
	// - Supported special characters: !#$%^&*_+=-}{|/.'`~
	if passwordRegexp.MatchString(password) {
		return true
	}
	return false
}

func ValidateNumber(num int) bool {
	if num >= MinNum && num <= MaxNum {
		return true
	}
	return false
}
