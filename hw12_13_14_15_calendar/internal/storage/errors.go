package storage

import (
	"errors"
)

var (
	ErrDateBusy       = errors.New("данное время уже занято другим событием")
	ErrEventNotExists = errors.New("событие с указанным ID не существует")

	ErrUserExists    = errors.New("пользователь с таким адресом уже существует")
	ErrUserNotExists = errors.New("пользователь с указанным ID не существует")
)
