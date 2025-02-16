package hw09structvalidator

import "errors"

var (
	ErrExpectedStruct  = errors.New("на вход ожидается структура")
	ErrValidStrContain = errors.New("строка должна состоять из цифр")
	ErrValidLenString  = errors.New("неверная длина строки")
)

var (
	strValidLenString = "%w. Длина строки должна быть ровно %v символа"
	ErrValidStrValue  = "строка должна входить в множество строк {%v}"
	ErrValidMinValue  = "число не может быть меньше %v"
	ErrValidMaxValue  = "число не может быть больше %v"
	ErrValidNumValue  = "число должно входить в множество чисел {%v}"
)
