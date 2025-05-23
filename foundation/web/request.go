package web

import (
	"fmt"
	"io"
	"net/http"
)

// Param returns the parameters from the request path
func Param(r *http.Request, key string) string {
	return r.PathValue(key)
}

// Decoder represents data that can be decoded
type Decoder interface {
	Decode(data []byte) error
}

// Validator represents data that can be validated
type validator interface {
	Validate() error
}

// Decode reads from the HTTP request, decodes the data into a structure, and optionally validates the data
func Decode(r *http.Request, v Decoder) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("cannot read request payload: %w", err)
	}

	if err := v.Decode(data); err != nil {
		return fmt.Errorf("cannot decode request payload: %w", err)
	}

	if v, ok := v.(validator); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}
