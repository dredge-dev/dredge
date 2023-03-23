package cmd

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dredge-dev/dredge/internal/api"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func init() {
	table.DefaultHeaderFormatter = color.New(color.Bold).SprintfFunc()
}

func format(output *api.CommandOutput) (string, error) {
	if output.Type.Name == "object" {
		return formatPlain(output.Type.IsArray, output.Output)
	}

	out := new(strings.Builder)
	tbl := table.New(formatHeader(output.Type)...).WithWriter(out)

	if output.Type.IsArray {
		s := reflect.ValueOf(output.Output)
		if s.Kind() != reflect.Slice {
			return "", fmt.Errorf("expected array type but provider returned object")
		}
		for i := 0; i < s.Len(); i++ {
			row, err := formatObject(output.Type, s.Index(i).Interface())
			if err != nil {
				return "", err
			}
			tbl.AddRow(row...)
		}
	} else {
		row, err := formatObject(output.Type, output.Output)
		if err != nil {
			return "", err
		}
		tbl.AddRow(row...)
	}

	tbl.Print()
	return out.String(), nil
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

	out := new(strings.Builder)
	tbl := table.New("Field", "Value").WithWriter(out)
	for _, key := range s.MapKeys() {
		tbl.AddRow(key.String(), formatField(s.MapIndex(key)))
	}
	tbl.Print()
	return out.String(), nil
}

func formatHeader(t *api.Type) []interface{} {
	var output []interface{}
	for _, f := range t.Fields {
		output = append(output, f.Name)
	}
	return output
}

func formatObject(t *api.Type, o interface{}) ([]interface{}, error) {
	s := reflect.ValueOf(o)
	if s.Kind() != reflect.Map {
		return nil, fmt.Errorf("provider did not return a proper object")
	}
	var output []interface{}
	for _, f := range t.Fields {
		val := s.MapIndex(reflect.ValueOf(f.Name))
		output = append(output, formatField(val))
	}
	return output, nil
}

func formatField(o reflect.Value) string {
	if !o.IsValid() {
		return "<empty>"
	}
	return fmt.Sprintf("%s", o)
}
