package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	tests := map[string]struct {
		dredgeFile *DredgeFile
		errorMsg   string
	}{
		"valid bucket": {
			dredgeFile: &DredgeFile{
				Buckets: []Bucket{
					{
						Name: "b1",
					},
				},
			},
			errorMsg: "",
		},
		"bucket without name": {
			dredgeFile: &DredgeFile{
				Buckets: []Bucket{
					{},
				},
			},
			errorMsg: "name field is required for bucket",
		},
		"bucket with import": {
			dredgeFile: &DredgeFile{
				Buckets: []Bucket{
					{
						Name: "b1",
						Import: &ImportBucket{
							Bucket: "b2",
						},
					},
				},
			},
			errorMsg: "",
		},
		"bucket with invalid import": {
			dredgeFile: &DredgeFile{
				Buckets: []Bucket{
					{
						Name:   "b1",
						Import: &ImportBucket{},
					},
				},
			},
			errorMsg: "bucket b1: bucket field is required for import",
		},
		"bucket with import and workflow": {
			dredgeFile: &DredgeFile{
				Buckets: []Bucket{
					{
						Name: "b1",
						Import: &ImportBucket{
							Bucket: "b2",
						},
						Workflows: []Workflow{
							{
								Name: "workflow",
							},
						},
					},
				},
			},
			errorMsg: "bucket b1: contains both workflows and an import",
		},
		"bucket with invalid workflow": {
			dredgeFile: &DredgeFile{
				Buckets: []Bucket{
					{
						Name: "b1",
						Workflows: []Workflow{
							{},
						},
					},
				},
			},
			errorMsg: "bucket b1: name field is required for workflow",
		},
		"workflow without name": {
			dredgeFile: &DredgeFile{
				Workflows: []Workflow{
					{},
				},
			},
			errorMsg: "name field is required for workflow",
		},
		"workflow with import": {
			dredgeFile: &DredgeFile{
				Workflows: []Workflow{
					{
						Name: "workflow",
						Import: &ImportWorkflow{
							Workflow: "test",
						},
					},
				},
			},
			errorMsg: "",
		},
		"workflow with step": {
			dredgeFile: &DredgeFile{
				Workflows: []Workflow{
					{
						Name: "workflow",
						Steps: []Step{
							{
								Shell: &ShellStep{
									Cmd: "test",
								},
							},
						},
					},
				},
			},
			errorMsg: "",
		},
		"workflow no steps no import": {
			dredgeFile: &DredgeFile{
				Workflows: []Workflow{
					{
						Name: "w1",
					},
				},
			},
			errorMsg: "workflow w1: no steps or import defined",
		},
		"workflow with steps and import": {
			dredgeFile: &DredgeFile{
				Workflows: []Workflow{
					{
						Name: "w1",
						Import: &ImportWorkflow{
							Workflow: "w2",
						},
						Steps: []Step{
							{
								Shell: &ShellStep{
									Cmd: "test",
								},
							},
						},
					},
				},
			},
			errorMsg: "workflow w1: contains both steps and an import",
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		err := test.dredgeFile.Validate()
		if test.errorMsg == "" {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, test.errorMsg, fmt.Sprint(err))
		}
	}
}

func TestStepValidate(t *testing.T) {
	tests := map[string]struct {
		step     Step
		errorMsg string
	}{
		"no children": {
			step:     Step{},
			errorMsg: "step  does not contain an action",
		},
		"only a name": {
			step:     Step{Name: "s1"},
			errorMsg: "step s1 does not contain an action",
		},
		"1 child": {
			step:     Step{Shell: &ShellStep{Cmd: "cmd", Runtime: "runtime"}},
			errorMsg: "",
		},
		"1 child and a name": {
			step:     Step{Name: "step", Shell: &ShellStep{Cmd: "cmd", Runtime: "runtime"}},
			errorMsg: "",
		},
		"2 children": {
			step:     Step{Shell: &ShellStep{Cmd: "cmd", Runtime: "runtime"}, Template: &TemplateStep{Input: "input", Dest: "dst"}, Browser: &BrowserStep{Url: "url"}},
			errorMsg: "step  contains more than 1 action",
		},
		"invalid shell": {
			step:     Step{Shell: &ShellStep{}},
			errorMsg: "cmd field is required for shell",
		},
		"invalid template": {
			step:     Step{Template: &TemplateStep{}},
			errorMsg: "input and dest fields are required for template",
		},
		"invalid browser": {
			step:     Step{Browser: &BrowserStep{}},
			errorMsg: "url field is required for browser",
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		err := test.step.Validate()
		if test.errorMsg == "" {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, test.errorMsg, fmt.Sprint(err))
		}
	}
}
