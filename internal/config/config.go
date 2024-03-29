package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	DEFAULT_HOME      = "/home"
	INPUT_TEXT        = "text"
	INPUT_SELECT      = "select"
	INSERT_BEGIN      = "begin"
	INSERT_END        = "end"
	INSERT_UNIQUE     = "unique"
	RUNTIME_NATIVE    = "native"
	RUNTIME_CONTAINER = "container"
	LOG_FATAL         = "fatal"
	LOG_ERROR         = "error"
	LOG_WARN          = "warn"
	LOG_INFO          = "info"
	LOG_DEBUG         = "debug"
	LOG_TRACE         = "trace"
)

type DredgeFile struct {
	Variables Variables  `yaml:",omitempty"`
	Runtimes  []Runtime  `yaml:",omitempty"`
	Workflows []Workflow `yaml:",omitempty"`
	Buckets   []Bucket   `yaml:",omitempty"`
	Resources Resources  `yaml:",omitempty"`
}

type Variables map[string]string
type SourcePath string

type Runtime struct {
	Name        string
	Type        string
	Image       string            `yaml:",omitempty"`
	Home        string            `yaml:",omitempty"`
	Cache       []string          `yaml:",omitempty"`
	GlobalCache []string          `yaml:"global_cache,omitempty"`
	Ports       []string          `yaml:",omitempty"`
	EnvVars     map[string]string `yaml:",omitempty"`
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
	DefaultValue string   `yaml:"default_value,omitempty"`
	Skip         string   `yaml:",omitempty"`
}

type Step struct {
	Name           string              `yaml:",omitempty"`
	Shell          *ShellStep          `yaml:",omitempty"`
	Template       *TemplateStep       `yaml:",omitempty"`
	Browser        *BrowserStep        `yaml:",omitempty"`
	EditDredgeFile *EditDredgeFileStep `yaml:"edit_dredgefile,omitempty"`
	If             *IfStep             `yaml:",omitempty"`
	Execute        *ExecuteStep        `yaml:",omitempty"`
	Set            *SetStep            `yaml:",omitempty"`
	Log            *LogStep            `yaml:",omitempty"`
	Confirm        *ConfirmStep        `yaml:",omitempty"`
}

type ShellStep struct {
	Cmd     string
	Runtime string `yaml:",omitempty"`
	StdOut  string `yaml:"stdout,omitempty"`
	StdErr  string `yaml:"stderr,omitempty"`
}

type TemplateStep struct {
	Source SourcePath `yaml:",omitempty"`
	Input  string     `yaml:",omitempty"`
	Dest   string
	Insert *Insert `yaml:",omitempty"`
}

type Insert struct {
	Section   string `yaml:",omitempty"`
	Placement string `yaml:",omitempty"`
}

type BrowserStep struct {
	Url string
}

type EditDredgeFileStep struct {
	AddVariables Variables  `yaml:"add_variables,omitempty"`
	AddWorkflows []Workflow `yaml:"add_workflows,omitempty"`
	AddBuckets   []Bucket   `yaml:"add_buckets,omitempty"`
}

type IfStep struct {
	Cond  string
	Steps []Step `yaml:",omitempty"`
}

type ExecuteStep struct {
	Resource string
	Command  string
	Register string `yaml:",omitempty"`
}

type SetStep map[string]string

type LogStep struct {
	Level   string
	Message string
}

type ConfirmStep struct {
	Message string
}

type Resources map[string]Resource

type Resource []ResourceProvider

type ResourceProvider struct {
	Provider string
	Config   map[string]string `yaml:",omitempty"`
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
		return fmt.Errorf("cannot write to non-local file %s", f)
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
	if r.Home == "" {
		return DEFAULT_HOME
	}
	return r.Home
}

func (i Input) HasValue(value string) bool {
	for _, v := range i.Values {
		if value == v {
			return true
		}
	}
	return false
}
