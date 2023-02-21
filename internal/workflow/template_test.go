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
)

func TestExecuteTemplate(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("drg-%d", rand.Intn(100000)))
	defer os.Remove(tmpFile)

	step := config.TemplateStep{
		Input: "Hello {{ .test }}",
		Dest:  tmpFile,
	}
	workflow := &Workflow{
		Name:        "workflow",
		Description: "My workflow",
		Inputs: []config.Input{
			{
				Name: "test",
			},
		},
		Steps: []config.Step{
			{
				Name:     "",
				Template: &step,
			},
		},
		Callbacks: &CallbacksMock{
			Env: map[string]interface{}{
				"test": "value",
			},
		},
	}

	err := workflow.Execute()
	assert.Nil(t, err)

	content, err := ioutil.ReadFile(tmpFile)
	assert.Nil(t, err)
	assert.Equal(t, "Hello value", string(content))
}

func TestExecuteTemplateFromSource(t *testing.T) {
	dstFile := filepath.Join(os.TempDir(), fmt.Sprintf("drg-%d", rand.Intn(100000)))
	defer os.Remove(dstFile)

	templateFile := "./test-execute-template"
	err := ioutil.WriteFile(templateFile, []byte("Hello {{ .test }}"), 0644)
	defer os.Remove(templateFile)
	assert.Nil(t, err)

	workflow := &Workflow{
		Name:        "workflow",
		Description: "My workflow",
		Inputs: []config.Input{
			{
				Name: "test",
			},
		},
		Steps: []config.Step{
			{
				Name: "",
				Template: &config.TemplateStep{
					Source: config.SourcePath(templateFile),
					Dest:   dstFile,
				},
			},
		},
		Callbacks: &CallbacksMock{
			Env: map[string]interface{}{
				"test": "value",
			},
		},
	}

	err = workflow.Execute()
	assert.Nil(t, err)

	content, err := ioutil.ReadFile(dstFile)
	assert.Nil(t, err)
	assert.Equal(t, "Hello value", string(content))
}

func TestInsert(t *testing.T) {
	dstPath := filepath.Join(os.TempDir(), fmt.Sprintf("drg-%d", rand.Intn(100000)))

	tests := map[string]struct {
		preContent  string
		insert      *config.Insert
		text        string
		dest        string
		postContent string
		errorMsg    string
	}{
		"no insert": {
			preContent:  "something to overwrite",
			text:        "hello",
			dest:        dstPath,
			postContent: "hello",
			errorMsg:    "",
		},
		"new file begin": {
			preContent:  "",
			insert:      &config.Insert{Placement: "begin"},
			text:        "hello",
			dest:        dstPath,
			postContent: "hello",
			errorMsg:    "",
		},
		"new file end": {
			preContent:  "",
			insert:      &config.Insert{Placement: "end"},
			text:        "hello",
			dest:        dstPath,
			postContent: "hello",
			errorMsg:    "",
		},
		"prefix content": {
			preContent:  "world",
			insert:      &config.Insert{Placement: "begin"},
			text:        "hello",
			dest:        dstPath,
			postContent: "hello\nworld",
			errorMsg:    "",
		},
		"suffix content": {
			preContent:  "hello",
			insert:      &config.Insert{Placement: "end"},
			text:        "world",
			dest:        dstPath,
			postContent: "hello\nworld",
			errorMsg:    "",
		},
		"unique content new": {
			preContent:  "hello\nworld",
			insert:      &config.Insert{Placement: "unique"},
			text:        "new",
			dest:        dstPath,
			postContent: "hello\nworld\nnew",
			errorMsg:    "",
		},
		"unique content exists": {
			preContent:  "hello\nworld",
			insert:      &config.Insert{Placement: "unique"},
			text:        "hello",
			dest:        dstPath,
			postContent: "hello\nworld",
			errorMsg:    "",
		},
		"unique content exists 2": {
			preContent:  "hello\nworld",
			insert:      &config.Insert{Placement: "unique"},
			text:        "world",
			dest:        dstPath,
			postContent: "hello\nworld",
			errorMsg:    "",
		},
		"default to suffix": {
			preContent:  "hello",
			insert:      &config.Insert{},
			text:        "world",
			dest:        dstPath,
			postContent: "hello\nworld",
			errorMsg:    "",
		},
		"go import": {
			preContent:  "package main\nimport \"fmt\"\nfunc main() {\n}\n",
			insert:      &config.Insert{Section: "import"},
			text:        "\"testing\"",
			dest:        dstPath + ".go",
			postContent: "package main\n\nimport (\n\t\"fmt\"\n\t\"testing\"\n)\n\nfunc main() {\n}",
			errorMsg:    "",
		},
		"invalid extension": {
			preContent: "hello",
			insert:     &config.Insert{Section: "import"},
			text:       " world",
			dest:       dstPath + ".java",
			errorMsg:   "unsupported extension java for insert (valid values: go)",
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		if test.preContent != "" {
			err := ioutil.WriteFile(test.dest, []byte(test.preContent), 0644)
			assert.Nil(t, err)
		}
		err := insert(test.insert, test.text, test.dest)
		if test.errorMsg == "" {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, test.errorMsg, fmt.Sprint(err))
		}
		if test.postContent != "" {
			bytes, err := ioutil.ReadFile(test.dest)
			assert.Nil(t, err)
			assert.Equal(t, test.postContent, string(bytes))
		}
		os.Remove(test.dest)
	}
}
