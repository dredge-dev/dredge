package workflow

import (
	"fmt"
	"os"
	"testing"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGetRuntime(t *testing.T) {
	home := "/home"
	wd, _ := os.Getwd()

	env := config.Env{
		Runtimes: []config.Runtime{
			{
				Name:  "build-container",
				Type:  "container",
				Image: "build-image:latest",
				Home:  &home,
				Cache: []string{"/go"},
				Ports: []string{"8080:8080"},
			},
			{
				Name: "native",
				Type: "native",
			},
		},
	}

	tests := map[string]struct {
		name    string
		command string
	}{
		"container": {
			name:    "build-container",
			command: fmt.Sprintf("docker run --rm  -v %s/.dredge/cache/go:/go -v %s:/home -p 8080:8080 -w /home -it build-image:latest ${cmd}", wd, wd),
		},
		"native": {
			name:    "native",
			command: "${cmd}",
		},
		"default": {
			name:    "",
			command: "${cmd}",
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		runtime, err := GetRuntime(env, test.name)
		assert.Nil(t, err)
		assert.Equal(t, test.command, runtime.Template)
	}
}

func TestGetCommand(t *testing.T) {
	env := config.Env{
		Variables: map[string]string{
			"var1": "test",
			"var2": "hello",
		},
	}

	tests := map[string]struct {
		runtime *Runtime
		command string
		output  string
	}{
		"empty env": {
			runtime: &Runtime{Env: config.Env{}, Template: "prefix ${cmd}"},
			command: "cmd",
			output:  "prefix cmd",
		},
		"env replace": {
			runtime: &Runtime{Env: env, Template: "prefix ${var1} ${cmd}"},
			command: "cmd ${var2}",
			output:  "prefix test cmd hello",
		},
		"env replace in runtime": {
			runtime: &Runtime{Env: env, Template: "prefix ${var1} ${cmd}"},
			command: "cmd",
			output:  "prefix test cmd",
		},
		"env replace in command": {
			runtime: &Runtime{Env: env, Template: "prefix ${cmd}"},
			command: "cmd ${var2}",
			output:  "prefix cmd hello",
		},
		"command with ||": {
			runtime: &Runtime{Env: env, Template: "${cmd}"},
			command: "test || out",
			output:  "test || out",
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		assert.Equal(t, test.output, test.runtime.getCommand(test.command))
	}
}
