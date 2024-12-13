package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

// RtTrim обрезает последний символ строки.
func RtTrim(sb *strings.Builder) error {
	trimStr := sb.String()
	sb.Reset()
	rnStr := []rune(trimStr)
	_, err := sb.WriteString(string(rnStr[0 : len(rnStr)-1]))
	if err != nil {
		return fmt.Errorf("%w WriteString", err)
	}
	return nil
}

// Unpack распаковывает строку.
func Unpack(str string) (string, error) {
	// Если пустая, сразу выходим
	if len(str) == 0 {
		return "", nil
	}

	// Пройдем по строке
	var (
		res    strings.Builder
		prevRn rune
	)

	// Если начали с цифры тогда сразу ошибка
	if unicode.IsDigit([]rune(str)[0]) {
		return "", ErrInvalidString
	}
	for _, item := range str {
		switch {
		case unicode.IsDigit(item) && unicode.IsDigit(prevRn):
			return "", ErrInvalidString
		case unicode.IsDigit(item):
			cnt, err := strconv.Atoi(string(item))
			if err != nil {
				return "", fmt.Errorf("%w strconv.Atoi", err)
			}
			switch cnt {
			case 1:
				// Ничего не делаем, символ добавлен в предыдущей итерации
				prevRn = item
				continue
			case 0:
				// Уберем последний добавленный символ
				err = RtTrim(&res)
				if err != nil {
					return "", fmt.Errorf("%w rtTrim", err)
				}
				prevRn = item
				continue
			default:
				// Добавим символ n-1 раз
				_, err = res.WriteString(strings.Repeat(string(prevRn), cnt-1))
				if err != nil {
					return "", fmt.Errorf("%w WriteString", err)
				}
				prevRn = item
			}
		default:
			prevRn = item
			res.WriteRune(item)
		}
	}

	return res.String(), nil
}
