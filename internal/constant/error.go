package constant

import "errors"

// Network errors.
var (
	InvalidJSON   = errors.New("invalid json")
	InvalidSchema = errors.New("invalid json schema")
)
