package hw09structvalidator

import "errors"

var (
	ErrExpectedStruct  = errors.New("на вход ожидается структура")
	ErrValidStrContain = errors.New("строка должна состоять из цифр")
	ErrValidLenString  = errors.New("неверная длина строки")
	ErrValidValue      = errors.New("неверное значение")
)

var (
	strValidLenString = "%w. Длина строки должна быть ровно %v символа"
	strValidMinValue  = "%w. Число не может быть меньше %v"
	strValidMaxValue  = "%w. Число не может быть больше %v"
	strValidSetValue  = "%w. Значение должно входить в допустимое множество {%v}"
)
