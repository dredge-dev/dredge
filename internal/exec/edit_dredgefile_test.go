package exec

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestExecuteEditDredgeFile(t *testing.T) {
	dredgeFile := fmt.Sprintf("./tmp-test-%d", rand.Intn(100000))
	importFile := fmt.Sprintf("%s-import", dredgeFile)
	defer os.Remove(dredgeFile)
	defer os.Remove(importFile)

	d, err := yaml.Marshal(config.DredgeFile{
		Workflows: []config.Workflow{
			{
				Name: "w1",
				Steps: []config.Step{
					{
						Browser: &config.BrowserStep{
							Url: "https://www.dredge.dev",
						},
					},
				},
			},
		},
		Buckets: []config.Bucket{
			{
				Name:      "b1",
				Workflows: []config.Workflow{},
			},
		},
	})
	assert.Nil(t, err)

	tests := map[string]struct {
		df      config.DredgeFile
		content config.DredgeFile
		errMsg  string
	}{
		"add workflow": {
			df: config.DredgeFile{
				Workflows: []config.Workflow{
					{
						Name: "add-workflow",
						Steps: []config.Step{
							{
								EditDredgeFile: &config.EditDredgeFileStep{
									AddWorkflows: []config.Workflow{
										{
											Name: "w2",
											Import: &config.ImportWorkflow{
												Workflow: "add-workflow",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			content: config.DredgeFile{
				Workflows: []config.Workflow{
					{
						Name: "w1",
						Steps: []config.Step{
							{
								Browser: &config.BrowserStep{
									Url: "https://www.dredge.dev",
								},
							},
						},
					},
					{
						Name: "w2",
						Import: &config.ImportWorkflow{
							Source:   config.SourcePath(importFile),
							Workflow: "add-workflow",
						},
					},
				},
				Buckets: []config.Bucket{
					{
						Name:      "b1",
						Workflows: []config.Workflow{},
					},
				},
			},
		},
		"add bucket": {
			df: config.DredgeFile{
				Workflows: []config.Workflow{
					{
						Name: "add-workflow",
						Steps: []config.Step{
							{
								EditDredgeFile: &config.EditDredgeFileStep{
									AddBuckets: []config.Bucket{
										{
											Name: "b2",
											Import: &config.ImportBucket{
												Bucket: "test-bucket",
											},
										},
									},
								},
							},
						},
					},
				},
				Buckets: []config.Bucket{
					{
						Name: "test-bucket",
					},
				},
			},
			content: config.DredgeFile{
				Workflows: []config.Workflow{
					{
						Name: "w1",
						Steps: []config.Step{
							{
								Browser: &config.BrowserStep{
									Url: "https://www.dredge.dev",
								},
							},
						},
					},
				},
				Buckets: []config.Bucket{
					{
						Name:      "b1",
						Workflows: []config.Workflow{},
					},
					{
						Name: "b2",
						Import: &config.ImportBucket{
							Source: config.SourcePath(importFile),
							Bucket: "test-bucket",
						},
					},
				},
			},
		},
		"add existing workflow": {
			df: config.DredgeFile{
				Workflows: []config.Workflow{
					{
						Name: "add-workflow",
						Steps: []config.Step{
							{
								EditDredgeFile: &config.EditDredgeFileStep{
									AddWorkflows: []config.Workflow{
										{
											Name: "w1",
											Import: &config.ImportWorkflow{
												Workflow: "add-workflow",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			errMsg: "workflow w1 already present",
		},
		"add existing bucket": {
			df: config.DredgeFile{
				Workflows: []config.Workflow{
					{
						Name: "add-workflow",
						Steps: []config.Step{
							{
								EditDredgeFile: &config.EditDredgeFileStep{
									AddBuckets: []config.Bucket{
										{
											Name: "b1",
											Import: &config.ImportBucket{
												Bucket: "b2",
											},
										},
									},
								},
							},
						},
					},
				},
				Buckets: []config.Bucket{
					{
						Name: "b2",
					},
				},
			},
			errMsg: "bucket b1 already present",
		},
		"add invalid workflow": {
			df: config.DredgeFile{
				Workflows: []config.Workflow{
					{
						Name: "add-workflow",
						Steps: []config.Step{
							{
								EditDredgeFile: &config.EditDredgeFileStep{
									AddWorkflows: []config.Workflow{
										{
											Name: "w2",
											Steps: []config.Step{
												{},
											},
											Import: &config.ImportWorkflow{
												Workflow: "add-workflow",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			errMsg: "workflow w2: contains both steps and an import",
		},
		"add variable": {
			df: config.DredgeFile{
				Workflows: []config.Workflow{
					{
						Name: "add-workflow",
						Steps: []config.Step{
							{
								EditDredgeFile: &config.EditDredgeFileStep{
									AddVariables: config.Variables{
										"hello": "world",
									},
								},
							},
						},
					},
				},
			},
			content: config.DredgeFile{
				Variables: config.Variables{
					"hello": "world",
				},
				Workflows: []config.Workflow{
					{
						Name: "w1",
						Steps: []config.Step{
							{
								Browser: &config.BrowserStep{
									Url: "https://www.dredge.dev",
								},
							},
						},
					},
				},
				Buckets: []config.Bucket{
					{
						Name:      "b1",
						Workflows: []config.Workflow{},
					},
				},
			},
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		err = ioutil.WriteFile(dredgeFile, d, 0644)
		assert.Nil(t, err)

		importContent, err := yaml.Marshal(test.df)
		assert.Nil(t, err)

		err = ioutil.WriteFile(importFile, importContent, 0644)
		assert.Nil(t, err)

		e, err := NewExec(config.SourcePath(dredgeFile), nil, nil)
		assert.Nil(t, err)

		de, err := e.Import(config.SourcePath(importFile))
		assert.Nil(t, err)

		w, err := de.GetWorkflow("", "add-workflow")
		assert.Nil(t, err)

		err = w.Execute()
		if test.errMsg == "" {
			assert.Nil(t, err)

			_, df, err := ReadDredgeFile(config.SourcePath(dredgeFile))
			assert.Nil(t, err)

			content, err := yaml.Marshal(df)
			assert.Nil(t, err)

			c, err := yaml.Marshal(test.content)
			assert.Nil(t, err)

			assert.Equal(t, string(c), string(content))
		} else {
			assert.Equal(t, test.errMsg, fmt.Sprint(err))
		}

		os.Remove(dredgeFile)
		os.Remove(importFile)
	}
}

func TestGetRootExec(t *testing.T) {
	root := &DredgeExec{}

	tests := map[string]struct {
		e *DredgeExec
	}{
		"root": {
			e: root,
		},
		"1 level": {
			e: &DredgeExec{
				Parent: root,
			},
		},
		"2 levels": {
			e: &DredgeExec{
				Parent: &DredgeExec{
					Parent: root,
				},
			},
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		assert.Equal(t, root, test.e.getRootExec())
	}
}
