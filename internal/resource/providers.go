package resource

import (
	"fmt"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/providers"
)

func CreateProvider(conf config.ResourceProvider) (Provider, error) {
	p, err := getProvider(conf.Provider)
	if err != nil {
		return nil, err
	}
	err = p.Init(conf.Config)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func getProvider(name string) (Provider, error) {
	if name == "github-releases" {
		return &providers.GithubReleasesProvider{}, nil
	}
	if name == "github-issues" {
		return &providers.GithubIssuesProvider{}, nil
	}
	if name == "local-doc" {
		return &providers.LocalDocProvider{}, nil
	}
	if name == "local-docker-compose" {
		return &providers.LocalDockerComposeProvider{}, nil
	}
	return nil, fmt.Errorf("could not find provider %s", name)
}
