package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

const DEFAULT_HOME = "/home"

type DredgeFile struct {
	Env       Env        `yaml:",omitempty"`
	Workflows []Workflow `yaml:",omitempty"`
	Buckets   []Bucket   `yaml:",omitempty"`
}

type Env struct {
	Variables map[string]string `yaml:",omitempty"`
	Runtimes  []Runtime         `yaml:",omitempty"`
}

type Runtime struct {
	Name  string
	Type  string
	Image string   `yaml:",omitempty"`
	Home  *string  `yaml:",omitempty"`
	Cache []string `yaml:",omitempty"`
	Ports []string `yaml:",omitempty"`
}

type Bucket struct {
	Name        string
	Description string        `yaml:",omitempty"`
	Workflows   []Workflow    `yaml:",omitempty"`
	Import      *ImportBucket `yaml:",omitempty"`
}

type ImportBucket struct {
	Source string
	Bucket string
}

type Workflow struct {
	Name        string
	Description string            `yaml:",omitempty"`
	Inputs      map[string]string `yaml:",omitempty"`
	Steps       []Step            `yaml:",omitempty"`
	Import      *ImportWorkflow   `yaml:",omitempty"`
}

type ImportWorkflow struct {
	Source   string
	Bucket   string
	Workflow string
}

type Step struct {
	Name           string              `yaml:",omitempty"`
	Shell          *ShellStep          `yaml:",omitempty"`
	Template       *TemplateStep       `yaml:",omitempty"`
	Browser        *BrowserStep        `yaml:",omitempty"`
	EditDredgeFile *EditDredgeFileStep `yaml:"edit_dredgefile,omitempty"`
}

type ShellStep struct {
	Cmd     string
	Runtime string `yaml:",omitempty"`
}

type TemplateStep struct {
	Input string
	Dest  string
}

type BrowserStep struct {
	Url string
}

type EditDredgeFileStep struct {
	AddWorkflows []Workflow `yaml:",omitempty"`
	AddBuckets   []Bucket   `yaml:",omitempty"`
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

func WriteDredgeFile(dredgeFile *DredgeFile, filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	defer encoder.Close()

	return encoder.Encode(dredgeFile)
}

func (r Runtime) GetHome() string {
	if r.Home == nil {
		return DEFAULT_HOME
	}
	return *r.Home
}
