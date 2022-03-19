package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const DEFAULT_HOME = "/home"

type DredgeFile struct {
	Variables Variables  `yaml:",omitempty"`
	Runtimes  []Runtime  `yaml:",omitempty"`
	Workflows []Workflow `yaml:",omitempty"`
	Buckets   []Bucket   `yaml:",omitempty"`
}

type Variables map[string]string
type SourcePath string

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
	Source SourcePath
	Bucket string
}

type Workflow struct {
	Name        string
	Description string          `yaml:",omitempty"`
	Inputs      []Input         `yaml:",omitempty"`
	Steps       []Step          `yaml:",omitempty"`
	Import      *ImportWorkflow `yaml:",omitempty"`
}

type ImportWorkflow struct {
	Source   SourcePath
	Bucket   string
	Workflow string
}

type Input struct {
	Name         string
	Description  string   `yaml:",omitempty"`
	Type         string   `yaml:",omitempty"`
	Values       []string `yaml:",omitempty"`
	DefaultValue string   `yaml:default_value",omitempty"`
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
	Source SourcePath `yaml:",omitempty"`
	Input  string     `yaml:",omitempty"`
	Dest   string
}

type BrowserStep struct {
	Url string
}

type EditDredgeFileStep struct {
	AddVariables Variables  `yaml:"add_variables,omitempty"`
	AddWorkflows []Workflow `yaml:"add_workflows,omitempty"`
	AddBuckets   []Bucket   `yaml:"add_buckets,omitempty"`
}

func NewDredgeFile(buf []byte) (*DredgeFile, error) {
	dredgeFile := &DredgeFile{}
	err := yaml.Unmarshal(buf, dredgeFile)
	if err != nil {
		return nil, err
	}
	err = dredgeFile.Validate()
	if err != nil {
		return nil, err
	}
	return dredgeFile, nil
}

func WriteDredgeFile(dredgeFile *DredgeFile, filename SourcePath) error {
	f := string(filename)
	if !strings.HasPrefix(f, "./") {
		return fmt.Errorf("Cannot write to non-local file %s", f)
	}

	file, err := os.OpenFile(f, os.O_CREATE|os.O_WRONLY, 0644)
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

func (i Input) HasValue(value string) bool {
	for _, v := range i.Values {
		if value == v {
			return true
		}
	}
	return false
}
