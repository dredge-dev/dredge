package resource

import (
	"fmt"
	"reflect"

	"github.com/dredge-dev/dredge/internal/api"
	"github.com/dredge-dev/dredge/internal/config"
)

type Resource struct {
	Definition *api.ResourceDefinition
	Providers  []Provider
}

type Provider interface {
	Name() string
	Discover(callbacks api.Callbacks) error
	Init(config map[string]string) error
	ExecuteCommand(commandName string, callbacks api.Callbacks) (interface{}, error)
}

func NewResource(rd *api.ResourceDefinition, r config.Resource) (*Resource, error) {
	var providers []Provider
	for _, p := range r {
		provider, err := CreateProvider(p)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}

	return &Resource{
		Definition: rd,
		Providers:  providers,
	}, nil
}

func (r *Resource) ExecuteCommand(command string, outputType *api.Type, c api.Callbacks) (*api.CommandOutput, error) {
	// TODO If the result is not an array, stop when the first provider returns non-nil value
	var outputs []interface{}
	for _, provider := range r.Providers {
		output, err := provider.ExecuteCommand(command, c)
		if err != nil {
			if _, ok := err.(*api.NoResult); !ok {
				return nil, err
			}
		} else {
			outputs = append(outputs, output)
		}
	}

	if outputType.IsArray {
		flat, err := flatten(outputs)
		if err != nil {
			return nil, err
		}
		return &api.CommandOutput{Type: outputType, Output: flat}, nil
	}

	if len(outputs) == 0 {
		return nil, fmt.Errorf("no result returned by provider(s)")
	}
	if len(outputs) > 1 {
		return nil, fmt.Errorf("1 result expected, more than 1 provider returned")
	}
	return &api.CommandOutput{Type: outputType, Output: outputs[0]}, nil
}

func flatten(outputs []interface{}) ([]interface{}, error) {
	var flat []interface{}
	for _, output := range outputs {
		s := reflect.ValueOf(output)
		if s.Kind() != reflect.Slice {
			return nil, fmt.Errorf("expected array type but provider returned object")
		}
		for i := 0; i < s.Len(); i++ {
			flat = append(flat, s.Index(i).Interface())
		}
	}
	return flat, nil
}
