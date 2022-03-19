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
	"github.com/dredge-dev/dredge/internal/exec"
	"github.com/stretchr/testify/assert"
)

func TestExecuteTemplate(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("drg-%d", rand.Intn(100000)))
	defer os.Remove(tmpFile)

	step := config.TemplateStep{
		Input: "Hello {{ .test }}",
		Dest:  tmpFile,
	}
	workflow := &exec.Workflow{
		Exec:        exec.EmptyExec(""),
		Name:        "workflow",
		Description: "My workflow",
		Inputs: []config.Input{
			{
				Name: "test",
			},
		},
		Steps: []config.Step{
			{
				Name:     "",
				Template: &step,
			},
		},
	}

	os.Setenv("test", "value")
	err := ExecuteWorkflow(workflow)
	assert.Nil(t, err)

	content, err := ioutil.ReadFile(tmpFile)
	assert.Nil(t, err)
	assert.Equal(t, "Hello value", string(content))
}

func TestExecuteTemplateFromSource(t *testing.T) {
	dstFile := filepath.Join(os.TempDir(), fmt.Sprintf("drg-%d", rand.Intn(100000)))
	defer os.Remove(dstFile)

	templateFile := "./test-execute-template"
	err := ioutil.WriteFile(templateFile, []byte("Hello {{ .test }}"), 0644)
	defer os.Remove(templateFile)
	assert.Nil(t, err)

	workflow := &exec.Workflow{
		Exec:        exec.EmptyExec("./Dredgefile"),
		Name:        "workflow",
		Description: "My workflow",
		Inputs: []config.Input{
			{
				Name: "test",
			},
		},
		Steps: []config.Step{
			{
				Name: "",
				Template: &config.TemplateStep{
					Source: config.SourcePath(templateFile),
					Dest:   dstFile,
				},
			},
		},
	}

	os.Setenv("test", "value")
	err = ExecuteWorkflow(workflow)
	assert.Nil(t, err)

	content, err := ioutil.ReadFile(dstFile)
	assert.Nil(t, err)
	assert.Equal(t, "Hello value", string(content))
}

func TestTemplate(t *testing.T) {
	tests := map[string]struct {
		input  string
		env    exec.Env
		output string
		err    error
	}{
		"no replaces": {
			input:  "test",
			env:    exec.NewEnv(),
			output: "test",
			err:    nil,
		},
		"variable": {
			input:  "hello {{ .test }}",
			env:    exec.Env{"test": "world"},
			output: "hello world",
			err:    nil,
		},
		"replace function": {
			input:  "{{replace .test \" \" \"-\" }}",
			env:    exec.Env{"test": "hello world"},
			output: "hello-world",
			err:    nil,
		},
		"date function": {
			input:  "{{ date \"2006-01-02\" }}",
			env:    exec.Env{},
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
