package validate

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseValidateOptions(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		options     validateOptions
		expectedErr error
	}{
		{
			name:        "required",
			input:       "required",
			options:     validateOptions{Required: true},
			expectedErr: nil,
		},
		{
			name:        "nothing",
			input:       "",
			options:     validateOptions{},
			expectedErr: nil,
		},
		{
			name:        "regex",
			input:       "regex=hello world",
			options:     validateOptions{Regex: "hello world"},
			expectedErr: nil,
		},
		{
			name:        "regex and required",
			input:       "regex=hello world,required",
			options:     validateOptions{Required: true, Regex: "hello world"},
			expectedErr: nil,
		},

		// TODO: this should throw an error
		{
			name:        "invalid",
			input:       "inv",
			options:     validateOptions{},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			options, err := parseValidateOptions(test.input)
			if err != test.expectedErr {
				t.Errorf("parseValidateOptions failed on test '%s'. Expected error '%v'. Got '%v'", test.name, test.expectedErr, err)
			}
			if !reflect.DeepEqual(options, test.options) {
				t.Errorf("parseValidateOptions failed on test '%s'. Expected options '%+v'. Got '%+v'", test.name, test.options, options)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	type myStruct struct {
		Col1 string `validate:"required"`
		Col2 string `validate:"regex=hello$"`
	}

	s := myStruct{
		Col1: "",
		Col2: "hello",
	}

	err := Validate(s)
	fmt.Println(err)
}
