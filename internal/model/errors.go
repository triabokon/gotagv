package model

import "fmt"

var (
	ErrInvalidArgument = fmt.Errorf("invalid argument")

	ErrNotFound      = fmt.Errorf("entity not found")
	ErrAlreadyExists = fmt.Errorf("entity already exists")
)
