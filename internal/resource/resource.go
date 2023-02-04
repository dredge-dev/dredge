package resource

import (
	"fmt"
	"reflect"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/exec"
)

type Resource struct {
	Exec       *exec.DredgeExec
	Name       string
	Definition *ResourceDefinition
	Providers  []Provider
}

type ProviderCreator func(*exec.DredgeExec, config.ResourceProvider) (Provider, error)

type CommandOutput struct {
	Type   *Type
	Output interface{}
}

func GetResources(e *exec.DredgeExec) ([]string, error) {
	var resources []string
	for name, _ := range e.DredgeFile.Resources {
		_, err := GetResourceDefinition(e, name)
		if err != nil {
			return nil, err
		}
		resources = append(resources, name)
	}
	return resources, nil
}

func GetResource(e *exec.DredgeExec, create ProviderCreator, resourceName string) (*Resource, error) {
	rd, err := GetResourceDefinition(e, resourceName)
	if err != nil {
		return nil, err
	}

	r, ok := e.DredgeFile.Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("could not find resource %s", resourceName)
	}
	if len(r) == 0 {
		return nil, fmt.Errorf("no provider specified for this resource")
	}

	var providers []Provider
	for _, p := range r {
		provider, err := create(e, p)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}

	return &Resource{
		Exec:       e,
		Name:       resourceName,
		Definition: rd,
		Providers:  providers,
	}, nil
}

func (r *Resource) ExecuteCommand(command string, callbacks Callbacks) (*CommandOutput, error) {
	c, err := r.Definition.GetCommand(command)
	if err != nil {
		return nil, err
	}

	outputType, err := GetType(r.Exec, c.OutputType)
	if err != nil {
		return nil, err
	}

	// TODO If the result is not an array, stop when the first provider returns non-nil value
	var outputs []interface{}
	for _, provider := range r.Providers {
		output, err := provider.ExecuteCommand(command, &DredgeEnvCallbacks{r.Name, r.Exec, callbacks})
		if err != nil {
			if _, ok := err.(*NoResult); !ok {
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
		return &CommandOutput{outputType, flat}, nil
	}

	if len(outputs) == 0 {
		return nil, fmt.Errorf("no result returned by provider(s)")
	}
	if len(outputs) > 1 {
		return nil, fmt.Errorf("1 result expected, more than 1 provider returned")
	}
	return &CommandOutput{outputType, outputs[0]}, nil
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

type DredgeEnvCallbacks struct {
	ResourceName string
	Exec         *exec.DredgeExec
	Callbacks    Callbacks
}

func (c *DredgeEnvCallbacks) Log(level LogLevel, msg string) error {
	return c.Callbacks.Log(level, msg)
}

func (c *DredgeEnvCallbacks) RequestInput(inputRequests []InputRequest) (map[string]string, error) {
	// TODO add inputs to the environment so it doesn't get asked twice
	inputs := make(map[string]string)
	var remainingRequests []InputRequest

	for _, inputRequest := range inputRequests {
		fullName := c.ResourceName + "." + inputRequest.Name
		if value, ok := c.Exec.Env[fullName]; ok {
			inputs[inputRequest.Name] = value
		} else {
			remainingRequests = append(remainingRequests, inputRequest)
		}
	}

	if len(remainingRequests) > 0 {
		remainingInputs, err := c.Callbacks.RequestInput(remainingRequests)
		if err != nil {
			return nil, err
		}

		for inputName, inputValue := range remainingInputs {
			inputs[inputName] = inputValue
		}
	}

	return inputs, nil
}
