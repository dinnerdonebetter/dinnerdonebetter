package database

import (
	"errors"
)

var (
	// ErrDatabaseNotReady indicates the given database is not ready.
	ErrDatabaseNotReady = errors.New("database is not ready")

	// ErrDatabaseCircuitBreakerTripped indicates the given database is not ready.
	ErrDatabaseCircuitBreakerTripped = errors.New("database circuit breaker is tripped")

	// ErrUserAlreadyExists indicates that a user with that username has already been created.
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrConfigRequired indicates that a null config was provided.
	ErrConfigRequired = errors.New("configuration required")

	// ErrAlreadyFinalized indicates a meal plan has already been finalized.
	ErrAlreadyFinalized = errors.New("meal plan already finalized")
)
