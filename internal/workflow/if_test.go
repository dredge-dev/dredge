package workflow

import (
	"fmt"
	"os"
	"testing"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/exec"
	"github.com/stretchr/testify/assert"
)

const TEST_FILE = "tmp-dredge-if-test"

func getTouchWorkflow(cond string) *exec.Workflow {
	w := &exec.Workflow{
		Exec: exec.EmptyExec(""),
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
	}
	w.Exec.Env["RUN"] = "true"
	w.Exec.Env["DONT"] = "false"
	return w
}

func TestExecuteIfStep(t *testing.T) {
	defer os.Remove(TEST_FILE)

	tests := map[string]struct {
		workflow *exec.Workflow
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
			errMsg:   "Failed to parse template: template: :1: unexpected \"}\" in operand",
			exists:   false,
		},
		"error in steps": {
			workflow: &exec.Workflow{
				Exec: exec.EmptyExec(""),
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
			},
			errMsg: "No execution found for step ",
			exists: false,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		err := ExecuteWorkflow(test.workflow)
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
