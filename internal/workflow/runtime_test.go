package workflow

import (
	"fmt"
	"os"
	"testing"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGetRuntime(t *testing.T) {
	r1 := config.Runtime{
		Name:  "build-container",
		Type:  "container",
		Image: "build-image:latest",
		Home:  "/home",
		Cache: []string{"/go"},
		Ports: []string{"8080:8080"},
	}
	r2 := config.Runtime{
		Name:  "port-container",
		Type:  "container",
		Image: "port-image:latest",
		Home:  "/home",
		Cache: []string{"/test"},
		Ports: []string{"{{ .PORTS }}"},
	}
	r3 := config.Runtime{
		Name: "native",
		Type: "native",
	}

	runtimes := []config.Runtime{r1, r2, r3}
	workflow := &Workflow{
		Runtimes:  runtimes,
		Callbacks: &CallbacksMock{},
	}

	tests := map[string]struct {
		name    string
		runtime config.Runtime
	}{
		"build-container": {
			name:    "build-container",
			runtime: r1,
		},
		"ports-container": {
			name:    "port-container",
			runtime: r2,
		},
		"native": {
			name:    "native",
			runtime: r3,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		runtime, err := workflow.GetRuntime(test.name)
		assert.Nil(t, err)
		assert.Equal(t, test.runtime, runtime.Config)
	}
}

func TestGetDefaultRuntime(t *testing.T) {
	workflow := &Workflow{
		Runtimes:  []config.Runtime{},
		Callbacks: &CallbacksMock{},
	}
	runtime, err := workflow.GetRuntime("")
	assert.Nil(t, err)
	assert.Equal(t, config.RUNTIME_NATIVE, runtime.Config.Type)
}

func TestGetCommand(t *testing.T) {
	wd, _ := os.Getwd()
	userHome, _ := os.UserHomeDir()

	buildContainer := config.Runtime{
		Name:    "build-container",
		Type:    "container",
		Image:   "build-image:latest",
		Home:    "/home",
		Cache:   []string{"/go"},
		Ports:   []string{"8080:8080"},
		EnvVars: map[string]string{"HI": "{{.HI}}", "ISSUES": "{{.ISSUES}}", "PORTS": "{{.PORTS}}"},
	}
	portContainer := config.Runtime{
		Name:    "port-container",
		Type:    "container",
		Image:   "port-image:latest",
		Home:    "/home",
		Cache:   []string{"/test"},
		Ports:   []string{"{{ .PORTS }}"},
		EnvVars: map[string]string{"HI": "{{.HI}}", "ISSUES": "{{.ISSUES}}", "PORTS": "{{.PORTS}}"},
	}
	globalCacheContainer := config.Runtime{
		Name:        "global-cache-container",
		Type:        "container",
		Image:       "gc:latest",
		Home:        "/home",
		GlobalCache: []string{"/gcache"},
		EnvVars:     map[string]string{"HI": "{{.HI}}", "ISSUES": "{{.ISSUES}}", "PORTS": "{{.PORTS}}"},
	}

	withEnv := &CallbacksMock{
		Env: map[string]interface{}{
			"PORTS":  "1234,80",
			"HI":     "hello",
			"ISSUES": "false",
		},
	}

	emptyEnv := &CallbacksMock{}

	tests := map[string]struct {
		runtime       *Runtime
		inputCommand  string
		interactive   bool
		outputCommand string
	}{
		"container": {
			runtime:       &Runtime{Config: buildContainer, Templater: emptyEnv.Template},
			inputCommand:  "cmd",
			interactive:   true,
			outputCommand: fmt.Sprintf("docker run --rm  -v %s/.dredge/cache/go:/go -v %s:/home -p 8080:8080 -w /home -it build-image:latest cmd", wd, wd),
		},
		"env replace in container command": {
			runtime:       &Runtime{Config: buildContainer, Templater: withEnv.Template},
			inputCommand:  "echo {{ .HI }}",
			interactive:   true,
			outputCommand: fmt.Sprintf("docker run --rm -e HI=hello -e ISSUES=false -e PORTS=1234,80 -v %s/.dredge/cache/go:/go -v %s:/home -p 8080:8080 -w /home -it build-image:latest echo hello", wd, wd),
		},
		"non-interactive container": {
			runtime:       &Runtime{Config: buildContainer, Templater: emptyEnv.Template},
			inputCommand:  "cmd",
			interactive:   false,
			outputCommand: fmt.Sprintf("docker run --rm  -v %s/.dredge/cache/go:/go -v %s:/home -p 8080:8080 -w /home  build-image:latest cmd", wd, wd),
		},
		"container with ports": {
			runtime:       &Runtime{Config: portContainer, Templater: withEnv.Template},
			inputCommand:  "cmd",
			interactive:   true,
			outputCommand: fmt.Sprintf("docker run --rm -e HI=hello -e ISSUES=false -e PORTS=1234,80 -v %s/.dredge/cache/test:/test -v %s:/home -p 1234:1234 -p 80:80 -w /home -it port-image:latest cmd", wd, wd),
		},
		"container without ports": {
			runtime:       &Runtime{Config: portContainer, Templater: emptyEnv.Template},
			inputCommand:  "cmd",
			interactive:   true,
			outputCommand: fmt.Sprintf("docker run --rm  -v %s/.dredge/cache/test:/test -v %s:/home  -w /home -it port-image:latest cmd", wd, wd),
		},
		"container with global cache": {
			runtime:       &Runtime{Config: globalCacheContainer, Templater: emptyEnv.Template},
			inputCommand:  "cmd",
			interactive:   true,
			outputCommand: fmt.Sprintf("docker run --rm  -v %s/.dredge/cache/global-cache-container/gcache:/gcache -v %s:/home  -w /home -it gc:latest cmd", userHome, wd),
		},
		"native": {
			runtime:       &Runtime{Config: config.Runtime{Type: "native"}, Templater: emptyEnv.Template},
			inputCommand:  "cmd",
			interactive:   true,
			outputCommand: "cmd",
		},
		"env replace in command": {
			runtime:       &Runtime{Config: config.Runtime{Type: "native"}, Templater: withEnv.Template},
			inputCommand:  "echo {{ .HI }}",
			interactive:   true,
			outputCommand: "echo hello",
		},
		"command with ||": {
			runtime:       &Runtime{Config: config.Runtime{Type: "native"}, Templater: withEnv.Template},
			inputCommand:  "test || out",
			interactive:   true,
			outputCommand: "test || out",
		},
		"container with ||": {
			runtime:       &Runtime{Config: buildContainer, Templater: withEnv.Template},
			inputCommand:  "test || out",
			interactive:   true,
			outputCommand: fmt.Sprintf("docker run --rm -e HI=hello -e ISSUES=false -e PORTS=1234,80 -v %s/.dredge/cache/go:/go -v %s:/home -p 8080:8080 -w /home -it build-image:latest test || out", wd, wd),
		},
		"command with if": {
			runtime:       &Runtime{Config: config.Runtime{Type: "native"}, Templater: withEnv.Template},
			inputCommand:  "gh repo create {{if .ISSUES}}--disable-issues{{end}}",
			interactive:   true,
			outputCommand: "gh repo create --disable-issues",
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		cmd, err := test.runtime.GetCommand(test.interactive, test.inputCommand)
		assert.Nil(t, err)
		assert.Equal(t, test.outputCommand, cmd)
	}
}
