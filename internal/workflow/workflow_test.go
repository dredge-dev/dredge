package workflow

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestExecuteShellStep(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("drg-%d", rand.Intn(100000)))
	defer os.Remove(tmpFile)

	dredgeFile := &config.DredgeFile{
		Workflows: []config.Workflow{
			{
				Name:        "workflow",
				Description: "perform work",
				Steps: []config.Step{
					{
						Shell: &config.ShellStep{
							Cmd: fmt.Sprintf("touch %s", tmpFile),
						},
					},
				},
			},
		},
	}

	workflow, err := dredgeFile.GetWorkflow("", "workflow")
	assert.Nil(t, err)

	err = ExecuteWorkflow(dredgeFile, *workflow)
	assert.Nil(t, err)

	_, err = os.Stat(tmpFile)
	assert.Nil(t, err)
}

func TestImportWorkflowSameDedgeFile(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("drg-%d", rand.Intn(100000)))
	defer os.Remove(tmpFile)

	dredgeFile := &config.DredgeFile{
		Workflows: []config.Workflow{
			{
				Name: "workflow",
				Import: &config.ImportWorkflow{
					Bucket:   "b1",
					Workflow: "w1",
				},
			},
		},
		Buckets: []config.Bucket{
			{
				Name: "b1",
				Workflows: []config.Workflow{
					{
						Name:        "w1",
						Description: "perform work",
						Steps: []config.Step{
							{
								Shell: &config.ShellStep{
									Cmd: fmt.Sprintf("touch %s", tmpFile),
								},
							},
						},
					},
				},
			},
		},
	}

	workflow, err := dredgeFile.GetWorkflow("", "workflow")
	assert.Nil(t, err)

	err = ExecuteWorkflow(dredgeFile, *workflow)
	assert.Nil(t, err)

	_, err = os.Stat(tmpFile)
	assert.Nil(t, err)
}

func TestImportWorkflowDedgeFile(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("drg-%d", rand.Intn(100000)))
	defer os.Remove(tmpFile)

	remoteDredgeFile := config.DredgeFile{
		Buckets: []config.Bucket{
			{
				Name: "b1",
				Workflows: []config.Workflow{
					{
						Name:        "w1",
						Description: "perform work",
						Steps: []config.Step{
							{
								Shell: &config.ShellStep{
									Cmd: fmt.Sprintf("touch %s", tmpFile),
								},
							},
						},
					},
				},
			},
		},
	}

	remoteDredgeFileContent, err := yaml.Marshal(remoteDredgeFile)
	assert.Nil(t, err)

	remoteDredgeFilePath := "./test-import-workflow.DredgeFile"
	err = ioutil.WriteFile(remoteDredgeFilePath, remoteDredgeFileContent, 0644)
	assert.Nil(t, err)
	defer os.Remove(remoteDredgeFilePath)

	dredgeFile := &config.DredgeFile{
		Workflows: []config.Workflow{
			{
				Name: "workflow",
				Import: &config.ImportWorkflow{
					Source:   remoteDredgeFilePath,
					Bucket:   "b1",
					Workflow: "w1",
				},
			},
		},
	}

	workflow, err := dredgeFile.GetWorkflow("", "workflow")
	assert.Nil(t, err)

	err = ExecuteWorkflow(dredgeFile, *workflow)
	assert.Nil(t, err)

	_, err = os.Stat(tmpFile)
	assert.Nil(t, err)
}
