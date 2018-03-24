package discogs

import (
	"fmt"
	"strings"
)

// Error represents a Discogs API error
type Error struct {
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s", strings.ToLower(e.Message))
}

// APIErrors
var (
	ErrCurrencyNotSupported = &Error{"currency does not supported"}
	ErrUserAgentInvalid     = &Error{"invalid user-agent"}
)
