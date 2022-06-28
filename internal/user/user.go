package user

import (
	"errors"
	"strings"
)

func IsValidUsername(username string) error {
	if username == "" {
		return errors.New("username must not be empty")
	}
	arr := strings.Split(username, " ")
	if len(arr) == 1 {
		return nil
	} else {
		return errors.New("username must not contain any space")
	}
}
