package workflow

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dredge-dev/dredge/internal/config"
)

func (workflow *Workflow) executeTemplate(step *config.TemplateStep) error {
	text, err := getTemplateText(workflow, step)
	if err != nil {
		return err
	}

	dest, err := workflow.Callbacks.Template(step.Dest)
	if err != nil {
		return fmt.Errorf("failed to template Dest: %s", err)
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

func getTemplateText(workflow *Workflow, step *config.TemplateStep) (string, error) {
	input := step.Input
	if step.Source != "" {
		path, err := workflow.Callbacks.RelativePathFromDredgefile((string)(step.Source))
		if err != nil {
			return "", err
		}
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			return "", err
		}
		input = string(buf)
	}
	return workflow.Callbacks.Template(input)
}
