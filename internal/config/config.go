package config

import (
	"io/ioutil"

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
	Import      *ImportBucket
}

type ImportBucket struct {
	Source string
	Bucket string
}

type Workflow struct {
	Name        string
	Description string
	Inputs      map[string]string
	Steps       []Step
	Import      *ImportWorkflow
}

type ImportWorkflow struct {
	Source   string
	Bucket   string
	Workflow string
}

type Step struct {
	Name     string
	Shell    *ShellStep
	Template *TemplateStep
	Browser  *BrowserStep
}

type ShellStep struct {
	Cmd     string
	Runtime string
}

type TemplateStep struct {
	Input string
	Dest  string
}

type BrowserStep struct {
	Url string
}

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

	err = dredgeFile.Validate()
	if err != nil {
		return nil, err
	}
	return dredgeFile, nil
}

func (r Runtime) GetHome() string {
	if r.Home == nil {
		return DEFAULT_HOME
	}
	return *r.Home
}
