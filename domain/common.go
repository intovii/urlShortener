package domain

import "errors"

var ErrNotFound = errors.New("row does not exist")
var ErrInvalidArgument = errors.New("invalid param")