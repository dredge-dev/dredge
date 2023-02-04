package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dredge-dev/dredge/internal/exec"
	"github.com/dredge-dev/dredge/internal/resource"
	"github.com/dredge-dev/dredge/internal/resource/providers"
	"github.com/spf13/cobra"
)

var textParam string

type ArgsParser func(args []string) (string, map[string]string, error)

func initResourceCommands(e *exec.DredgeExec, rootCmd *cobra.Command) error {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "get <resource>",
		Short: "Get all resources of the provided type",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runResourceCommand("get", args, ResourceArgsParser, e, os.Stdin, os.Stdout)
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "create <resource>",
		Short: "Create a resource of the provided type",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runResourceCommand("create", args, ResourceArgsParser, e, os.Stdin, os.Stdout)
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "describe <resource>/<name>",
		Short: "Describe a resource with the provided type and name",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runResourceCommand("describe", args, ResourceAndNameArgsParser, e, os.Stdin, os.Stdout)
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "update <resource>/<name>",
		Short: "Update a resource with the provided type and name",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runResourceCommand("update", args, ResourceAndNameArgsParser, e, os.Stdin, os.Stdout)
		},
	})
	searchCmd := &cobra.Command{
		Use:   "search <resource>",
		Short: "Search for a resource of the provided type",

		RunE: func(cmd *cobra.Command, args []string) error {
			return runResourceCommand("search", args, ResourceAndTextArgsParser, e, os.Stdin, os.Stdout)
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
		parts[0] + ".name": parts[1],
	}, nil
}

func ResourceAndTextArgsParser(args []string) (string, map[string]string, error) {
	if len(args) < 1 {
		return "", nil, fmt.Errorf("not enough arguments: missing <resource>")
	}
	resourceName := args[0]
	inputs := make(map[string]string)
	if textParam != "" {
		inputs[resourceName+".text"] = textParam
	}
	return resourceName, inputs, nil
}

func runResourceCommand(command string, args []string, argsParser ArgsParser, e *exec.DredgeExec, reader io.Reader, writer io.Writer) error {
	resourceName, namedArgs, err := argsParser(args)
	if err != nil {
		return err
	}

	e.Env.AddInputs(namedArgs)

	r, err := resource.GetResource(e, providers.CreateProvider, resourceName)
	if err != nil {
		return err
	}

	output, err := r.ExecuteCommand(command, CliCallbacks{reader, writer})
	if err != nil {
		return err
	}

	return formatted(output, writer)
}

func formatted(output *resource.CommandOutput, writer io.Writer) error {
	formatted, err := format(output)
	if err != nil {
		return err
	}

	fmt.Fprint(writer, formatted)
	return nil
}
