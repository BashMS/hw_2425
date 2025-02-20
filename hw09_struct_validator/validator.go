package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
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

func validateStrLen(rf reflect.Value, fieldName string, tag string) error {
	if rf.Kind() != reflect.String {
		return fmt.Errorf(strFieldTypeNotMatchValidationType, fieldName, "Len")
	}

	tagVal, err := strconv.Atoi(tag)
	if err != nil {
		return fmt.Errorf("strconv.Atoi: %w", err)
	}
	if len(rf.String()) != tagVal {
		return fmt.Errorf(strValidLenString, ErrValidValue, tagVal)
	}

	return nil
}

func validateMin(rf reflect.Value, fieldName string, tag string) error {
	if rf.Kind() != reflect.Int {
		return fmt.Errorf(strFieldTypeNotMatchValidationType, fieldName, "Min")
	}
	tagVal, err := strconv.ParseInt(tag, 10, 64)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt: %w", err)
	}
	if rf.Int() < tagVal {
		return fmt.Errorf(strValidMinValue, ErrValidValue, tagVal)
	}

	return nil
}

func validateMax(rf reflect.Value, fieldName string, tag string) error {
	if rf.Kind() != reflect.Int {
		return fmt.Errorf(strFieldTypeNotMatchValidationType, fieldName, "Max")
	}
	tagVal, err := strconv.ParseInt(tag, 10, 64)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt: %w", err)
	}
	if rf.Int() > tagVal {
		return fmt.Errorf(strValidMaxValue, ErrValidValue, tagVal)
	}

	return nil
}

func validateIn(val string, tag string) error {
	// Получим доступное множестов значений
	tagVals := strings.Split(tag, ",")
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

func validateRegexp(rf reflect.Value, fieldName string, tag string) error {
	if rf.Kind() != reflect.String {
		return fmt.Errorf(strFieldTypeNotMatchValidationType, fieldName, "Regexp")
	}
	// Получим доступное множестов значений
	re, err := regexp.Compile(tag)
	if err != nil {
		return fmt.Errorf("invalid regexp.Compile: %w", err)
	}
	if !re.MatchString(rf.String()) {
		return fmt.Errorf(strValidEmail, ErrValidValue, rf.String())
	}

	return nil
}

func validateItem(tag string, rf reflect.Value, fieldName string) error {
	// разберем теги валидации
	tags := strings.Split(tag, "|")
	for _, tgItem := range tags {
		switch {
		case strings.Contains(tgItem, "len:"):
			return validateStrLen(rf, fieldName, strings.TrimLeft(tgItem, "len:"))
		case strings.Contains(tgItem, "min:"):
			return validateMin(rf, fieldName, strings.TrimLeft(tgItem, "min:"))
		case strings.Contains(tgItem, "max:"):
			return validateMax(rf, fieldName, strings.TrimLeft(tgItem, "max:"))
		case strings.Contains(tgItem, "in:"):
			return validateIn(rf.String(), strings.TrimLeft(tgItem, "in:"))
		case strings.Contains(tgItem, "regexp:"):
			return validateRegexp(rf, fieldName, strings.TrimLeft(tgItem, "regxp:"))
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
			err := validateItem(tag, fieldValue, fieldType.Name)
			if err != nil && !errors.Is(err, ErrValidValue) {
				return err
			}
			if err != nil {
				valResult = append(valResult, ValidationError{
					Field: fieldType.Name,
					Err:   err,
				},
				)
			}
			continue
		}
		// слайс валидируем по значениям
		for i := 0; i < fieldValue.Len(); i++ {
			err := validateItem(tag, fieldValue.Index(i), fieldType.Name)
			if err != nil && !errors.Is(err, ErrValidValue) {
				return err
			}
			if err != nil {
				valResult = append(valResult, ValidationError{
					Field: fieldType.Name,
					Err:   err,
				},
				)
			}
		}
	}

	return valResult
}
