package workflow

import (
	"github.com/dredge-dev/dredge/internal/config"
)

func (workflow *Workflow) executeEditDredgeFile(edit *config.EditDredgeFileStep) error {
	if len(edit.AddVariables) > 0 {
		toAdd := make(map[string]string)
		for variable, value := range edit.AddVariables {
			templatedValue, err := workflow.Callbacks.Template(value)
			if err != nil {
				return err
			}
			toAdd[variable] = templatedValue
		}
		err := workflow.Callbacks.AddVariablesToDredgefile(toAdd)
		if err != nil {
			return err
		}
	}

	for _, w := range edit.AddWorkflows {
		err := workflow.Callbacks.AddWorkflowToDredgefile(w)
		if err != nil {
			return err
		}
	}

	for _, b := range edit.AddBuckets {
		err := workflow.Callbacks.AddBucketToDredgefile(b)
		if err != nil {
			return err
		}
	}

	return nil
}
