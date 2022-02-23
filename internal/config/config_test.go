package config

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestGetDredgeFile(t *testing.T) {
	df := DredgeFile{
		Buckets: []Bucket{
			{
				Name: "b1",
			},
		},
	}

	tmpFile := fmt.Sprintf("./tmp-test-%d", rand.Intn(100000))
	defer os.Remove(tmpFile)

	data, err := yaml.Marshal(df)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(tmpFile, data, 0644)
	if err != nil {
		panic(err)
	}

	tests := map[string]struct {
		source     string
		dredgeFile *DredgeFile
		errMsg     string
	}{
		"local file": {
			source:     tmpFile,
			dredgeFile: &df,
			errMsg:     "",
		},
		"file not found": {
			source:     "./non-existing-file",
			dredgeFile: nil,
			errMsg:     "Error while parsing ./non-existing-file: open ./non-existing-file: no such file or directory",
		},
		"unsupported": {
			source:     "/hello",
			dredgeFile: nil,
			errMsg:     "Sources should start with ./",
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		d, err := GetDredgeFile(test.source)
		if test.dredgeFile == nil {
			assert.Nil(t, d)
		} else {
			assert.Equal(t, test.dredgeFile.Buckets[0].Name, d.Buckets[0].Name)
		}
		if test.errMsg == "" {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, test.errMsg, fmt.Sprint(err))
		}
	}
}

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
			step:     Step{Shell: &ShellStep{Cmd: "cmd", Runtime: "runtime"}, Template: &TemplateStep{Input: "input", Dest: "dst"}},
			errorMsg: "step  contains more than 1 action",
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
	w3 := Workflow{Name: "third"}
	w4 := Workflow{Name: "fourth"}
	dredgeFile := &DredgeFile{
		Workflows: []Workflow{
			w1,
			w2,
		},
		Buckets: []Bucket{
			{
				Name: "b1",
				Workflows: []Workflow{
					w3,
				},
			},
			{
				Name: "b2",
				Workflows: []Workflow{
					w4,
				},
			},
		},
	}

	tests := map[string]struct {
		bucketName   string
		workflowName string
		workflow     *Workflow
		errMsg       string
	}{
		"first": {
			bucketName:   "",
			workflowName: "first",
			workflow:     &w1,
			errMsg:       "",
		},
		"second": {
			bucketName:   "",
			workflowName: "second",
			workflow:     &w2,
			errMsg:       "",
		},
		"third is in default bucket": {
			bucketName:   "",
			workflowName: "third",
			workflow:     nil,
			errMsg:       "Could not find workflow /third",
		},
		"third should be in b1": {
			bucketName:   "b1",
			workflowName: "third",
			workflow:     &w3,
			errMsg:       "",
		},
		"third should not be in b2": {
			bucketName:   "b2",
			workflowName: "third",
			workflow:     nil,
			errMsg:       "Could not find workflow b2/third",
		},
		"fourth should be in b2": {
			bucketName:   "b2",
			workflowName: "fourth",
			workflow:     &w4,
			errMsg:       "",
		},
		"fifth does not exist in default": {
			bucketName:   "",
			workflowName: "fifth",
			workflow:     nil,
			errMsg:       "Could not find workflow /fifth",
		},
		"fifth does not exist in b1": {
			bucketName:   "b1",
			workflowName: "fifth",
			workflow:     nil,
			errMsg:       "Could not find workflow b1/fifth",
		},
		"bucket not found": {
			bucketName:   "b3",
			workflowName: "first",
			workflow:     nil,
			errMsg:       "Could not find workflow b3/first",
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		w, err := dredgeFile.GetWorkflow(test.bucketName, test.workflowName)
		assert.Equal(t, test.workflow, w)
		if test.errMsg == "" {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, test.errMsg, fmt.Sprint(err))
		}
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
		errMsg     string
	}{
		"b1": {
			bucketName: "b1",
			bucket:     &b1,
			errMsg:     "",
		},
		"b2": {
			bucketName: "b2",
			bucket:     &b2,
			errMsg:     "",
		},
		"b3 not found": {
			bucketName: "b3",
			bucket:     nil,
			errMsg:     "Could not find bucket b3",
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		b, err := dredgeFile.GetBucket(test.bucketName)
		assert.Equal(t, test.bucket, b)
		if test.errMsg == "" {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, test.errMsg, fmt.Sprint(err))
		}
	}
}

func TestMergeSources(t *testing.T) {
	tests := map[string]struct {
		parent string
		child  string
		result string
	}{
		"parent without dir, child without dir": {
			parent: "./test.Dredgefile",
			child:  "./second.Dredgefile",
			result: "./second.Dredgefile",
		},
		"parent with dir, child without dir": {
			parent: "./parent/test.Dredgefile",
			child:  "./second.Dredgefile",
			result: "./parent/second.Dredgefile",
		},
		"parent without dir, child with dir": {
			parent: "./test.Dredgefile",
			child:  "./child/second.Dredgefile",
			result: "./child/second.Dredgefile",
		},
		"parent with dir, child with dir": {
			parent: "./parent/test.Dredgefile",
			child:  "./child/second.Dredgefile",
			result: "./parent/child/second.Dredgefile",
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		assert.Equal(t, test.result, mergeSources(test.parent, test.child))
	}
}
