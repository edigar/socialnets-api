package errorType

import (
	"errors"
	"fmt"
)

type ErrorUserValidation struct {
	Err error
}

func NewErrorUserValidation(text string) *ErrorUserValidation {
	return &ErrorUserValidation{
		Err: errors.New(text),
	}
}

func (uve *ErrorUserValidation) Error() string {
	return fmt.Sprintf("%s", uve.Err)
}
