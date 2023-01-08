package workflow

import (
	"fmt"
	"os"
	"testing"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/exec"
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
		runtime, err := GetRuntime(exec.NewEnv(), runtimes, test.name)
		assert.Nil(t, err)
		assert.Equal(t, test.runtime, runtime.Config)
	}
}

func TestGetDefaultRuntime(t *testing.T) {
	runtime, err := GetRuntime(exec.NewEnv(), []config.Runtime{}, "")
	assert.Nil(t, err)
	assert.Equal(t, config.RUNTIME_NATIVE, runtime.Config.Type)
}

func TestGetCommand(t *testing.T) {
	wd, _ := os.Getwd()
	userHome, _ := os.UserHomeDir()

	buildContainer := config.Runtime{
		Name:  "build-container",
		Type:  "container",
		Image: "build-image:latest",
		Home:  "/home",
		Cache: []string{"/go"},
		Ports: []string{"8080:8080"},
	}
	portContainer := config.Runtime{
		Name:  "port-container",
		Type:  "container",
		Image: "port-image:latest",
		Home:  "/home",
		Cache: []string{"/test"},
		Ports: []string{"{{ .PORTS }}"},
	}
	globalCacheContainer := config.Runtime{
		Name:        "global-cache-container",
		Type:        "container",
		Image:       "gc:latest",
		Home:        "/home",
		GlobalCache: []string{"/gcache"},
	}

	env := exec.NewEnv()
	env["PORTS"] = "1234,80"
	env["HI"] = "hello"
	env["ISSUES"] = "false"

	tests := map[string]struct {
		runtime       *Runtime
		inputCommand  string
		interactive   bool
		outputCommand string
	}{
		"container": {
			runtime:       &Runtime{Env: exec.NewEnv(), Config: buildContainer},
			inputCommand:  "cmd",
			interactive:   true,
			outputCommand: fmt.Sprintf("docker run --rm  -v %s/.dredge/cache/go:/go -v %s:/home -p 8080:8080 -w /home -it build-image:latest cmd", wd, wd),
		},
		"env replace in container command": {
			runtime:       &Runtime{Env: env, Config: buildContainer},
			inputCommand:  "echo {{ .HI }}",
			interactive:   true,
			outputCommand: fmt.Sprintf("docker run --rm -e HI=hello -e ISSUES=false -e PORTS=1234,80 -v %s/.dredge/cache/go:/go -v %s:/home -p 8080:8080 -w /home -it build-image:latest echo hello", wd, wd),
		},
		"non-interactive container": {
			runtime:       &Runtime{Env: exec.NewEnv(), Config: buildContainer},
			inputCommand:  "cmd",
			interactive:   false,
			outputCommand: fmt.Sprintf("docker run --rm  -v %s/.dredge/cache/go:/go -v %s:/home -p 8080:8080 -w /home  build-image:latest cmd", wd, wd),
		},
		"container with ports": {
			runtime:       &Runtime{Env: env, Config: portContainer},
			inputCommand:  "cmd",
			interactive:   true,
			outputCommand: fmt.Sprintf("docker run --rm -e HI=hello -e ISSUES=false -e PORTS=1234,80 -v %s/.dredge/cache/test:/test -v %s:/home -p 1234:1234 -p 80:80 -w /home -it port-image:latest cmd", wd, wd),
		},
		"container without ports": {
			runtime:       &Runtime{Env: exec.NewEnv(), Config: portContainer},
			inputCommand:  "cmd",
			interactive:   true,
			outputCommand: fmt.Sprintf("docker run --rm  -v %s/.dredge/cache/test:/test -v %s:/home  -w /home -it port-image:latest cmd", wd, wd),
		},
		"container with global cache": {
			runtime:       &Runtime{Env: exec.NewEnv(), Config: globalCacheContainer},
			inputCommand:  "cmd",
			interactive:   true,
			outputCommand: fmt.Sprintf("docker run --rm  -v %s/.dredge/cache/global-cache-container/gcache:/gcache -v %s:/home  -w /home -it gc:latest cmd", userHome, wd),
		},
		"native": {
			runtime:       &Runtime{Env: exec.NewEnv(), Config: config.Runtime{Type: "native"}},
			inputCommand:  "cmd",
			interactive:   true,
			outputCommand: "cmd",
		},
		"env replace in command": {
			runtime:       &Runtime{Env: env, Config: config.Runtime{Type: "native"}},
			inputCommand:  "echo {{ .HI }}",
			interactive:   true,
			outputCommand: "echo hello",
		},
		"command with ||": {
			runtime:       &Runtime{Env: env, Config: config.Runtime{Type: "native"}},
			inputCommand:  "test || out",
			interactive:   true,
			outputCommand: "test || out",
		},
		"container with ||": {
			runtime:       &Runtime{Env: env, Config: buildContainer},
			inputCommand:  "test || out",
			interactive:   true,
			outputCommand: fmt.Sprintf("docker run --rm -e HI=hello -e ISSUES=false -e PORTS=1234,80 -v %s/.dredge/cache/go:/go -v %s:/home -p 8080:8080 -w /home -it build-image:latest test || out", wd, wd),
		},
		"command with if": {
			runtime:       &Runtime{Env: env, Config: config.Runtime{Type: "native"}},
			inputCommand:  "gh repo create {{if isFalse .ISSUES}}--disable-issues{{end}}",
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
