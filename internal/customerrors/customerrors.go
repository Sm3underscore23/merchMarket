package customerrors

import "fmt"

var (
	ErrUserNotFound = fmt.Errorf("user not found")
	ErrWrongPasswod = fmt.Errorf("wrong password")
)
