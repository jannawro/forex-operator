package forex

import (
	"errors"
	"fmt"
	"strings"
)

// TODO: Error messages
var (
	ErrNoAuthToken          = errors.New("")
	ErrOpenExchangeAPIError = errors.New("")
)

type UnsupportedCurrenciesError struct {
	Currencies []string
}

func (e *UnsupportedCurrenciesError) Error() string {
	return fmt.Sprintf("unsupported currencies: %s", strings.Join(e.Currencies, ", "))
}
