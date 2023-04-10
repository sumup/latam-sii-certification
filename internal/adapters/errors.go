package adapters

import "errors"

var (
	errCantBuildRequest  = errors.New("err cant build request")
	errRequestCallFailed = errors.New("err request call failed")
	errCantReadBody      = errors.New("err cant read body")
)
