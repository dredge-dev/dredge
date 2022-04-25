package workflow

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/exec"
)

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

func executeTemplate(workflow *exec.Workflow, step *config.TemplateStep) error {
	text, err := getTemplateText(workflow, step)
	if err != nil {
		return err
	}

	dest, err := Template(step.Dest, workflow.Exec.Env)
	if err != nil {
		return fmt.Errorf("Failed to template Dest: %s", err)
	}

	return insert(step.Insert, text, dest)
}

func insert(insert *config.Insert, text string, dest string) error {
	if insert == nil {
		return ioutil.WriteFile(dest, []byte(text), 0644)
	}

	currentContent, err := readFileIfExists(dest)
	if err != nil {
		return err
	}

	if insert.Section == "" {
		if len(currentContent) == 0 {
			return ioutil.WriteFile(dest, []byte(text), 0644)
		} else if insert.Placement == "" || insert.Placement == config.INSERT_END {
			return ioutil.WriteFile(dest, []byte(currentContent+"\n"+text), 0644)
		} else if insert.Placement == config.INSERT_BEGIN {
			return ioutil.WriteFile(dest, []byte(text+"\n"+currentContent), 0644)
		} else if insert.Placement == config.INSERT_UNIQUE {
			found := false
			for _, line := range strings.Split(currentContent, "\n") {
				if text == line {
					found = true
				}
			}
			if !found {
				return ioutil.WriteFile(dest, []byte(currentContent+"\n"+text), 0644)
			}
			return nil
		}
	}

	ext := getExtension(dest)
	if ext == "go" {
		output, err := insertGo(insert, currentContent, text)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(dest, []byte(output), 0644)
	} else {
		return fmt.Errorf("unsupported extension %s for insert (valid values: go)", ext)
	}
}

func readFileIfExists(src string) (string, error) {
	_, err := os.Stat(src)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		} else {
			return "", err
		}
	} else {
		bytes, err := ioutil.ReadFile(src)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}
}

func getExtension(dest string) string {
	parts := strings.Split(dest, ".")
	return parts[len(parts)-1]
}

func getTemplateText(workflow *exec.Workflow, step *config.TemplateStep) (string, error) {
	input := step.Input
	if step.Source != "" {
		buf, err := workflow.Exec.ReadSource(step.Source)
		if err != nil {
			return "", err
		}
		input = string(buf)
	}

	return Template(input, workflow.Exec.Env)
}

func Template(input string, env exec.Env) (string, error) {
	t, err := template.New("").Option("missingkey=zero").Funcs(TEMPLATE_FUNCTIONS).Parse(string(input))
	if err != nil {
		return "", fmt.Errorf("Failed to parse template: %s", err)
	}

	var buffer bytes.Buffer
	if err := t.Execute(&buffer, env); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
