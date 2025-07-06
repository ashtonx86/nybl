package diabeticerrors

import "fmt"

const (
	ALREADY_EXISTS = iota 
)

// Error related to data and store.
type DataError struct {
	Code int
	Provider string 
	Err error 
}

func (de DataError) Error() string {
	return fmt.Sprintf("%s: %s", de.Provider, de.Err)
}

func (de DataError) Unwrap() error {
	return de.Err
}

var AlreadyExistsError = func(provider string, err error) error { 
	return DataError{
		Code: ALREADY_EXISTS,
		Provider: provider,
		Err: fmt.Errorf("data error :: Code %d, Provider: %s, Error: %w", ALREADY_EXISTS, provider, err),
	}
}
