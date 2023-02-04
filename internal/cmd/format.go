package cmd

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dredge-dev/dredge/internal/resource"
)

func format(output *resource.CommandOutput) (string, error) {
	if output.Type.Name == "object" {
		return formatPlain(output.Type.IsArray, output.Output)
	}

	formatted, err := formatHeader(output.Type)
	if err != nil {
		return "", nil
	}

	if output.Type.IsArray {
		s := reflect.ValueOf(output.Output)
		if s.Kind() != reflect.Slice {
			return "", fmt.Errorf("expected array type but provider returned object")
		}
		for i := 0; i < s.Len(); i++ {
			line, err := formatObject(output.Type, s.Index(i).Interface())
			if err != nil {
				return "", err
			}
			formatted = formatted + "\n" + line
		}
	} else {
		line, err := formatObject(output.Type, output.Output)
		if err != nil {
			return "", err
		}
		formatted = formatted + "\n" + line
	}

	return formatted + "\n", nil
}

func formatPlain(isArray bool, o interface{}) (string, error) {
	if isArray {
		var formatted string

		s := reflect.ValueOf(o)
		if s.Kind() != reflect.Slice {
			return "", fmt.Errorf("expected array type but provider returned object")
		}
		for i := 0; i < s.Len(); i++ {
			line, err := formatPlainObject(s.Index(i).Interface())
			if err != nil {
				return "", err
			}
			formatted = formatted + "\n" + line
		}
		return formatted + "\n", nil
	}
	return formatPlainObject(o)
}

func formatPlainObject(o interface{}) (string, error) {
	s := reflect.ValueOf(o)
	if s.Kind() != reflect.Map {
		return "", fmt.Errorf("provider did not return a proper object")
	}

	var output string
	for _, key := range s.MapKeys() {
		if len(output) > 0 {
			output = output + "\n"
		}
		val := s.MapIndex(key)
		fieldValue, err := formatField(val)
		if err != nil {
			return "", err
		}
		output = output + key.String() + ":" + "\t" + fieldValue
	}
	return output + "\n", nil
}

func formatHeader(t *resource.Type) (string, error) {
	var output string
	for _, f := range t.Fields {
		if len(output) > 0 {
			output = output + "\t"
		}
		output = output + strings.ToUpper(f.Name)
	}
	return output, nil
}

func formatObject(t *resource.Type, o interface{}) (string, error) {
	s := reflect.ValueOf(o)
	if s.Kind() != reflect.Map {
		return "", fmt.Errorf("provider did not return a proper object")
	}

	var output string
	for _, f := range t.Fields {
		if len(output) > 0 {
			output = output + "\t"
		}
		val := s.MapIndex(reflect.ValueOf(f.Name))
		fieldValue, err := formatField(val)
		if err != nil {
			return "", err
		}
		output = output + fieldValue
	}
	return output, nil
}

func formatField(o reflect.Value) (string, error) {
	if !o.IsValid() {
		return "<empty>", nil
	}
	return fmt.Sprintf("%s", o), nil
}
