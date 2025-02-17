package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var result string
	for _, ve := range v {
		result = fmt.Sprintf("%s; %s=>'%s'", result, ve.Field, ve.Err.Error())
	}
	return strings.TrimLeft(result, "; ")
}

func validateStrLen(str string, tag string) error {
	tagVal, err := strconv.Atoi(strings.TrimLeft(tag, "len:"))
	if err != nil {
		return fmt.Errorf("strconv.Atoi: %w", err)
	}
	if len(str) != tagVal {
		return fmt.Errorf(strValidLenString, ErrValidLenString, tagVal)
	}

	return nil
}

func validateMin(val int64, tag string) error {
	tagVal, err := strconv.ParseInt(strings.TrimLeft(tag, "min:"), 10, 64)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt: %w", err)
	}
	if val < tagVal {
		return fmt.Errorf(strValidMinValue, ErrValidValue, tagVal)
	}

	return nil
}

func validateMax(val int64, tag string) error {
	tagVal, err := strconv.ParseInt(strings.TrimLeft(tag, "max:"), 10, 64)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt: %w", err)
	}
	if val > tagVal {
		return fmt.Errorf(strValidMaxValue, ErrValidValue, tagVal)
	}

	return nil
}

func validateIn(val interface{}, tag string) error {
	// Получим доступное множестов значений
	tagVals := strings.Split(strings.TrimLeft(tag, "in:"), ",")
	exists := false
	for _, tagVal := range tagVals {
		if val == tagVal {
			exists = true
			break
		}
	}
	if !exists {
		return fmt.Errorf(strValidSetValue, ErrValidValue, tagVals)
	}

	return nil
}

func validateItem(tag string, rf reflect.Value) error {
	// разберем теги валидации
	tags := strings.Split(tag, "|")
	for _, tgItem := range tags {
		switch {
		case strings.Contains(tgItem, "len:"):
			return validateStrLen(rf.String(), tgItem)
		case strings.Contains(tgItem, "min:"):
			return validateMin(rf.Int(), tgItem)
		case strings.Contains(tgItem, "min:"):
			return validateMax(rf.Int(), tgItem)
		case strings.Contains(tgItem, "in:"):
			return validateIn(rf, tgItem)
		}
	}
	return nil
}

func Validate(v interface{}) error {
	sv := reflect.ValueOf(v)
	t := reflect.TypeOf(v)

	if sv.Kind() != reflect.Struct {
		return ErrExpectedStruct
	}

	var valResult ValidationErrors

	for i := 0; i < sv.NumField(); i++ {
		fieldValue := sv.Field(i)
		fieldType := t.Field(i)
		tag := fieldType.Tag.Get("validate")
		if len(tag) == 0 {
			continue
		}
		// если значение не слайс тогда сразу валидируем
		if fieldValue.Kind() != reflect.Slice {
			err := validateItem(tag, fieldValue)
			if err != nil {
				valResult = append(valResult, ValidationError{Field: fieldType.Name,
					Err: err})
			}
			continue
		}
		// слайс валидируем по значениям
		for i := 0; i < fieldValue.Len(); i++ {
			err := validateItem(tag, fieldValue.Index(i))
			if err != nil {
				valResult = append(valResult, ValidationError{Field: fieldType.Name,
					Err: err})
			}
		}
	}

	return valResult
}
