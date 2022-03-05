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

func TestNewExec(t *testing.T) {
	df := config.DredgeFile{
		Buckets: []config.Bucket{
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
		source     config.SourcePath
		dredgeFile *config.DredgeFile
		errMsg     string
	}{
		"local file": {
			source:     config.SourcePath(tmpFile),
			dredgeFile: &df,
			errMsg:     "",
		},
		"file not found": {
			source:     "./non-existing-file",
			dredgeFile: nil,
			errMsg:     "stat ./non-existing-file: no such file or directory",
		},
		"unsupported": {
			source:     "/hello",
			dredgeFile: nil,
			errMsg:     "sources should start with ./",
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		d, err := NewExec(test.source)
		if test.dredgeFile == nil {
			assert.Nil(t, d)
		} else {
			buckets, err := d.GetBuckets()
			assert.Nil(t, err)
			assert.Equal(t, test.dredgeFile.Buckets[0].Name, buckets[0].Name)
		}
		if test.errMsg == "" {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, test.errMsg, fmt.Sprint(err))
		}
	}
}

func TestGetWorkflow(t *testing.T) {
	de := &DredgeExec{
		Source: "./Dredgefile",
		Env:    NewEnv(),
		DredgeFile: &config.DredgeFile{
			Workflows: []config.Workflow{
				{Name: "first"},
				{Name: "second"},
			},
			Buckets: []config.Bucket{
				{
					Name: "b1",
					Workflows: []config.Workflow{
						{Name: "third"},
					},
				},
				{
					Name: "b2",
					Workflows: []config.Workflow{
						{Name: "fourth"},
					},
				},
			},
		},
	}

	tests := map[string]struct {
		bucketName   string
		workflowName string
		errMsg       string
	}{
		"first": {
			bucketName:   "",
			workflowName: "first",
			errMsg:       "",
		},
		"second": {
			bucketName:   "",
			workflowName: "second",
			errMsg:       "",
		},
		"third is in default bucket": {
			bucketName:   "",
			workflowName: "third",
			errMsg:       "Could not find workflow /third",
		},
		"third should be in b1": {
			bucketName:   "b1",
			workflowName: "third",
			errMsg:       "",
		},
		"third should not be in b2": {
			bucketName:   "b2",
			workflowName: "third",
			errMsg:       "Could not find workflow b2/third",
		},
		"fourth should be in b2": {
			bucketName:   "b2",
			workflowName: "fourth",
			errMsg:       "",
		},
		"fifth does not exist in default": {
			bucketName:   "",
			workflowName: "fifth",
			errMsg:       "Could not find workflow /fifth",
		},
		"fifth does not exist in b1": {
			bucketName:   "b1",
			workflowName: "fifth",
			errMsg:       "Could not find workflow b1/fifth",
		},
		"bucket not found": {
			bucketName:   "b3",
			workflowName: "first",
			errMsg:       "Could not find workflow b3/first",
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		w, err := de.GetWorkflow(test.bucketName, test.workflowName)
		if test.errMsg == "" {
			assert.Nil(t, err)
			assert.Equal(t, test.workflowName, w.Name)
		} else {
			assert.Equal(t, test.errMsg, fmt.Sprint(err))
		}
	}
}

func TestGetBucket(t *testing.T) {
	de := &DredgeExec{
		Source: "./Dredgefile",
		Env:    NewEnv(),
		DredgeFile: &config.DredgeFile{
			Buckets: []config.Bucket{
				{Name: "b1"},
				{Name: "b2"},
			},
		},
	}
	tests := map[string]struct {
		bucketName string
		errMsg     string
	}{
		"b1": {
			bucketName: "b1",
			errMsg:     "",
		},
		"b2": {
			bucketName: "b2",
			errMsg:     "",
		},
		"b3 not found": {
			bucketName: "b3",
			errMsg:     "Could not find bucket b3",
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		b, err := de.GetBucket(test.bucketName)
		if test.errMsg == "" {
			assert.Nil(t, err)
			assert.Equal(t, test.bucketName, b.Name)
		} else {
			assert.Equal(t, test.errMsg, fmt.Sprint(err))
		}
	}
}

func TestImportWorkflowSameDedgeFile(t *testing.T) {
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
									Cmd: "touch tmp",
								},
							},
						},
					},
				},
			},
		},
	}

	de := &DredgeExec{
		Source:     "./Dredgefile",
		DredgeFile: dredgeFile,
		Env:        NewEnv(),
	}

	w, err := de.GetWorkflow("", "workflow")
	assert.Nil(t, err)
	assert.Equal(t, "workflow", w.Name)
	assert.Equal(t, "perform work", w.Description)
	assert.Equal(t, "touch tmp", w.Steps[0].Shell.Cmd)
}

func TestImport(t *testing.T) {
	remoteDredgeFile := config.DredgeFile{
		Buckets: []config.Bucket{
			{
				Name:        "b1",
				Description: "a bucket of workflows",
				Workflows: []config.Workflow{
					{
						Name:        "w1",
						Description: "perform work",
						Steps: []config.Step{
							{
								Shell: &config.ShellStep{
									Cmd: "touch tmp",
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

	df := &config.DredgeFile{
		Workflows: []config.Workflow{
			{
				Name: "workflow",
				Import: &config.ImportWorkflow{
					Source:   config.SourcePath(remoteDredgeFilePath),
					Bucket:   "b1",
					Workflow: "w1",
				},
			},
		},
		Buckets: []config.Bucket{
			{
				Name: "bucket",
				Import: &config.ImportBucket{
					Source: config.SourcePath(remoteDredgeFilePath),
					Bucket: "b1",
				},
			},
		},
	}

	de := &DredgeExec{
		Source:     "./Dredgefile",
		DredgeFile: df,
		Env:        NewEnv(),
	}

	w, err := de.GetWorkflow("", "workflow")
	assert.Nil(t, err)
	assert.Equal(t, "workflow", w.Name)
	assert.Equal(t, "perform work", w.Description)
	assert.Equal(t, "touch tmp", w.Steps[0].Shell.Cmd)

	b, err := de.GetBucket("bucket")
	assert.Nil(t, err)
	assert.Equal(t, "bucket", b.Name)
	assert.Equal(t, "a bucket of workflows", b.Description)
}
