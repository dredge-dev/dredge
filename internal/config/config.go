package config

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

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

func readDredgeFile(filename string) (*DredgeFile, error) {
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
	return dredgeFile, err
}

func GetDredgeFile(source string) (*DredgeFile, error) {
	filename := source

	if !strings.HasPrefix(filename, "./") {
		return nil, fmt.Errorf("Sources should start with ./")
	}

	dredgeFile, err := readDredgeFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing %s: %s", filename, err)
	}

	return dredgeFile, nil
}

func (dredgeFile *DredgeFile) GetWorkflow(bucketName, workflowName string) *Workflow {
	if bucketName == "" {
		for _, w := range dredgeFile.Workflows {
			if w.Name == workflowName {
				return &w
			}
		}
	} else {
		for _, b := range dredgeFile.Buckets {
			if b.Name == bucketName {
				for _, w := range b.Workflows {
					if w.Name == workflowName {
						return &w
					}
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

func (dredgeFile *DredgeFile) Validate() error {
	for _, w := range dredgeFile.Workflows {
		if err := w.Validate(); err != nil {
			return err
		}
	}
	for _, b := range dredgeFile.Buckets {
		if err := b.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (b Bucket) Validate() error {
	if b.Name == "" {
		return fmt.Errorf("name field is required for bucket")
	}
	for _, w := range b.Workflows {
		if err := w.Validate(); err != nil {
			return fmt.Errorf("bucket %s: %v", b.Name, err)
		}
	}
	return nil
}

func (w Workflow) Validate() error {
	if w.Name == "" {
		return fmt.Errorf("name field is required for workflow")
	}
	if w.Import != nil {
		if len(w.Steps) > 0 {
			return fmt.Errorf("workflow %s: contains both steps and an import", w.Name)
		}
		if err := w.Import.Validate(); err != nil {
			return fmt.Errorf("workflow %s: %v", w.Name, err)
		}
		return nil
	}
	if len(w.Steps) == 0 {
		return fmt.Errorf("workflow %s: no steps or import defined", w.Name)
	}
	for _, s := range w.Steps {
		if err := s.Validate(); err != nil {
			return fmt.Errorf("workflow %s: %v", w.Name, err)
		}
	}
	return nil
}

func (i ImportWorkflow) Validate() error {
	if i.Workflow == "" {
		return fmt.Errorf("workflow field is required for import")
	}
	return nil
}

func (s Step) Validate() error {
	// TODO Add validate for each step type
	numFields := 0

	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		if v.Type().Field(i).Name != "Name" {
			if !v.Field(i).IsNil() {
				numFields += 1
			}
		}
	}

	if numFields == 0 {
		return fmt.Errorf("step %s does not contain an action", s.Name)
	} else if numFields == 1 {
		return nil
	} else {
		return fmt.Errorf("step %s contains more than 1 action", s.Name)
	}
}

func (w Workflow) GetDescription() string {
	description := w.Description
	if description == "" && w.Import != nil && w.Import.Source != "" {
		if idf, _ := GetDredgeFile(w.Import.Source); idf != nil {
			if iw := idf.GetWorkflow(w.Import.Bucket, w.Import.Workflow); iw != nil {
				description = iw.GetDescription()
			}
		}
	}
	return description
}

func (r Runtime) GetHome() string {
	if r.Home == nil {
		return DEFAULT_HOME
	}
	return *r.Home
}
