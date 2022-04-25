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
		"runtime validation": {
			dredgeFile: &DredgeFile{
				Runtimes: []Runtime{
					{},
				},
			},
			errorMsg: "name field is required for runtime",
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

func TestRuntimeValidate(t *testing.T) {
	tests := map[string]struct {
		runtime  Runtime
		errorMsg string
	}{
		"invalid runtime type": {
			runtime: Runtime{
				Name: "c1",
				Type: "cool",
			},
			errorMsg: "unknown runtime type: cool (valid options are native, container)",
		},
		"native runtime type": {
			runtime: Runtime{
				Name: "n",
				Type: "native",
			},
			errorMsg: "",
		},
		"runtime without name": {
			runtime: Runtime{
				Type: "native",
			},
			errorMsg: "name field is required for runtime",
		},
		"container runtime type": {
			runtime: Runtime{
				Name:  "c",
				Type:  "container",
				Image: "my-image",
			},
			errorMsg: "",
		},
		"container missing image": {
			runtime: Runtime{
				Name: "c",
				Type: "container",
			},
			errorMsg: "image field is required for container runtimes",
		},
		"native with container fields": {
			runtime: Runtime{
				Name:  "n",
				Type:  "native",
				Image: "out-of-place",
			},
			errorMsg: "image, home, cache, global_cache and ports fields are only applicable to container runtimes",
		},
	}
	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		err := test.runtime.Validate()
		if test.errorMsg == "" {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, test.errorMsg, fmt.Sprint(err))
		}
	}
}

func TestInputValidate(t *testing.T) {
	tests := map[string]struct {
		input    Input
		errorMsg string
	}{
		"missing name": {
			input: Input{
				Name: "",
			},
			errorMsg: "name field is required on inputs",
		},
		"invalid type": {
			input: Input{
				Name: "test",
				Type: "invalid",
			},
			errorMsg: "input test: unknown input type: invalid (valid options are: text, select)",
		},
		"simple input": {
			input: Input{
				Name: "inputName",
			},
			errorMsg: "",
		},
		"default input": {
			input: Input{
				Name:         "test",
				Description:  "some input",
				DefaultValue: "world",
			},
			errorMsg: "",
		},
		"text input": {
			input: Input{
				Name:         "test",
				Type:         "text",
				Description:  "some input",
				DefaultValue: "world",
			},
			errorMsg: "",
		},
		"default input with values": {
			input: Input{
				Name: "test",
				Values: []string{
					"hello", "world",
				},
			},
			errorMsg: "input test: values for input can only be provided for the select type",
		},
		"text input with values": {
			input: Input{
				Name: "test",
				Type: "text",
				Values: []string{
					"hello", "world",
				},
			},
			errorMsg: "input test: values for input can only be provided for the select type",
		},
		"select input": {
			input: Input{
				Name:        "select-input",
				Type:        "select",
				Description: "select an input",
				Values: []string{
					"value1", "value2", "value3",
				},
			},
			errorMsg: "",
		},
		"select input without values": {
			input: Input{
				Name:        "select-input",
				Type:        "select",
				Description: "select an input",
				Values:      []string{},
			},
			errorMsg: "input select-input: no values are provided, values are required for the select type",
		},
		"select input with default value": {
			input: Input{
				Name:        "select-input",
				Type:        "select",
				Description: "select an input",
				Values: []string{
					"value1", "value2",
				},
				DefaultValue: "value1",
			},
			errorMsg: "input select-input: default value can only be provided for the text type",
		},
	}
	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		err := test.input.Validate()
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
		"template": {
			step: Step{Template: &TemplateStep{
				Source: "file",
				Dest:   "test",
				Insert: &Insert{
					Section: "import",
				},
			}},
			errorMsg: "",
		},
		"template without dest": {
			step:     Step{Template: &TemplateStep{}},
			errorMsg: "dest field is required for template",
		},
		"template with both input and source": {
			step: Step{Template: &TemplateStep{
				Source: "file",
				Input:  "hello",
			}},
			errorMsg: "either input or source should be set for template",
		},
		"template with invalid insert placement": {
			step: Step{Template: &TemplateStep{
				Source: "file",
				Dest:   "test",
				Insert: &Insert{
					Section:   "import",
					Placement: "middle",
				},
			}},
			errorMsg: "unknown placement in insert: middle (valid options are: begin, end, unique)",
		},
		"browser": {
			step: Step{Browser: &BrowserStep{
				Url: "https://dredge.dev/",
			}},
			errorMsg: "",
		},
		"invalid browser": {
			step:     Step{Browser: &BrowserStep{}},
			errorMsg: "url field is required for browser",
		},
		"if": {
			step: Step{
				If: &IfStep{
					Cond: "{{ .VALID }}",
					Steps: []Step{
						{Browser: &BrowserStep{Url: "https://www.google.com"}},
					},
				},
			},
			errorMsg: "",
		},
		"if without cond": {
			step: Step{
				If: &IfStep{
					Steps: []Step{
						{Browser: &BrowserStep{Url: "https://www.google.com"}},
					},
				},
			},
			errorMsg: "cond field is required for if",
		},
		"if without steps": {
			step: Step{
				If: &IfStep{
					Cond:  "{{ .VALID }}",
					Steps: []Step{},
				},
			},
			errorMsg: "1 or more steps are required for if",
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
