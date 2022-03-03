package exec

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/dredge-dev/dredge/internal/config"
)

type DredgeExec struct {
	Parent     *DredgeExec
	Source     config.SourcePath
	DredgeFile *config.DredgeFile
	Env        Env
}

type Bucket struct {
	Exec        *DredgeExec
	Name        string
	Description string
	workflows   []config.Workflow
}

type Workflow struct {
	Exec        *DredgeExec
	Name        string
	Description string
	Inputs      map[string]string
	Steps       []config.Step
}

func EmptyExec(source config.SourcePath) *DredgeExec {
	return &DredgeExec{
		Source:     source,
		DredgeFile: &config.DredgeFile{},
		Env:        NewEnv(),
	}
}

func NewExec(source config.SourcePath) (*DredgeExec, error) {
	content, err := ReadSource(source)
	if err != nil {
		return nil, err
	}

	dredgeFile, err := config.NewDredgeFile(content)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing %s: %s", source, err)
	}

	env := NewEnv()
	env.AddVariables(dredgeFile.Variables)

	return &DredgeExec{
		Source:     source,
		DredgeFile: dredgeFile,
		Env:        env,
	}, nil
}

func (exec *DredgeExec) Import(source config.SourcePath) (*DredgeExec, error) {
	fullSource := MergeSources(exec.Source, source)

	content, err := ReadSource(fullSource)
	if err != nil {
		return nil, err
	}

	imported, err := config.NewDredgeFile(content)
	if err != nil {
		return nil, err
	}

	env := exec.Env.Clone()
	env.AddVariables(imported.Variables)

	return &DredgeExec{
		Parent:     exec,
		Source:     fullSource,
		DredgeFile: imported,
		Env:        env,
	}, nil
}

func MergeSources(parent config.SourcePath, child config.SourcePath) config.SourcePath {
	c := string(child)
	p := string(parent)
	if c == "" {
		return parent
	} else if strings.HasPrefix(c, "./") {
		if strings.HasPrefix(p, "./") {
			parentPath := strings.Split(p, "/")
			parentDir := parentPath[:len(parentPath)-1]
			parts := append(parentDir, c[2:])
			return config.SourcePath(strings.Join(parts, "/"))
		}
	}
	return child
}

func ReadSource(source config.SourcePath) ([]byte, error) {
	s := string(source)
	if !strings.HasPrefix(s, "./") {
		return nil, fmt.Errorf("Sources should start with ./")
	}
	return ioutil.ReadFile(s)
}

func (exec *DredgeExec) ReadSource(source config.SourcePath) ([]byte, error) {
	return ReadSource(MergeSources(exec.Source, source))
}

func (exec *DredgeExec) GetWorkflows() ([]*Workflow, error) {
	var workflows []*Workflow
	for _, w := range exec.DredgeFile.Workflows {
		workflow, err := exec.resolveWorkflow(w)
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, workflow)
	}
	return workflows, nil
}

func (exec *DredgeExec) GetWorkflow(bucketName, workflowName string) (*Workflow, error) {
	if bucketName == "" {
		for _, w := range exec.DredgeFile.Workflows {
			if w.Name == workflowName {
				return exec.resolveWorkflow(w)
			}
		}
	} else {
		for _, b := range exec.DredgeFile.Buckets {
			if b.Name == bucketName {
				bucket, err := exec.resolveBucket(b)
				if err != nil {
					return nil, err
				}
				for _, w := range bucket.workflows {
					if w.Name == workflowName {
						return exec.resolveWorkflow(w)
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("Could not find workflow %s/%s", bucketName, workflowName)
}

func (exec *DredgeExec) GetBuckets() ([]*Bucket, error) {
	var buckets []*Bucket
	for _, b := range exec.DredgeFile.Buckets {
		bucket, err := exec.resolveBucket(b)
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, bucket)
	}
	return buckets, nil
}

func (exec *DredgeExec) GetBucket(bucketName string) (*Bucket, error) {
	for _, b := range exec.DredgeFile.Buckets {
		if b.Name == bucketName {
			return exec.resolveBucket(b)
		}
	}
	return nil, fmt.Errorf("Could not find bucket %s", bucketName)
}

func (exec *DredgeExec) resolveBucket(b config.Bucket) (*Bucket, error) {
	if b.Import != nil {
		de := exec
		if b.Import.Source != "" {
			var err error
			de, err = exec.Import(b.Import.Source)
			if err != nil {
				return nil, fmt.Errorf("Could not load Dredgefile %s: %v", b.Import.Source, err)
			}
		}
		bucket, err := de.GetBucket(b.Import.Bucket)
		if err != nil {
			return nil, err
		}
		bucket.Name = b.Name
		if b.Description != "" {
			bucket.Description = b.Description
		}
		return bucket, nil
	}
	return &Bucket{
		Exec:        exec,
		Name:        b.Name,
		Description: b.Description,
		workflows:   b.Workflows,
	}, nil
}

func (b *Bucket) GetWorkflows() ([]*Workflow, error) {
	var workflows []*Workflow
	for _, w := range b.workflows {
		workflow, err := b.Exec.resolveWorkflow(w)
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, workflow)
	}
	return workflows, nil
}

func (exec *DredgeExec) resolveWorkflow(w config.Workflow) (*Workflow, error) {
	if w.Import != nil {
		de := exec
		if w.Import.Source != "" {
			var err error
			de, err = exec.Import(w.Import.Source)
			if err != nil {
				return nil, fmt.Errorf("Could not load Dredgefile %s: %v", w.Import.Source, err)
			}
		}
		workflow, err := de.GetWorkflow(w.Import.Bucket, w.Import.Workflow)
		if err != nil {
			return nil, err
		}
		workflow.Name = w.Name
		if w.Description != "" {
			workflow.Description = w.Description
		}
		return workflow, nil
	}
	return &Workflow{
		Exec:        exec,
		Name:        w.Name,
		Description: w.Description,
		Inputs:      w.Inputs,
		Steps:       w.Steps,
	}, nil
}
