package hw02unpackstring

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "a0b0c0", expected: ""},
		{input: "буквв0а1", expected: "буква"},
		{input: "буквыA1B2C3D4", expected: "буквыABBCCCDDDD"},

		// uncomment if task with asterisk completed
		// {input: `qwe\4\5`, expected: `qwe45`},
		// {input: `qwe\45`, expected: `qwe44444`},
		// {input: `qwe\\5`, expected: `qwe\\\\\`},
		// {input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func TestRtTrim(t *testing.T) {
	var (
		in1  strings.Builder
		in2  strings.Builder
		out1 strings.Builder
		out2 strings.Builder
	)
	in1.WriteString("aaaabccddddde")
	out1.WriteString("aaaabccddddd")
	in2.WriteString("a")
	out2.WriteString("")

	tests := []struct {
		input    *strings.Builder
		expected *strings.Builder
	}{
		{input: &in1, expected: &out1},
		{input: &in2, expected: &out2},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input.String(), func(t *testing.T) {
			err := RtTrim(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, tc.input)
		})
	}
}
