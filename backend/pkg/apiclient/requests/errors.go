package requests

import (
	"errors"
)

var (
	// ErrNoURLProvided is a handy error to return when we expect a *url.URL and don't receive one.
	ErrNoURLProvided = errors.New("no URL provided")

	// ErrNilEncoderProvided indicates a nil encoder was provided to the constructor.
	ErrNilEncoderProvided = errors.New("nil encoder provided")

	// ErrEmptyInputProvided indicates empty input was provided in an unacceptable context.
	ErrEmptyInputProvided = errors.New("empty input provided")

	// ErrNilInputProvided indicates nil input was provided in an unacceptable context.
	ErrNilInputProvided = errors.New("nil input provided")

	// ErrInvalidIDProvided indicates nil input was provided in an unacceptable context.
	ErrInvalidIDProvided = errors.New("required ID provided is zero")

	// ErrEmptyEmailAddressProvided indicates the user provided an empty username for search.
	ErrEmptyEmailAddressProvided = errors.New("empty email address provided")

	// ErrCookieRequired indicates a cookie is required.
	ErrCookieRequired = errors.New("cookie required for request")

	// ErrInvalidPhotoEncodingForUpload indicates the provided photo upload is of the wrong encoding.
	ErrInvalidPhotoEncodingForUpload = errors.New("invalid photo encoding")
)
