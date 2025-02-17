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
	ErrValidStrValue  = "%w.строка должна входить в множество строк {%v}"
	strValidMinValue  = "%w.число не может быть меньше %v"
	strValidMaxValue  = "%w.число не может быть больше %v"
	ErrValidNumValue  = "%w.число должно входить в множество чисел {%v}"
)
