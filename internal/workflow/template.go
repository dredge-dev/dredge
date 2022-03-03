package workflow

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"
	"time"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/exec"
)

const TEMPLATE_NAME = "root"

var TEMPLATE_FUNCTIONS = template.FuncMap{
	"replace": func(s, old, new string) string {
		return strings.Replace(s, old, new, -1)
	},
	"date": func(format string) string {
		return time.Now().Format(format)
	},
}

func executeTemplate(workflow *exec.Workflow, step *config.TemplateStep) error {
	env := workflow.Exec.Env

	input := step.Input
	if step.Source != "" {
		buf, err := workflow.Exec.ReadSource(step.Source)
		if err != nil {
			return err
		}
		input = string(buf)
	}

	out, err := Template(input, env)
	if err != nil {
		return err
	}

	dest, err := Template(step.Dest, env)
	if err != nil {
		return fmt.Errorf("Failed to template Dest: %s", err)
	}

	return ioutil.WriteFile(dest, []byte(out), 0644)
}

func Template(input string, env exec.Env) (string, error) {
	t, err := template.New(TEMPLATE_NAME).Funcs(TEMPLATE_FUNCTIONS).Parse(string(input))
	if err != nil {
		return "", fmt.Errorf("Failed to parse template: %s", err)
	}

	var buffer bytes.Buffer
	if err := t.Execute(&buffer, env); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
