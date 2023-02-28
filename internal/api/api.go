package api

import "github.com/dredge-dev/dredge/internal/config"

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

type CommandOutput struct {
	Type   *Type
	Output interface{}
}

type Type struct {
	Name    string
	IsArray bool
	Fields  []Field
}

type UserInteractionCallbacks interface {
	Log(level LogLevel, msg string) error
	RequestInput(inputRequests []InputRequest) (map[string]string, error)
	OpenUrl(url string) error
	Confirm(msg string) error
}

type ExecutionCallbacks interface {
	ExecuteResourceCommand(resource string, command string) (*CommandOutput, error)
	SetEnv(name string, value interface{}) error
	Template(input string) (string, error)
}

type DredgefileCallbacks interface {
	AddVariablesToDredgefile(variable map[string]string) error
	AddWorkflowToDredgefile(workflow config.Workflow) error
	AddBucketToDredgefile(bucket config.Bucket) error
	RelativePathFromDredgefile(path string) (string, error)
}

type Callbacks interface {
	UserInteractionCallbacks
	ExecutionCallbacks
	DredgefileCallbacks
}

type LogLevel int

const (
	Fatal LogLevel = iota
	Error
	Warn
	Info
	Debug
	Trace
)

func (l LogLevel) String() string {
	return [...]string{"Fatal", "Error", "Warn", "Info", "Debug", "Trace"}[l]
}

type NoResult struct{}

func (n *NoResult) Error() string {
	return "no result"
}

type InputType int

const (
	Text InputType = iota
	Select
)

func (i InputType) String() string {
	return [...]string{"Text", "Select"}[i]
}

type InputRequest struct {
	Name         string
	Description  string
	Type         InputType
	Values       []string
	DefaultValue string
}
