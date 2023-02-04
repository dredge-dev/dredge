package resource

import (
	"fmt"

	"github.com/dredge-dev/dredge/internal/exec"
)

type ResourceDefinition struct {
	Name     string
	Fields   []Field
	Commands []Command
}

type Field struct {
	Name        string
	Description string
	Type        string
}

type Command struct {
	Name       string
	Inputs     []string
	OutputType string
}

func GetResourceDefinition(de *exec.DredgeExec, resourceName string) (*ResourceDefinition, error) {
	defaults := GetDefaultResourceDefinitions()
	for _, rd := range defaults {
		if rd.Name == resourceName {
			return &rd, nil
		}
	}
	return nil, fmt.Errorf("could not find resource definition for %s", resourceName)
}

func (r *ResourceDefinition) GetCommand(name string) (*Command, error) {
	for _, c := range r.Commands {
		if c.Name == name {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("could not find %s command for %s resource", name, r.Name)
}
