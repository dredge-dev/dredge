package exec

import (
	"fmt"

	"github.com/dredge-dev/dredge/internal/api"
	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/workflow"
)

type DredgeExec struct {
	Parent              *DredgeExec
	Source              config.SourcePath
	DredgeFile          *config.DredgeFile
	Env                 Env
	ResourceDefinitions []api.ResourceDefinition
	callbacks           api.UserInteractionCallbacks
}

func EmptyExec(source config.SourcePath, rd []api.ResourceDefinition, c api.UserInteractionCallbacks) *DredgeExec {
	return &DredgeExec{
		Source:              source,
		DredgeFile:          &config.DredgeFile{},
		Env:                 NewEnv(),
		ResourceDefinitions: rd,
		callbacks:           c,
	}
}

func NewExec(source config.SourcePath, rd []api.ResourceDefinition, c api.UserInteractionCallbacks) (*DredgeExec, error) {
	actualSource, dredgeFile, err := ReadDredgeFile(source)
	if err != nil {
		return nil, err
	}

	env := NewEnv()
	env.AddVariables(dredgeFile.Variables)

	return &DredgeExec{
		Source:              actualSource,
		DredgeFile:          dredgeFile,
		Env:                 env,
		ResourceDefinitions: rd,
		callbacks:           c,
	}, nil
}

func (exec *DredgeExec) Import(source config.SourcePath) (*DredgeExec, error) {
	fullSource := MergeSources(exec.Source, source)

	actualSource, imported, err := ReadDredgeFile(fullSource)
	if err != nil {
		return nil, err
	}

	env := exec.Env.Clone()
	env.AddVariables(imported.Variables)

	return &DredgeExec{
		Parent:              exec,
		Source:              actualSource,
		DredgeFile:          imported,
		Env:                 env,
		ResourceDefinitions: exec.ResourceDefinitions,
		callbacks:           exec.callbacks,
	}, nil
}

func (exec *DredgeExec) GetWorkflows() ([]*workflow.Workflow, error) {
	var workflows []*workflow.Workflow
	for _, w := range exec.DredgeFile.Workflows {
		workflow, err := exec.resolveWorkflow(w)
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, workflow)
	}
	return workflows, nil
}

func (exec *DredgeExec) GetWorkflow(bucketName, workflowName string) (*workflow.Workflow, error) {
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
				for _, w := range bucket.Workflows {
					if w.Name == workflowName {
						return exec.resolveWorkflow(w)
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("could not find workflow %s/%s", bucketName, workflowName)
}

func (exec *DredgeExec) GetBuckets() ([]*workflow.Bucket, error) {
	var buckets []*workflow.Bucket
	for _, b := range exec.DredgeFile.Buckets {
		bucket, err := exec.resolveBucket(b)
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, bucket)
	}
	return buckets, nil
}

func (exec *DredgeExec) GetBucket(bucketName string) (*workflow.Bucket, error) {
	for _, b := range exec.DredgeFile.Buckets {
		if b.Name == bucketName {
			return exec.resolveBucket(b)
		}
	}
	return nil, fmt.Errorf("could not find bucket %s", bucketName)
}

func (exec *DredgeExec) resolveBucket(b config.Bucket) (*workflow.Bucket, error) {
	if b.Import != nil {
		de := exec
		if b.Import.Source != "" {
			var err error
			de, err = exec.Import(b.Import.Source)
			if err != nil {
				return nil, fmt.Errorf("could not load Dredgefile %s: %v", b.Import.Source, err)
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
	return &workflow.Bucket{
		Name:        b.Name,
		Description: b.Description,
		Workflows:   b.Workflows,
		Callbacks:   exec,
	}, nil
}

func (exec *DredgeExec) GetWorkflowsInBucket(b *workflow.Bucket) ([]*workflow.Workflow, error) {
	var workflows []*workflow.Workflow
	for _, w := range b.Workflows {
		workflow, err := exec.resolveWorkflow(w)
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, workflow)
	}
	return workflows, nil
}

func (exec *DredgeExec) resolveWorkflow(w config.Workflow) (*workflow.Workflow, error) {
	if w.Import != nil {
		de := exec
		if w.Import.Source != "" {
			var err error
			de, err = exec.Import(w.Import.Source)
			if err != nil {
				return nil, fmt.Errorf("could not load Dredgefile %s: %v", w.Import.Source, err)
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
	return &workflow.Workflow{
		Name:        w.Name,
		Description: w.Description,
		Inputs:      w.Inputs,
		Steps:       w.Steps,
		Runtimes:    exec.DredgeFile.Runtimes, // TODO I think this breaks with imports
		Callbacks:   exec,
	}, nil
}

func (e *DredgeExec) getRootExec() *DredgeExec {
	exec := e
	for exec.Parent != nil {
		exec = exec.Parent
	}
	return exec
}

func (e *DredgeExec) getRootExecAndDredgeFile() (*DredgeExec, *config.DredgeFile) {
	rootExec := e.getRootExec()
	return rootExec, rootExec.DredgeFile
}
