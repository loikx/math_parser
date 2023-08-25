package validation

import "errors"

var (
	ErrInvalidBracesSequence = errors.New("invalid braces sequence")
	ErrInvalidRegex          = errors.New("invalid regex")
	ErrInvalidOperators      = errors.New("invalid operators")
)
