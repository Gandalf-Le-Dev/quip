package domain

import "errors"

var (
    ErrNotFound      = errors.New("not found")
    ErrExpired       = errors.New("content has expired")
    ErrLimitExceeded = errors.New("download/view limit exceeded")
    ErrInvalidInput  = errors.New("invalid input")
)