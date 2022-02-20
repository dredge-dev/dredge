package config

import (
	"io/ioutil"
	"reflect"

	"gopkg.in/yaml.v3"
)

const DEFAULT_HOME = "/home"

type DredgeFile struct {
	Env       Env
	Workflows []Workflow
	Buckets   []Bucket
}

type Env struct {
	Variables map[string]string
	Runtimes  []Runtime
}

type Runtime struct {
	Name  string
	Type  string
	Image string
	Home  *string
	Cache []string
	Ports []string
}

type Bucket struct {
	Name        string
	Description string
	Workflows   []Workflow
}

type Workflow struct {
	Name        string
	Description string
	Inputs      map[string]string
	Steps       []Step
}

type Step struct {
	Name     *string
	Shell    *ShellStep
	Template *TemplateStep
	Browser  *string
}

type ShellStep struct {
	Cmd     string
	Runtime string
}

type TemplateStep struct {
	Input string
	Dest  TemplateString
}

type TemplateString string

func ReadDredgeFile(filename string) (*DredgeFile, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	dredgeFile := &DredgeFile{}
	err = yaml.Unmarshal(buf, dredgeFile)
	if err != nil {
		return nil, err
	}

	return dredgeFile, nil
}

func (dredgeFile *DredgeFile) GetWorkflow(workflowName string) *Workflow {
	for _, w := range dredgeFile.Workflows {
		if w.Name == workflowName {
			return &w
		}
	}
	return nil
}

func (dredgeFile *DredgeFile) GetWorkflowInBucket(bucketName, workflowName string) *Workflow {
	for _, b := range dredgeFile.Buckets {
		if b.Name == bucketName {
			for _, w := range b.Workflows {
				if w.Name == workflowName {
					return &w
				}
			}
		}
	}
	return nil
}

func (dredgeFile *DredgeFile) GetBucket(bucketName string) *Bucket {
	for _, b := range dredgeFile.Buckets {
		if b.Name == bucketName {
			return &b
		}
	}
	return nil
}

func (s Step) Validate() bool {
	numFields := 0

	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		if v.Type().Field(i).Name != "Name" {
			if !v.Field(i).IsNil() {
				numFields += 1
			}
		}
	}

	return numFields == 1
}

func (r Runtime) GetHome() string {
	if r.Home == nil {
		return DEFAULT_HOME
	}
	return *r.Home
}
