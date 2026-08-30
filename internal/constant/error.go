package constant

import "errors"

// Network errors.
var (
	InvalidJSON        = errors.New("invalid json")
	InvalidSchema      = errors.New("invalid json schema")
	InvalidParsedToken = errors.New("invalid parsed token")
)

// Database errors.
var (
	NoAffectedRows = errors.New("no affected rows")
)

// Bussiness validation errors.
var (
	InvalidURIPrefix = errors.New("invalid uri prefix")
	InvalidFileExt   = errors.New("invalid file extension")
	ErrOutOfRange    = errors.New("the value is out of range")
)
