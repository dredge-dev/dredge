package exec

import (
	"github.com/dredge-dev/dredge/internal/api"
)

func (e *DredgeExec) Log(level api.LogLevel, msg string) error {
	return e.callbacks.Log(level, msg)
}

func (e *DredgeExec) RequestInput(inputRequests []api.InputRequest) (map[string]string, error) {
	// TODO add inputs to the environment so it doesn't get asked twice
	inputs := make(map[string]string)
	var remainingRequests []api.InputRequest

	for _, inputRequest := range inputRequests {
		if value, ok := e.Env[inputRequest.Name]; ok {
			inputs[inputRequest.Name] = value
		} else {
			remainingRequests = append(remainingRequests, inputRequest)
		}
	}

	if len(remainingRequests) > 0 {
		remainingInputs, err := e.callbacks.RequestInput(remainingRequests)
		if err != nil {
			return nil, err
		}

		for inputName, inputValue := range remainingInputs {
			inputs[inputName] = inputValue
		}
	}

	return inputs, nil
}

func (e *DredgeExec) OpenUrl(url string) error {
	return e.callbacks.OpenUrl(url)
}

func (e *DredgeExec) ExecuteResourceCommand(resourceName string, commandName string) (*api.CommandOutput, error) {
	r, err := e.GetResource(resourceName)
	if err != nil {
		return nil, err
	}

	commDef, err := r.Definition.GetCommand(commandName)
	if err != nil {
		return nil, err
	}

	outputType, err := e.GetType(commDef.OutputType)
	if err != nil {
		return nil, err
	}

	return r.ExecuteCommand(commandName, outputType, e)
}
