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
		return fmt.Errorf("Not enough arguments: missing source")
	}

	dredgeFile, err := config.GetDredgeFile(args[0])
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
			w = dredgeFile.GetWorkflow("", args[1])
		} else {
			w = dredgeFile.GetWorkflow(args[1], args[2])
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
