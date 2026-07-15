package storage

import "errors"

var ErrSessionRowMismatch = errors.New("Error with session")
var ErrSessionNotFound = errors.New("Session doesnt exists")
var ErrUsernameTaken = errors.New("Username already taken.")
var ErrUserNotFound = errors.New("Username not found")
var ErrUpdatingRequest = errors.New("Error Updating or User Already Rejected")
var ErrExpiredSession = errors.New("Session exists but is already expired")
var ErrFolderNotFound = errors.New("Folder not found")
var ErrUpdatingFolder = errors.New("Error Updating Folder")
var ErrFolderMismatch = errors.New("Error with Folder Query")
var ErrForbidden = errors.New("You do not own this")
