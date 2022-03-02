package workflow

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/exec"
	"github.com/stretchr/testify/assert"
)

func TestExecuteShellStep(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("drg-%d", rand.Intn(100000)))
	defer os.Remove(tmpFile)

	workflow := &exec.Workflow{
		Exec:        exec.EmptyExec(""),
		Name:        "workflow",
		Description: "perform work",
		Steps: []config.Step{
			{
				Shell: &config.ShellStep{
					Cmd: fmt.Sprintf("touch %s", tmpFile),
				},
			},
		},
	}

	err := ExecuteWorkflow(workflow)
	assert.Nil(t, err)

	_, err = os.Stat(tmpFile)
	assert.Nil(t, err)
}
