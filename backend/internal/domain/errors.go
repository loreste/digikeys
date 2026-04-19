package domain

import "errors"

var (
	ErrNotFound        = errors.New("resource not found")
	ErrAlreadyExists   = errors.New("resource already exists")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrInvalidInput    = errors.New("invalid input")
	ErrBiometricFailed = errors.New("biometric operation failed")
	ErrCardExpired     = errors.New("card has expired")
	ErrInvalidMRZ      = errors.New("invalid MRZ data")
	ErrSyncFailed      = errors.New("data synchronization failed")
)
