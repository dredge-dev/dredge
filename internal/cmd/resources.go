package cmd

import (
	"fmt"
	"strings"

	"github.com/dredge-dev/dredge/internal/exec"
	"github.com/spf13/cobra"
)

var textParam string

type ArgsParser func(args []string) (string, map[string]string, error)

func GetResourcesHelp(e *exec.DredgeExec) string {
	resources, err := e.GetResources()
	if err != nil || len(resources) == 0 {
		return ""
	}
	return "\n\nResources: " + strings.Join(resources, ", ")
}

func initResourceCommands(e *exec.DredgeExec, rootCmd *cobra.Command) error {
	rootCmd.AddCommand(&cobra.Command{
		Use:     "get <resource>",
		Short:   "Get all resources of the provided type",
		GroupID: "resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			return printOrErr(runResourceCommand("get", args, ResourceArgsParser, e))
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:     "create <resource>",
		Short:   "Create a resource of the provided type",
		GroupID: "resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			return printOrErr(runResourceCommand("create", args, ResourceArgsParser, e))
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:     "describe <resource>/<name>",
		Short:   "Describe a resource with the provided type and name",
		GroupID: "resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			return printOrErr(runResourceCommand("describe", args, ResourceAndNameArgsParser, e))
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:     "update <resource>/<name>",
		Short:   "Update a resource with the provided type and name",
		GroupID: "resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			return printOrErr(runResourceCommand("update", args, ResourceAndNameArgsParser, e))
		},
	})
	searchCmd := &cobra.Command{
		Use:     "search <resource>",
		Short:   "Search for a resource of the provided type",
		GroupID: "resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			return printOrErr(runResourceCommand("search", args, ResourceAndTextArgsParser, e))
		},
	}
	searchCmd.Flags().StringVar(&textParam, "text", "", "text to search")
	rootCmd.AddCommand(searchCmd)
	return nil
}

func ResourceArgsParser(args []string) (string, map[string]string, error) {
	if len(args) < 1 {
		return "", nil, fmt.Errorf("not enough arguments: missing <resource>")
	}
	return args[0], map[string]string{}, nil
}

func ResourceAndNameArgsParser(args []string) (string, map[string]string, error) {
	if len(args) < 1 {
		return "", nil, fmt.Errorf("not enough arguments: missing <resource>/<name>")
	}
	parts := strings.SplitN(args[0], "/", 2)
	if len(parts) != 2 {
		return "", nil, fmt.Errorf("argument not in format <resource>/<name>")
	}
	return parts[0], map[string]string{
		"name": parts[1],
	}, nil
}

func ResourceAndTextArgsParser(args []string) (string, map[string]string, error) {
	if len(args) < 1 {
		return "", nil, fmt.Errorf("not enough arguments: missing <resource>")
	}
	inputs := make(map[string]string)
	if textParam != "" {
		inputs["text"] = textParam
	}
	return args[0], inputs, nil
}

func runResourceCommand(command string, args []string, argsParser ArgsParser, e *exec.DredgeExec) (string, error) {
	resourceName, namedArgs, err := argsParser(args)
	if err != nil {
		return "", err
	}

	e.Env.AddInputs(namedArgs)
	output, err := e.ExecuteResourceCommand(resourceName, command)
	if err != nil {
		return "", err
	}

	return format(output)
}

func printOrErr(output string, err error) error {
	if err != nil {
		return err
	}
	fmt.Print(output)
	return nil
}
