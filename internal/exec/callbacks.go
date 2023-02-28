package exec

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/dredge-dev/dredge/internal/api"
	"github.com/dredge-dev/dredge/internal/config"
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
			inputs[inputRequest.Name] = fmt.Sprintf("%v", value)
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

func (e *DredgeExec) Confirm(msg string) error {
	return e.callbacks.Confirm(msg)
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

func (e *DredgeExec) SetEnv(name string, value interface{}) error {
	e.Env[name] = value
	return nil
}

var TEMPLATE_FUNCTIONS = template.FuncMap{
	"replace": func(s, old, new string) string {
		return strings.Replace(s, old, new, -1)
	},
	"date": func(format string) string {
		return time.Now().Format(format)
	},
	"join": func(s1, s2, sep string) string {
		if len(s1) == 0 {
			return s2
		}
		if len(s2) == 0 {
			return s1
		}
		return s1 + sep + s2
	},
	"trimSpace": func(s string) string {
		return strings.TrimSpace(s)
	},
	"isTrue":  isTrue,
	"isFalse": isFalse,
}

func isTrue(s string) bool {
	l := strings.ToLower(s)
	return l == "1" || l == "t" || l == "true" || l == "yes"
}

func isFalse(s string) bool {
	l := strings.ToLower(s)
	return l == "0" || l == "f" || l == "false" || l == "no"
}

func (e *DredgeExec) Template(input string) (string, error) {
	t, err := template.New("").Option("missingkey=zero").Funcs(TEMPLATE_FUNCTIONS).Parse(string(input))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %s", err)
	}

	var buffer bytes.Buffer
	if err := t.Execute(&buffer, e.Env); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (e *DredgeExec) AddVariablesToDredgefile(variables map[string]string) error {
	rootExec, df := e.getRootExecAndDredgeFile()

	if df.Variables == nil {
		df.Variables = make(config.Variables)
	}
	for variable, value := range variables {
		if _, ok := df.Variables[variable]; !ok {
			df.Variables[variable] = value
		} else {
			return fmt.Errorf("variable %s already present", variable)
		}
	}

	return validateAndWriteDredgefile(df, rootExec.Source)
}

func validateAndWriteDredgefile(df *config.DredgeFile, path config.SourcePath) error {
	if err := df.Validate(); err != nil {
		return err
	}
	return config.WriteDredgeFile(df, path)
}

func (e *DredgeExec) AddWorkflowToDredgefile(w config.Workflow) error {
	rootExec, df := e.getRootExecAndDredgeFile()

	if w.Import != nil {
		w.Import.Source = MergeSources(e.Source, w.Import.Source)
	}
	if f, _ := rootExec.GetWorkflow("", w.Name); f != nil {
		return fmt.Errorf("workflow %s already present", w.Name)
	} else {
		df.Workflows = append(df.Workflows, w)
	}

	return validateAndWriteDredgefile(df, rootExec.Source)
}

func (e *DredgeExec) AddBucketToDredgefile(b config.Bucket) error {
	rootExec, df := e.getRootExecAndDredgeFile()

	if b.Import != nil {
		b.Import.Source = MergeSources(e.Source, b.Import.Source)
	}
	if f, _ := rootExec.GetBucket(b.Name); f != nil {
		return fmt.Errorf("bucket %s already present", b.Name)
	} else {
		df.Buckets = append(df.Buckets, b)
	}

	return validateAndWriteDredgefile(df, rootExec.Source)
}

func (e *DredgeExec) RelativePathFromDredgefile(path string) (string, error) {
	return resolvePath(MergeSources(e.Source, config.SourcePath(path)))
}
