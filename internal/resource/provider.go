package resource

import "github.com/dredge-dev/dredge/internal/callbacks"

type Provider interface {
	Name() string
	Init(config map[string]string) error
	ExecuteCommand(commandName string, callbacks callbacks.Callbacks) (interface{}, error)
}
