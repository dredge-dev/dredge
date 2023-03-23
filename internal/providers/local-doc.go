package providers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dredge-dev/dredge/internal/api"
)

type LocalDocProvider struct {
	Path string
}

func (l *LocalDocProvider) Name() string {
	return "local-doc"
}

func (l *LocalDocProvider) Discover(callbacks api.Callbacks) error {
	for _, path := range []string{"docs", "documentation"} {
		info, err := os.Stat(path)
		if err == nil && info.IsDir() {
			confirmed, err := callbacks.Confirm("Local docs in detected in %s, do you want to add local-docs?", path)
			if err != nil {
				return err
			}
			if confirmed {
				err = callbacks.Log(api.Info, "Adding local-doc as a provider for %s", path)
				if err != nil {
					return err
				}
				err = callbacks.AddProviderToDredgefile("doc", "local-doc", map[string]string{
					"path": path,
				})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (l *LocalDocProvider) Init(config map[string]string) error {
	err := checkConfig(config, []string{"path"})
	if err != nil {
		return err
	}
	l.Path = config["path"]
	return nil
}

func (l *LocalDocProvider) ExecuteCommand(commandName string, c api.Callbacks) (interface{}, error) {
	if commandName == "search" {
		return l.Search(c)
	}
	return nil, fmt.Errorf("could not find command %s", commandName)
}

func (l *LocalDocProvider) Search(c api.Callbacks) (interface{}, error) {
	inputs, err := c.RequestInput([]api.InputRequest{
		{
			Name:        "text",
			Description: "Search text",
			Type:        api.Text,
		},
	})
	if err != nil {
		return nil, err
	}

	docs, err := l.search(inputs["text"])
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, doc := range docs {
		name := filepath.Base(doc)
		location, err := filepath.Abs(doc)
		if err != nil {
			return nil, err
		}
		result = append(result, map[string]interface{}{
			"name":     name,
			"author":   "",
			"location": location,
			"date":     "",
		})
	}
	return result, nil
}

func (l *LocalDocProvider) search(text string) ([]string, error) {
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("grep -R -i -l '%s' %s", text, l.Path))
	output, err := cmd.Output()
	if err != nil {
		eerr, ok := err.(*exec.ExitError)
		if ok && eerr.ExitCode() == 1 {
			return []string{}, nil
		} else {
			return nil, err
		}
	}
	return strings.Split(strings.TrimSuffix(string(output), "\n"), "\n"), nil
}
