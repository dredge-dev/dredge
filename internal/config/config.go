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
	Source    string
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
	dredgeFile.Source = source

	return dredgeFile, nil
}

func (dredgeFile *DredgeFile) ImportDredgeFile(source string) (*DredgeFile, error) {
	imported, err := GetDredgeFile(source)
	if err != nil {
		return imported, err
	}
	imported.Source = mergeSources(dredgeFile.Source, imported.Source)
	for _, w := range imported.Workflows {
		if w.Import != nil {
			w.Import.Source = mergeSources(imported.Source, w.Import.Source)
		}
	}
	for _, b := range imported.Buckets {
		if b.Import != nil {
			b.Import.Source = mergeSources(imported.Source, b.Import.Source)
		}
		for _, w := range b.Workflows {
			if w.Import != nil {
				w.Import.Source = mergeSources(imported.Source, b.Import.Source)
			}
		}
	}
	return imported, nil
}

func mergeSources(parent, child string) string {
	if strings.HasPrefix(child, "./") {
		if strings.HasPrefix(parent, "./") {
			parentPath := strings.Split(parent, "/")
			parentDir := parentPath[:len(parentPath)-1]
			parts := append(parentDir, child[2:])
			return strings.Join(parts, "/")
		}
	}
	return child
}

func (dredgeFile *DredgeFile) GetWorkflow(bucketName, workflowName string) (*Workflow, error) {
	if bucketName == "" {
		for _, w := range dredgeFile.Workflows {
			if w.Name == workflowName {
				return &w, nil
			}
		}
	} else {
		for _, b := range dredgeFile.Buckets {
			if b.Name == bucketName {
				workflows, err := dredgeFile.GetWorkflows(b)
				if err != nil {
					return nil, err
				}
				for _, w := range workflows {
					if w.Name == workflowName {
						return &w, nil
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("Could not find workflow %s/%s", bucketName, workflowName)
}

func (dredgeFile *DredgeFile) GetBucket(bucketName string) (*Bucket, error) {
	for _, b := range dredgeFile.Buckets {
		if b.Name == bucketName {
			return &b, nil
		}
	}
	return nil, fmt.Errorf("Could not find bucket %s", bucketName)
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
	if b.Import != nil {
		if len(b.Workflows) > 0 {
			return fmt.Errorf("bucket %s: contains both workflows and an import", b.Name)
		}
		if err := b.Import.Validate(); err != nil {
			return fmt.Errorf("bucket %s: %v", b.Name, err)
		}
		return nil
	}
	for _, w := range b.Workflows {
		if err := w.Validate(); err != nil {
			return fmt.Errorf("bucket %s: %v", b.Name, err)
		}
	}
	return nil
}

func (i ImportBucket) Validate() error {
	if i.Bucket == "" {
		return fmt.Errorf("bucket field is required for import")
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

func (dredgeFile *DredgeFile) ResolveBucket(b *Bucket) (*DredgeFile, *Bucket, error) {
	if b.Import != nil {
		source := dredgeFile
		if b.Import.Source != "" {
			var err error
			source, err = dredgeFile.ImportDredgeFile(b.Import.Source)
			if err != nil {
				return nil, nil, fmt.Errorf("Could not load Dredgefile %s", b.Import.Source)
			}
		}

		bucket, err := source.GetBucket(b.Import.Bucket)
		if err != nil {
			return nil, nil, err
		}

		return source.ResolveBucket(bucket)
	}
	return dredgeFile, b, nil
}

func (dredgeFile *DredgeFile) GetWorkflows(bucket Bucket) ([]Workflow, error) {
	_, b, err := dredgeFile.ResolveBucket(&bucket)
	if err != nil {
		return nil, err
	}
	return b.Workflows, nil
}

func (dredgeFile *DredgeFile) GetBucketDescription(bucket Bucket) (string, error) {
	if bucket.Description == "" {
		_, b, err := dredgeFile.ResolveBucket(&bucket)
		if err != nil {
			return "", err
		}
		return b.Description, nil
	}
	return bucket.Description, nil
}

func (dredgeFile *DredgeFile) ResolveWorkflow(w *Workflow) (*DredgeFile, *Workflow, error) {
	if w.Import != nil {
		source := dredgeFile
		if w.Import.Source != "" {
			var err error
			source, err = dredgeFile.ImportDredgeFile(w.Import.Source)
			if err != nil {
				return nil, nil, fmt.Errorf("Could not load Dredgefile %s", w.Import.Source)
			}
		}

		workflow, err := source.GetWorkflow(w.Import.Bucket, w.Import.Workflow)
		if err != nil {
			return nil, nil, err
		}

		return source.ResolveWorkflow(workflow)
	}
	return dredgeFile, w, nil
}

func (dredgeFile *DredgeFile) GetWorkflowDescription(workflow Workflow) (string, error) {
	if workflow.Description == "" {
		_, w, err := dredgeFile.ResolveWorkflow(&workflow)
		if err != nil {
			return "", err
		}
		return w.Description, nil
	}
	return workflow.Description, nil
}

func (r Runtime) GetHome() string {
	if r.Home == nil {
		return DEFAULT_HOME
	}
	return *r.Home
}
