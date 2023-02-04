package resource

type Callbacks interface {
	Log(level LogLevel, msg string) error
	RequestInput(inputRequests []InputRequest) (map[string]string, error)
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
