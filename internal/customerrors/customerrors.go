package customerrors

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrReadInConfig = fmt.Errorf("imopsible to read in config error")

	ErrInvalidInputBody = fmt.Errorf("invalid input body error")

	ErrUserNotFound = fmt.Errorf("user not found error")     // StatusBadRequest
	ErrUserIdNotInt = fmt.Errorf("user id is not int error") // StatusInternalServerError

	ErrWrongPassword = fmt.Errorf("wrong password error") // StatusUnauthorized

	ErrParseInventory = fmt.Errorf("imopsible to parse transaction error")    // StatusInternalServerError
	ErrParseTrx       = fmt.Errorf("imopsible to parse inventory item error") // StatusInternalServerError

	ErrGetBalance    = fmt.Errorf("get balance error")      // StatusInternalServerError
	ErrChangeBalance = fmt.Errorf("not enoghf coins error") // StatusBadRequest
	ErrUpdateBalance = fmt.Errorf("update balance error")   // StatusInternalServerError

	ErrSendCoinsToYousel = fmt.Errorf("it is impossible to send coins to yourself error") // StatusBadRequest

	ErrProductNotFound = fmt.Errorf("product not found error") // StatusBadRequest

	ErrTxStart = fmt.Errorf("start db transaction error") // StatusInternalServerError
	ErrTxStop  = fmt.Errorf("stop db transaction error")  // StatusInternalServerError

	ErrDatabase = fmt.Errorf("data base error") // StatusInternalServerError
)

func ClassifyError(err error) (int, string) {
	if errors.Is(err, ErrUserIdNotInt) ||
		errors.Is(err, ErrDatabase) ||
		errors.Is(err, ErrGetBalance) ||
		errors.Is(err, ErrUpdateBalance) ||

		errors.Is(err, ErrParseInventory) ||
		errors.Is(err, ErrParseTrx) ||

		errors.Is(err, ErrTxStart) ||
		errors.Is(err, ErrTxStop) {
		return http.StatusInternalServerError, err.Error()
	}

	if errors.Is(err, ErrUserNotFound) ||
		errors.Is(err, ErrChangeBalance) ||
		errors.Is(err, ErrSendCoinsToYousel) ||
		errors.Is(err, ErrProductNotFound) {
		return http.StatusBadRequest, err.Error()
	}

	if errors.Is(err, ErrWrongPassword) {
		return http.StatusUnauthorized, err.Error()
	}

	return http.StatusTeapot, err.Error()
}
