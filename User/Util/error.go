package Util

import "fmt"

var (
	ErrEmailExists    = fmt.Errorf("email already exists")
	ErrInternal       = fmt.Errorf("internal error")
	ErrUserEmailData  = fmt.Errorf("email and username cannot be empty")
	ErrHash           = fmt.Errorf("error generating hashed password")
	ErrBcrypt         = fmt.Errorf("password doesnot match")
	ErrCredentials    = fmt.Errorf("invalid credentials")
	ErrInvalidToken   = fmt.Errorf("invalid token")
	ErrExpiredToken   = fmt.Errorf("expired token")
	ErrInvalidRequest = fmt.Errorf("invalid request")
	ErrRoomExist      = fmt.Errorf("room not exist")
	ErrAlreadyJoined  = fmt.Errorf("user already in room")
	ErrRoomExists     = fmt.Errorf("room already exist")
	ErrRoomPublish    = fmt.Errorf("error publishing room")
)
