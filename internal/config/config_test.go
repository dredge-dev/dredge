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

func TestGetWorkflow(t *testing.T) {
	w1 := Workflow{Name: "first"}
	w2 := Workflow{Name: "second"}
	dredgeFile := &DredgeFile{
		Workflows: []Workflow{
			w1,
			w2,
		},
	}

	tests := map[string]struct {
		name     string
		workflow *Workflow
	}{
		"first": {
			name:     "first",
			workflow: &w1,
		},
		"second": {
			name:     "second",
			workflow: &w2,
		},
		"not found": {
			name:     "third",
			workflow: nil,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		w := dredgeFile.GetWorkflow(test.name)
		assert.Equal(t, test.workflow, w)
	}
}

func TestGetWorkflowInBucket(t *testing.T) {
	w1 := Workflow{Name: "first"}
	w2 := Workflow{Name: "second"}
	w3 := Workflow{Name: "third"}
	dredgeFile := &DredgeFile{
		Buckets: []Bucket{
			{
				Name: "b1",
				Workflows: []Workflow{
					w1,
					w2,
				},
			},
			{
				Name: "b2",
				Workflows: []Workflow{
					w3,
				},
			},
		},
	}

	tests := map[string]struct {
		bucketName   string
		workflowName string
		workflow     *Workflow
	}{
		"first should be in b1": {
			bucketName:   "b1",
			workflowName: "first",
			workflow:     &w1,
		},
		"first should not be in b2": {
			bucketName:   "b2",
			workflowName: "first",
			workflow:     nil,
		},
		"second": {
			bucketName:   "b1",
			workflowName: "second",
			workflow:     &w2,
		},
		"third": {
			bucketName:   "b2",
			workflowName: "third",
			workflow:     &w3,
		},
		"bucket not found": {
			bucketName:   "b3",
			workflowName: "first",
			workflow:     nil,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		w := dredgeFile.GetWorkflowInBucket(test.bucketName, test.workflowName)
		assert.Equal(t, test.workflow, w)
	}
}

func TestGetBucket(t *testing.T) {
	b1 := Bucket{
		Name: "b1",
	}
	b2 := Bucket{
		Name: "b2",
	}
	dredgeFile := &DredgeFile{
		Buckets: []Bucket{
			b1,
			b2,
		},
	}
	tests := map[string]struct {
		bucketName string
		bucket     *Bucket
	}{
		"b1": {
			bucketName: "b1",
			bucket:     &b1,
		},
		"b2": {
			bucketName: "b2",
			bucket:     &b2,
		},
		"b3 not found": {
			bucketName: "b3",
			bucket:     nil,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		b := dredgeFile.GetBucket(test.bucketName)
		assert.Equal(t, test.bucket, b)
	}
}
