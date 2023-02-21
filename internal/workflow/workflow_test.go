package workflow

import (
	"bytes"
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/dredge-dev/dredge/internal/api"
	"github.com/dredge-dev/dredge/internal/config"
	"github.com/stretchr/testify/assert"
)

type CallbacksMock struct {
	MLog                        func(level api.LogLevel, msg string) error
	MRequestInput               func(inputRequests []api.InputRequest) (map[string]string, error)
	MOpenUrl                    func(url string) error
	MExecuteResourceCommand     func(resource string, command string) (*api.CommandOutput, error)
	MSetEnv                     func(name string, value interface{}) error
	MTemplate                   func(input string) (string, error)
	MAddVariablesToDredgefile   func(variable map[string]string) error
	MAddWorkflowToDredgefile    func(workflow config.Workflow) error
	MAddBucketToDredgefile      func(bucket config.Bucket) error
	MRelativePathFromDredgefile func(path string) (string, error)
	Env                         map[string]interface{}
}

func (c *CallbacksMock) Log(level api.LogLevel, msg string) error {
	if c.MLog != nil {
		return c.MLog(level, msg)
	}
	return fmt.Errorf("Log not mocked")
}
func (c *CallbacksMock) RequestInput(inputRequests []api.InputRequest) (map[string]string, error) {
	if c.MRequestInput != nil {
		return c.MRequestInput(inputRequests)
	}
	ret := make(map[string]string)
	for _, request := range inputRequests {
		ret[request.Name] = fmt.Sprintf("%s", c.Env[request.Name])
	}
	return ret, nil
}
func (c *CallbacksMock) OpenUrl(url string) error {
	if c.MOpenUrl != nil {
		return c.MOpenUrl(url)
	}
	return fmt.Errorf("OpenUrl not mocked")
}
func (c *CallbacksMock) ExecuteResourceCommand(resource string, command string) (*api.CommandOutput, error) {
	if c.MExecuteResourceCommand != nil {
		return c.MExecuteResourceCommand(resource, command)
	}
	return nil, fmt.Errorf("ExecuteResourceCommand not mocked")
}
func (c *CallbacksMock) SetEnv(name string, value interface{}) error {
	if c.MSetEnv != nil {
		return c.MSetEnv(name, value)
	}
	if c.Env == nil {
		c.Env = make(map[string]interface{})
	}
	c.Env[name] = value
	return nil
}
func (c *CallbacksMock) Template(input string) (string, error) {
	if c.MTemplate != nil {
		return c.MTemplate(input)
	}
	t, err := template.New("").Option("missingkey=zero").Parse(input)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %s", err)
	}
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c.Env); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
func (c *CallbacksMock) AddVariablesToDredgefile(variable map[string]string) error {
	if c.MAddVariablesToDredgefile != nil {
		return c.MAddVariablesToDredgefile(variable)
	}
	return fmt.Errorf("AddVariablesToDredgefile not mocked")
}
func (c *CallbacksMock) AddWorkflowToDredgefile(workflow config.Workflow) error {
	if c.MAddWorkflowToDredgefile != nil {
		return c.MAddWorkflowToDredgefile(workflow)
	}
	return fmt.Errorf("AddWorkflowToDredgefile not mocked")
}
func (c *CallbacksMock) AddBucketToDredgefile(bucket config.Bucket) error {
	if c.MAddBucketToDredgefile != nil {
		return c.MAddBucketToDredgefile(bucket)
	}
	return fmt.Errorf("AddBucketToDredgefile not mocked")
}
func (c *CallbacksMock) RelativePathFromDredgefile(path string) (string, error) {
	if c.MRelativePathFromDredgefile != nil {
		return c.MRelativePathFromDredgefile(path)
	}
	return path, nil
}

func TestExecuteShellStep(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("drg-%d", rand.Intn(100000)))
	defer os.Remove(tmpFile)

	c := &CallbacksMock{}
	workflow := &Workflow{
		Name:        "workflow",
		Description: "perform work",
		Steps: []config.Step{
			{
				Shell: &config.ShellStep{
					Cmd:    fmt.Sprintf("touch %s && echo hello && echo world >&2", tmpFile),
					StdOut: "OUTPUT",
					StdErr: "ERR",
				},
			},
		},
		Callbacks: c,
	}

	err := workflow.Execute()
	assert.Nil(t, err)

	_, err = os.Stat(tmpFile)
	assert.Nil(t, err)

	assert.Equal(t, "hello\n", c.Env["OUTPUT"])
	assert.Equal(t, "world\n", c.Env["ERR"])
}
