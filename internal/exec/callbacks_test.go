package exec

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestAddProviderToDredgefile(t *testing.T) {
	file := "./add-provider-test-dredgefile"
	source := config.SourcePath(file)
	os.Remove(file)

	tests := map[string]struct {
		e              *DredgeExec
		resource       string
		provider       string
		providerConfig map[string]string
		errMsg         string
		output         *config.DredgeFile
	}{
		"empty resource": {
			e:              &DredgeExec{},
			resource:       "",
			provider:       "my-release-provider",
			providerConfig: nil,
			errMsg:         "empty resource cannot be added to Dredgefile",
		},
		"empty provider": {
			e:              &DredgeExec{},
			resource:       "release",
			provider:       "",
			providerConfig: nil,
			errMsg:         "empty provider cannot be added to Dredgefile",
		},
		"with config": {
			e: &DredgeExec{
				Source:     source,
				DredgeFile: &config.DredgeFile{},
			},
			resource: "release",
			provider: "my-release-provider",
			providerConfig: map[string]string{
				"url": "test.com",
			},
			output: &config.DredgeFile{
				Resources: map[string]config.Resource{
					"release": []config.ResourceProvider{
						{
							Provider: "my-release-provider",
							Config: map[string]string{
								"url": "test.com",
							},
						},
					},
				},
			},
		},
		"without config": {
			e: &DredgeExec{
				Source:     source,
				DredgeFile: &config.DredgeFile{},
			},
			resource:       "release",
			provider:       "my-release-provider",
			providerConfig: nil,
			output: &config.DredgeFile{
				Resources: map[string]config.Resource{
					"release": []config.ResourceProvider{
						{
							Provider: "my-release-provider",
							Config:   nil,
						},
					},
				},
			},
		},
		"new resource": {
			e: &DredgeExec{
				Source: source,
				DredgeFile: &config.DredgeFile{
					Resources: map[string]config.Resource{
						"doc": []config.ResourceProvider{
							{
								Provider: "my-docs-provider",
								Config:   nil,
							},
						},
					},
				},
			},
			resource: "release",
			provider: "my-release-provider",
			providerConfig: map[string]string{
				"url": "test.com",
			},
			output: &config.DredgeFile{
				Resources: map[string]config.Resource{
					"doc": []config.ResourceProvider{
						{
							Provider: "my-docs-provider",
							Config:   nil,
						},
					},
					"release": []config.ResourceProvider{
						{
							Provider: "my-release-provider",
							Config: map[string]string{
								"url": "test.com",
							},
						},
					},
				},
			},
		},
		"existing resource": {
			e: &DredgeExec{
				Source: source,
				DredgeFile: &config.DredgeFile{
					Resources: map[string]config.Resource{
						"release": []config.ResourceProvider{
							{
								Provider: "existing-provider",
								Config:   nil,
							},
						},
					},
				},
			},
			resource: "release",
			provider: "my-release-provider",
			providerConfig: map[string]string{
				"url": "test.com",
			},
			output: &config.DredgeFile{
				Resources: map[string]config.Resource{
					"release": []config.ResourceProvider{
						{
							Provider: "existing-provider",
							Config:   nil,
						},
						{
							Provider: "my-release-provider",
							Config: map[string]string{
								"url": "test.com",
							},
						},
					},
				},
			},
		},
		"existing provider": {
			e: &DredgeExec{
				Source: source,
				DredgeFile: &config.DredgeFile{
					Resources: map[string]config.Resource{
						"release": []config.ResourceProvider{
							{
								Provider: "existing-provider",
								Config:   nil,
							},
						},
					},
				},
			},
			resource:       "release",
			provider:       "existing-provider",
			providerConfig: nil,
			errMsg:         "provider 'existing-provider' already defined",
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		err := test.e.AddProviderToDredgefile(test.resource, test.provider, test.providerConfig)
		if test.errMsg == "" {
			assert.Nil(t, err)
			actual, err := ioutil.ReadFile(file)
			if err != nil {
				panic(err)
			}
			err = config.WriteDredgeFile(test.output, source)
			if err != nil {
				panic(err)
			}
			expected, err := ioutil.ReadFile(file)
			if err != nil {
				panic(err)
			}
			assert.Equal(t, string(expected), string(actual))
		} else {
			assert.Equal(t, test.errMsg, fmt.Sprint(err))
		}
		os.Remove(file)
	}
}
