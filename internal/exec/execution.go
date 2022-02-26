package exec

import (
	"fmt"
	"strings"

	"github.com/dredge-dev/dredge/internal/config"
)

type DredgeExec struct {
	Source     string
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

func EmptyExec() *DredgeExec {
	return &DredgeExec{
		Source:     "",
		DredgeFile: &config.DredgeFile{},
		Env:        NewEnv(),
	}
}

func NewExec(source string) (*DredgeExec, error) {
	if !strings.HasPrefix(source, "./") {
		return nil, fmt.Errorf("Sources should start with ./")
	}

	dredgeFile, err := config.ReadDredgeFile(source)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing %s: %s", source, err)
	}

	env := NewEnv()
	env.AddVariables(dredgeFile.Env)

	return &DredgeExec{
		Source:     source,
		DredgeFile: dredgeFile,
		Env:        env,
	}, nil
}

func (exec *DredgeExec) _import(source string) (*DredgeExec, error) {
	fullSource := mergeSources(exec.Source, source)

	imported, err := config.ReadDredgeFile(fullSource)
	if err != nil {
		return nil, err
	}

	env := exec.Env.Clone()
	env.AddVariables(imported.Env)

	return &DredgeExec{
		Source:     fullSource,
		DredgeFile: imported,
		Env:        env,
	}, nil
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
			de, err = exec._import(b.Import.Source)
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
			de, err = exec._import(w.Import.Source)
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
