package storage

import "errors"

var ErrSessionRowMismatch = errors.New("Error with session")
var ErrSessionNotFound = errors.New("Session doesnt exists")
var ErrUsernameTaken = errors.New("Username already taken.")
var ErrUserNotFound = errors.New("Username not found")
var ErrUpdatingRequest = errors.New("Error Updating or User Already Rejected")
var ErrExpiredSession = errors.New("Session exists but is already expired")
