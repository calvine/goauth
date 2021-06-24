package errors

import "fmt"

type WrongTypeError struct {
	Expected, Actual string
}

func (e WrongTypeError) Error() string {
	return fmt.Sprintf("Wrong type found: Expected: %s - Actual: %s", e.Expected, e.Actual)
}
