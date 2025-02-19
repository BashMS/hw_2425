package hw09structvalidator

import "errors"

var (
	ErrExpectedStruct = errors.New("на вход ожидается структура")
	ErrValidValue     = errors.New("ошибка валидации")
)

var (
	strValidLenString = "%w. Длина строки должна быть ровно %v символа"
	strValidMinValue  = "%w. Число не может быть меньше %v"
	strValidMaxValue  = "%w. Число не может быть больше %v"
	strValidSetValue  = "%w. Значение должно входить в допустимое множество {%v}"
	strValidEmail     = "%w. Невалидный адрес почты %v"
)
