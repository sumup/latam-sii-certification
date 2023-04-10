package utils

import "errors"

var (
	ErrCountryNotConfigured = errors.New("country not configured")
	errCantCanonicalizeXML  = errors.New("cannot canonicalize xml")
)
