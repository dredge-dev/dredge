package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStepValidate(t *testing.T) {
	stepName := "step"

	tests := map[string]struct {
		step  Step
		valid bool
	}{
		"no children": {
			step:  Step{},
			valid: false,
		},
		"only a name": {
			step:  Step{Name: &stepName},
			valid: false,
		},
		"1 child": {
			step:  Step{Shell: &ShellStep{Cmd: "cmd", Runtime: "runtime"}},
			valid: true,
		},
		"1 child and a name": {
			step:  Step{Name: &stepName, Shell: &ShellStep{Cmd: "cmd", Runtime: "runtime"}},
			valid: true,
		},
		"2 children": {
			step:  Step{Shell: &ShellStep{Cmd: "cmd", Runtime: "runtime"}, Template: &TemplateStep{Input: "input", Dest: "dst"}},
			valid: false,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		assert.Equal(t, test.valid, test.step.Validate())
	}
}

func TestGetHome(t *testing.T) {
	customHome := "/my-home"

	tests := map[string]struct {
		runtime Runtime
		home    string
	}{
		"default home": {
			runtime: Runtime{Home: nil},
			home:    DEFAULT_HOME,
		},
		"custom home": {
			runtime: Runtime{Home: &customHome},
			home:    customHome,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		home := test.runtime.GetHome()
		assert.Equal(t, test.home, home)
	}
}
