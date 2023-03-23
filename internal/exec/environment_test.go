package exec

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestAddInput(t *testing.T) {
	tests := map[string]struct {
		input     config.Input
		reader    io.Reader
		osEnv     map[string]string
		outputEnv map[string]string
	}{
		"textInEnv": {
			input:     config.Input{Name: "title"},
			reader:    bytes.NewReader([]byte("")),
			osEnv:     map[string]string{"title": "the title"},
			outputEnv: map[string]string{"title": "the title"},
		},
		"textFromInput": {
			input:     config.Input{Name: "title"},
			reader:    bytes.NewReader([]byte("the title")),
			osEnv:     map[string]string{},
			outputEnv: map[string]string{"title": "the title"},
		},
		"selectInEnv": {
			input:     config.Input{Name: "city", Type: "select", Values: []string{"Brussels", "Barcelona"}},
			reader:    nil,
			osEnv:     map[string]string{"city": "Brussels"},
			outputEnv: map[string]string{"city": "Brussels"},
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		env := NewEnv()
		for k, v := range test.osEnv {
			os.Setenv(k, v)
		}
		err := env.AddInput(test.input, test.reader)
		assert.Nil(t, err)
		for k, v := range test.outputEnv {
			assert.Equal(t, v, env[k])
		}
		for k, _ := range test.osEnv {
			os.Unsetenv(k)
		}
	}
}

func TestAddVariables(t *testing.T) {
	env := NewEnv()
	env["first"] = "set"

	env.AddVariables(config.Variables{
		"first":  "overwritten",
		"second": "from vars",
	})

	assert.Equal(t, "set", env["first"])
	assert.Equal(t, "from vars", env["second"])
	assert.Equal(t, len(env), 2)
}

func TestTemplate(t *testing.T) {
	tests := map[string]struct {
		input  string
		env    Env
		output string
		err    error
	}{
		"no replaces": {
			input:  "test",
			env:    NewEnv(),
			output: "test",
			err:    nil,
		},
		"variable": {
			input:  "hello {{ .test }}",
			env:    Env{"test": "world"},
			output: "hello world",
			err:    nil,
		},
		"replace function": {
			input:  "{{replace .test \" \" \"-\" }}",
			env:    Env{"test": "hello world"},
			output: "hello-world",
			err:    nil,
		},
		"date function": {
			input:  "{{ date \"2006-01-02\" }}",
			env:    Env{},
			output: time.Now().Format("2006-01-02"),
			err:    nil,
		},
		"join with empty": {
			input:  "{{ join .arr .new \",\" }}",
			env:    Env{"new": "5000", "arr": ""},
			output: "5000",
			err:    nil,
		},
		"join empty": {
			input:  "{{ join .arr .new \",\" }}",
			env:    Env{"new": "", "arr": "5000"},
			output: "5000",
			err:    nil,
		},
		"join with one": {
			input:  "{{ join .arr .new \",\" }}",
			env:    Env{"new": "5000", "arr": "8000"},
			output: "8000,5000",
			err:    nil,
		},
		"join to list": {
			input:  "{{ join .arr .new \",\" }}",
			env:    Env{"new": "5000", "arr": "80,1234"},
			output: "80,1234,5000",
			err:    nil,
		},
		"isTrue true": {
			input:  "{{ if isTrue .val }}Hello{{ end }} {{ if isTrue .val2 }}world{{ end }}",
			env:    Env{"val": "true", "val2": "yes"},
			output: "Hello world",
			err:    nil,
		},
		"isTrue false": {
			input:  "{{ if isTrue .val }}Hello{{ end }} {{ if isTrue .val2 }}world{{ end }}",
			env:    Env{"val": "false", "val2": "no"},
			output: " ",
			err:    nil,
		},
		"isFalse true": {
			input:  "{{ if isFalse .val }}Hello{{ end }} {{ if isFalse .val2 }}world{{ end }}",
			env:    Env{"val": "false", "val2": "no"},
			output: "Hello world",
			err:    nil,
		},
		"isFalse false": {
			input:  "{{ if isFalse .val }}Hello{{ end }} {{ if isFalse .val2 }}world{{ end }}",
			env:    Env{"val": "true", "val2": "yes"},
			output: " ",
			err:    nil,
		},
		"trimSpace": {
			input:  "{{ trimSpace \" hello \" }}",
			env:    NewEnv(),
			output: "hello",
			err:    nil,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		e := &DredgeExec{Env: test.env}
		output, err := e.Template(test.input)
		assert.Equal(t, test.output, output)
		assert.Equal(t, test.err, err)
	}
}
