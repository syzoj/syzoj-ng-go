package server

import (
    "errors"
)

var ErrBusy = errors.New("Server busy")
var ErrNotImplemented = errors.New("Not implemented")
var ErrNotFound = errors.New("Not found")
var ErrCSRF = errors.New("CSRF token doesn't match")
