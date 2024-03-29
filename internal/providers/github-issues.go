package providers

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/dredge-dev/dredge/internal/api"
)

var URL_RE = regexp.MustCompile(`/issues/([0-9]+)`)

type GithubIssuesProvider struct {
}

func (g *GithubIssuesProvider) Name() string {
	return "github-issues"
}

func (g *GithubIssuesProvider) Discover(callbacks api.Callbacks) error {
	cmd := exec.Command("/bin/bash", "-c", "git remote -v")
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	if strings.Contains(string(output), "github.com") {
		confirmed, err := callbacks.Confirm("GitHub repo detected, do you use GitHub issues?")
		if err != nil {
			return err
		}
		if confirmed {
			err = callbacks.Log(api.Info, "Adding github-issues as a provider")
			if err != nil {
				return err
			}
			err = callbacks.AddProviderToDredgefile("issue", "github-issues", nil)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *GithubIssuesProvider) Init(config map[string]string) error {
	return nil
}

func (g *GithubIssuesProvider) ExecuteCommand(commandName string, callbacks api.Callbacks) (interface{}, error) {
	if commandName == "get" {
		return g.Get(callbacks)
	} else if commandName == "create" {
		return g.Create(callbacks)
	}
	return nil, fmt.Errorf("could not find command %s", commandName)
}

type GithubIssue struct {
	Author    GithubAuthor
	Labels    []GithubLabel
	CreatedAt string
	Number    int
	State     string
	Title     string
}

type GithubAuthor struct {
	Login string
}

type GithubLabel struct {
	Name string
}

func (g *GithubIssuesProvider) Get(callbacks api.Callbacks) (interface{}, error) {
	cmd := exec.Command("/bin/bash", "-c", "gh issue list --json number,title,author,state,createdAt,labels")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var issues []GithubIssue
	err = json.Unmarshal(output, &issues)
	if err != nil {
		return nil, err
	}
	var out []map[string]interface{}
	for _, issue := range issues {
		issueType := "issue"
		for _, label := range issue.Labels {
			if label.Name == "bug" {
				issueType = "bug"
			}
			if label.Name == "enhancement" {
				issueType = "feature"
			}
		}
		i := map[string]interface{}{
			"name":  fmt.Sprintf("%d", issue.Number),
			"title": issue.Title,
			"type":  issueType,
			"state": issue.State,
			"date":  issue.CreatedAt,
		}
		out = append(out, i)
	}
	return out, nil
}

func (g *GithubIssuesProvider) Create(c api.Callbacks) (interface{}, error) {
	inputs, err := c.RequestInput([]api.InputRequest{
		{
			Name:        "title",
			Description: "",
			Type:        api.Text,
		},
		{
			Name:         "type",
			Description:  "",
			Type:         api.Select,
			Values:       []string{"bug", "feature"},
			DefaultValue: "bug",
		},
		{
			Name:        "description",
			Description: "",
			Type:        api.Text,
		},
	})
	if err != nil {
		return nil, err
	}
	title := inputs["title"]
	body := inputs["description"]
	label := inputs["type"]
	if label == "feature" {
		label = "enhancement"
	}
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("gh issue create --title '%s' --body '%s' --label '%s'", title, body, label))
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	p := URL_RE.FindStringSubmatch(string(output))
	if len(p) < 1 {
		return nil, fmt.Errorf("format error in gh output")
	}
	name := p[1]
	return map[string]interface{}{
		"name":  name,
		"title": inputs["title"],
		"type":  inputs["type"],
		"state": "open",
		"date":  "now",
	}, nil
}
