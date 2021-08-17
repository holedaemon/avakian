package httpx

import "errors"

var ErrStatusCode = errors.New("httpx: unexpected status code returned")
