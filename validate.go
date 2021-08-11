package validate

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const validateTagName = "validate"

type validateOptions struct {
	Required bool
	Regex    string
}

func Validate(obj interface{}) error {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		fieldValue := v.Field(i)

		if validateValue, ok := fieldType.Tag.Lookup(validateTagName); ok {
			options, err := parseValidateOptions(validateValue)
			if err != nil {
				return fmt.Errorf("Failed to parseValidateOptions:%w", err)
			}

			if err := validate(fieldValue, fieldType, options); err != nil {
				return err
			}
		}
	}

	return nil
}

func parseValidateOptions(tagVal string) (validateOptions, error) {
	var options validateOptions
	values := strings.Split(tagVal, ",")

	for _, v := range values {
		kv := strings.Split(v, "=")
		switch kv[0] {
		case "required":
			options.Required = true
		case "regex":
			if len(kv) != 2 {
				return options, fmt.Errorf("parseValidOptions: regex option did not have an expression")
			}
			options.Regex = kv[1]
		}
	}

	return options, nil

}

func validate(field reflect.Value, ftype reflect.StructField, options validateOptions) error {
	strVal := fmt.Sprintf("%v", field.Interface())

	if options.Required && strVal == "" {
		return fmt.Errorf("required field '%s' did not have a value", ftype.Name)
	}

	if options.Regex != "" {
		matched, err := regexp.MatchString(options.Regex, strVal)
		if err != nil {
			return err
		}
		if !matched {
			return fmt.Errorf("field '%s' failed to pass regex '%s': '%s'", ftype.Name, options.Regex, strVal)
		}
	}

	return nil
}
