package cmd

import (
	"fmt"
	"strings"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/workflow"
	"github.com/spf13/cobra"
)

func runExecCommand(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Not enough arguments: missing target")
	}

	dredgeFile, err := getDredgeFile(args[0])
	if err != nil {
		return err
	}

	if len(args) == 1 {
		newCmd := &cobra.Command{
			Use: "exec",
		}
		addWorkflows(dredgeFile, newCmd)
		return newCmd.Help()
	} else {
		var w *config.Workflow
		if len(args) == 2 {
			w = dredgeFile.GetWorkflow(args[1])
		} else {
			w = dredgeFile.GetWorkflowInBucket(args[1], args[2])
		}
		if w != nil {
			return workflow.ExecuteWorkflow(dredgeFile, *w)
		}
		if len(args) == 2 {
			b := dredgeFile.GetBucket(args[1])
			if b != nil {
				cmd := createBucketCommand(dredgeFile, *b)
				return cmd.Help()
			}
		}
		return fmt.Errorf("Could not find workflow %s in %s", strings.Join(args[1:], "/"), args[0])
	}
}

func getDredgeFile(target string) (*config.DredgeFile, error) {
	filename := target

	if !strings.HasPrefix(filename, "./") {
		return nil, fmt.Errorf("Targets should start with ./")
	}

	dredgeFile, err := config.ReadDredgeFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing %s: %s", filename, err)
	}

	return dredgeFile, nil
}
