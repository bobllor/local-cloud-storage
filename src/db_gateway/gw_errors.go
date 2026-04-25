package dbgateway

import (
	"errors"
	"fmt"
)

var (
	UsernameEmptyErr            = errors.New("username cannot be empty")
	UsernameLenOutOfRangeErr    = fmt.Errorf("username must be between %d to %d characters long", usernameMinLength, usernameMaxLength)
	UsernameInvalidFirstCharErr = errors.New("username must start with a letter")
	UsernameInvalidEndCharErr   = errors.New("username must end with an alphanumeric character")
	UsernameIsInvalidErr        = errors.New("username may only consist of alphanumeric characters and single periods")
)

// IsUsernameError compares if the given error is a username validation error.
func IsUsernameError(err error) bool {
	errs := []error{
		UsernameEmptyErr,
		UsernameInvalidEndCharErr,
		UsernameInvalidFirstCharErr,
		UsernameIsInvalidErr,
		UsernameLenOutOfRangeErr,
	}

	for _, e := range errs {
		if errors.Is(err, e) {
			return true
		}
	}

	return false
}

// IsPasswordError compares if the given error is a password validation error.
func IsPasswordError(err error) bool {
	return false
}
