package workflow

import (
	"fmt"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/exec"
)

func executeEditDredgeFile(workflow *exec.Workflow, edit *config.EditDredgeFileStep) error {
	rootExec := getRootExec(workflow.Exec)
	df := rootExec.DredgeFile
	changed := false

	if len(edit.AddVariables) > 0 {
		if df.Variables == nil {
			df.Variables = make(config.Variables)
		}
		for variable, value := range edit.AddVariables {
			if _, ok := df.Variables[variable]; !ok {
				templatedValue, err := Template(value, workflow.Exec.Env)
				if err != nil {
					return err
				}
				df.Variables[variable] = templatedValue
				changed = true
			} else {
				fmt.Printf("Skipping adding variable %s to %s, already present.\n", variable, rootExec.Source)
			}
		}
	}

	for _, w := range edit.AddWorkflows {
		if w.Import != nil {
			w.Import.Source = exec.MergeSources(workflow.Exec.Source, w.Import.Source)
		}
		if f, _ := rootExec.GetWorkflow("", w.Name); f != nil {
			fmt.Printf("Skipping adding workflow %s to %s, already present.\n", w.Name, rootExec.Source)
		} else {
			df.Workflows = append(df.Workflows, w)
			changed = true
		}
	}

	for _, b := range edit.AddBuckets {
		if b.Import != nil {
			b.Import.Source = exec.MergeSources(workflow.Exec.Source, b.Import.Source)
		}
		if f, _ := rootExec.GetBucket(b.Name); f != nil {
			fmt.Printf("Skipping adding bucket %s to %s, already present.\n", b.Name, rootExec.Source)
		} else {
			df.Buckets = append(df.Buckets, b)
			changed = true
		}
	}

	if changed {
		if err := df.Validate(); err != nil {
			return err
		}
		return config.WriteDredgeFile(df, rootExec.Source)
	}
	return nil
}

func getRootExec(e *exec.DredgeExec) *exec.DredgeExec {
	exec := e
	for exec.Parent != nil {
		exec = exec.Parent
	}
	return exec
}
