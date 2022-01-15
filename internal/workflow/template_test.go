package workflow

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestExecuteTemplate(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("drg-%d", rand.Intn(100000)))
	defer os.Remove(tmpFile)

	step := config.TemplateStep{
		Input: "Hello {{ .test }}",
		Dest:  config.TemplateString(tmpFile),
	}
	workflow := config.Workflow{
		Name:        "workflow",
		Description: "My workflow",
		Inputs: map[string]string{
			"test": "world",
		},
		Steps: []config.Step{
			{
				Name:     nil,
				Template: &step,
			},
		},
	}
	dredgeFile := config.DredgeFile{
		Env:       config.Env{},
		Workflows: []config.Workflow{workflow},
	}

	env := NewEnv()
	env["test"] = "world"

	err := executeTemplate(&dredgeFile, workflow, &step, env)
	assert.Nil(t, err)

	content, err := ioutil.ReadFile(tmpFile)
	assert.Nil(t, err)
	assert.Equal(t, "Hello world", string(content))
}

func TestTemplate(t *testing.T) {
	tests := map[string]struct {
		input  config.TemplateString
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
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		output, err := Template(test.input, test.env)
		assert.Equal(t, test.output, output)
		assert.Equal(t, test.err, err)
	}
}
