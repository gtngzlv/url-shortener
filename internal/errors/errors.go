package errors

import "errors"

var (
	ErrAlreadyExist    = errors.New("provided URL already shorted and stored in DB")
	ErrNoBatchByUserID = errors.New("no batches by provided userID")
)
