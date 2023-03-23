package api

import "fmt"

func (r *ResourceDefinition) GetCommand(name string) (*Command, error) {
	for _, c := range r.Commands {
		if c.Name == name {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("could not find %s command for %s resource", name, r.Name)
}
