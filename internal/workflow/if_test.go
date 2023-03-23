package workflow

import (
	"fmt"
	"os"
	"testing"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/stretchr/testify/assert"
)

const TEST_FILE = "tmp-dredge-if-test"

func getTouchWorkflow(cond string) *Workflow {
	c := &CallbacksMock{
		Env: map[string]interface{}{
			"RUN":  "true",
			"DONT": "false",
		},
	}
	w := &Workflow{
		Name: "workflow",
		Steps: []config.Step{
			{
				If: &config.IfStep{
					Cond: cond,
					Steps: []config.Step{
						{
							Shell: &config.ShellStep{
								Cmd: fmt.Sprintf("touch %s", TEST_FILE),
							},
						},
					},
				},
			},
		},
		Callbacks: c,
	}
	return w
}

func TestExecuteIfStep(t *testing.T) {
	defer os.Remove(TEST_FILE)

	tests := map[string]struct {
		workflow *Workflow
		errMsg   string
		exists   bool
	}{
		"if true": {
			workflow: getTouchWorkflow("{{ .RUN }}"),
			errMsg:   "",
			exists:   true,
		},
		"if false": {
			workflow: getTouchWorkflow("{{ .DONT }}"),
			errMsg:   "",
			exists:   false,
		},
		"if other": {
			workflow: getTouchWorkflow("junk"),
			errMsg:   "",
			exists:   false,
		},
		"bad template": {
			workflow: getTouchWorkflow("{{ .BAD }"),
			errMsg:   "failed to parse template: template: :1: unexpected \"}\" in operand",
			exists:   false,
		},
		"error in steps": {
			workflow: &Workflow{
				Name: "workflow",
				Steps: []config.Step{
					{
						If: &config.IfStep{
							Cond: "true",
							Steps: []config.Step{
								{},
							},
						},
					},
				},
				Callbacks: &CallbacksMock{},
			},
			errMsg: "no execution found for step ",
			exists: false,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		err := test.workflow.Execute()
		if test.errMsg == "" {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, test.errMsg, fmt.Sprint(err))
		}
		_, err = os.Stat(TEST_FILE)
		if test.exists {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
		os.Remove(TEST_FILE)
	}
}
