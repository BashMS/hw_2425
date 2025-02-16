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
		if strings.Contains(tag, "len:") && fieldValue.Kind() == reflect.String {
			tagVal, err := strconv.Atoi(strings.TrimLeft(tag, "len:"))
			if err != nil {
				return fmt.Errorf("strconv.Atoi: %w", err)
			}
			if len(fieldValue.String()) != tagVal {
				valResult = append(valResult, ValidationError{Field: fieldType.Name,
					Err: fmt.Errorf(strValidLenString, ErrValidLenString, tagVal)})
			}
		}
	}

	return valResult
}
