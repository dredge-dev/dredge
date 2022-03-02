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
		name        string
		description string
		input       io.Reader
		osEnv       map[string]string
		outputEnv   map[string]string
	}{
		"inEnv": {
			name:        "title",
			description: "Title of the thing",
			input:       bytes.NewReader([]byte("")),
			osEnv:       map[string]string{"title": "the title"},
			outputEnv:   map[string]string{"title": "the title"},
		},
		"fromInput": {
			name:        "title",
			description: "Title of the thing",
			input:       bytes.NewReader([]byte("the title")),
			osEnv:       map[string]string{},
			outputEnv:   map[string]string{"title": "the title"},
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		env := NewEnv()
		for k, v := range test.osEnv {
			os.Setenv(k, v)
		}
		err := env.AddInput(test.name, test.description, test.input)
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
