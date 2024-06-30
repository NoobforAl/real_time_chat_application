package appErrors

import "errors"

var (
	// TODO: setup more errors
	ErrNoDocuments      = errors.New("Not Found records")
	ErrDatabaseIndex    = errors.New("some filed is exist!")
	ErrNoFieldsToUpdate = errors.New("no filed change to update!")
	ErrNotValidId       = errors.New("Not valid Id")
)
