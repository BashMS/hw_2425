package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: App{Version: "lan"},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   fmt.Errorf(strValidLenString, ErrValidLenString, int(5)),
				},
			},
		},
		{
			in: Response{Code: 50, Body: "test"},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   fmt.Errorf(strValidSetValue, ErrValidValue, []string{"200", "404", "500"}),
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			if err := Validate(tt.in); err != nil {
				if err.Error() != tt.expectedErr.Error() {
					t.Errorf("Validate() error = %v, expectedErr %v", err, tt.expectedErr)
				}
			}
			_ = tt
		})
	}
}
