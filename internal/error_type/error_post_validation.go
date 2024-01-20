package errorType

import (
	"errors"
	"fmt"
)

type ErrorPostValidation struct {
	Err error
}

func NewErrorPostValidation(text string) *ErrorPostValidation {
	return &ErrorPostValidation{
		Err: errors.New(text),
	}
}

func (uve *ErrorPostValidation) Error() string {
	return fmt.Sprintf("%s", uve.Err)
}
