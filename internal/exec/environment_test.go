package exec

import (
	"bytes"
	"io"
	"os"
	"testing"

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
