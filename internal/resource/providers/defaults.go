package providers

import (
	"fmt"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/exec"
	"github.com/dredge-dev/dredge/internal/resource"
)

func CreateProvider(de *exec.DredgeExec, provider config.ResourceProvider) (resource.Provider, error) {
	p, err := getProvider(provider.Provider)
	if err != nil {
		return nil, err
	}
	err = p.Init(provider.Config)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func getProvider(name string) (resource.Provider, error) {
	if name == "github-releases" {
		return &GithubReleasesProvider{}, nil
	}
	if name == "github-issues" {
		return &GithubIssuesProvider{}, nil
	}
	if name == "local-doc" {
		return &LocalDocProvider{}, nil
	}
	if name == "local-docker-compose" {
		return &LocalDockerComposeProvider{}, nil
	}
	return nil, fmt.Errorf("could not find provider %s", name)
}

func checkConfig(config map[string]string, configKeys []string) error {
	for _, key := range configKeys {
		if _, ok := config[key]; !ok {
			return fmt.Errorf("could not find field %s in config", key)
		}
	}
	return nil
}
