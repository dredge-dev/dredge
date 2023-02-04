package resource

type NoResult struct{}

func (n *NoResult) Error() string {
	return "no result"
}

type Provider interface {
	Name() string
	Init(config map[string]string) error
	ExecuteCommand(commandName string, callbacks Callbacks) (interface{}, error)
}
