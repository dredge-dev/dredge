package resource

import (
	"fmt"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/providers"
)

var PROVIDERS = map[string]Provider{
	"github-releases":      &providers.GithubReleasesProvider{},
	"github-issues":        &providers.GithubIssuesProvider{},
	"local-doc":            &providers.LocalDocProvider{},
	"local-docker-compose": &providers.LocalDockerComposeProvider{},
}

func GetProviders() ([]Provider, error) {
	var providers []Provider
	for _, provider := range PROVIDERS {
		providers = append(providers, provider)
	}
	return providers, nil
}

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
	if provider, ok := PROVIDERS[name]; ok {
		return provider, nil
	}
	return nil, fmt.Errorf("could not find provider %s", name)
}
