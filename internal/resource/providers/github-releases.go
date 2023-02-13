package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"

	"github.com/creack/pty"
	"github.com/dredge-dev/dredge/internal/callbacks"
)

var TERM_CHARS_RE = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))")

type GithubReleasesProvider struct {
}

func (g *GithubReleasesProvider) Name() string {
	return "github-releases"
}

func (g *GithubReleasesProvider) Init(config map[string]string) error {
	return nil
}

func (g *GithubReleasesProvider) ExecuteCommand(commandName string, c callbacks.Callbacks) (interface{}, error) {
	if commandName == "get" {
		return g.Get(c)
	} else if commandName == "describe" {
		return g.Describe(c)
	}
	return nil, fmt.Errorf("could not find command %s", commandName)
}

func (g *GithubReleasesProvider) Get(c callbacks.Callbacks) (interface{}, error) {
	cmd := exec.Command("/bin/bash", "-c", "gh release list")
	f, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	output := TERM_CHARS_RE.ReplaceAllString(string(bytes), "")
	lines := strings.Split(string(output), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no output from gh")
	}
	titleIndex := strings.Index(lines[0], "TITLE")
	typeIndex := strings.Index(lines[0], "TYPE")
	tagIndex := strings.Index(lines[0], "TAG NAME")
	publishedIndex := strings.Index(lines[0], "PUBLISHED")
	if titleIndex != 0 || typeIndex <= titleIndex || tagIndex <= typeIndex || publishedIndex <= tagIndex {
		return nil, fmt.Errorf("format error in gh output")
	}
	var out []map[string]interface{}
	for _, line := range lines[1:] {
		if len(line) > publishedIndex {
			v := map[string]interface{}{
				"title": strings.Trim(line[titleIndex:typeIndex], " "),
				"name":  strings.Trim(line[tagIndex:publishedIndex], " "),
				"date":  strings.Trim(line[publishedIndex:len(line)-1], " "),
			}
			out = append(out, v)
		}
	}
	return out, nil
}

type GithubRelease struct {
	Author      GithubAuthor
	PublishedAt string
	Url         string
	Name        string
	Body        string
}

func (g *GithubReleasesProvider) Describe(c callbacks.Callbacks) (interface{}, error) {
	inputs, err := c.RequestInput([]callbacks.InputRequest{
		{
			Name:        "name",
			Description: "Name",
			Type:        callbacks.Text,
		},
	})
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("gh release view '%s' --json name,body,author,publishedAt,url", inputs["name"]))
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var release GithubRelease
	err = json.Unmarshal(output, &release)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"name":        release.Name,
		"description": release.Body,
		"url":         release.Url,
		"date":        release.PublishedAt,
		"author":      release.Author.Login,
	}, nil
}
