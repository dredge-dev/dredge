package providers

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/dredge-dev/dredge/internal/resource"
)

type GithubReleasesProvider struct {
}

func (g *GithubReleasesProvider) Name() string {
	return "github-releases"
}

func (g *GithubReleasesProvider) Init(config map[string]string) error {
	return nil
}

func (g *GithubReleasesProvider) ExecuteCommand(commandName string, callbacks resource.Callbacks) (interface{}, error) {
	if commandName == "get" {
		return g.Get(callbacks)
	} else if commandName == "search" {
		return g.Search(callbacks)
	} else if commandName == "describe" {
		return g.Describe(callbacks)
	}
	return nil, fmt.Errorf("could not find command %s", commandName)
}

func (g *GithubReleasesProvider) Get(callbacks resource.Callbacks) (interface{}, error) {
	cmd := exec.Command("/bin/bash", "-c", "SHELL=/bin/bash gh release list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no output from gh")
	}
	titleIndex := strings.Index(lines[0], "TITLE")
	typeIndex := strings.Index(lines[0], "TYPE")
	tagIndex := strings.Index(lines[0], "TAG NAME")
	publishedIndex := strings.Index(lines[0], "PUBLISHED")
	if titleIndex != 0 || typeIndex <= titleIndex || tagIndex <= typeIndex || publishedIndex <= tagIndex {
		return nil, fmt.Errorf("format error in gh output lines[0]:%s, titleIndex:%d typeIndex:%d tagIndex:%d publishedIndex:%d", lines, titleIndex, typeIndex, tagIndex, publishedIndex)
	}
	var out []map[string]interface{}
	for _, line := range lines[1:] {
		v := map[string]interface{}{
			"title": strings.Trim(line[titleIndex:typeIndex], " "),
			"name":  strings.Trim(line[tagIndex:publishedIndex], " "),
			"date":  strings.Trim(line[publishedIndex:], " "),
		}
		out = append(out, v)
	}
	return out, nil
}

func (g *GithubReleasesProvider) Search(callbacks resource.Callbacks) (interface{}, error) {
	inputs, err := callbacks.RequestInput([]resource.InputRequest{
		{
			Name:        "text",
			Description: "Search text",
			Type:        resource.Text,
		},
	})
	if err != nil {
		return nil, err
	}
	var ret []map[string]interface{}
	ret = append(ret, map[string]interface{}{
		"name":  inputs["text"],
		"date":  "yesterday",
		"title": "first version",
		"notes": "https://github.com/dredge-dev/dredge/releases/tag/v0.0.6",
	})
	return ret, nil
}

func (g *GithubReleasesProvider) Describe(callbacks resource.Callbacks) (interface{}, error) {
	inputs, err := callbacks.RequestInput([]resource.InputRequest{
		{
			Name:        "name",
			Description: "Release name",
			Type:        resource.Text,
		},
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"name":  inputs["name"],
		"date":  "yesterday",
		"title": "first version",
		"notes": "https://github.com/dredge-dev/dredge/releases/tag/v0.0.6",
	}, nil
}
