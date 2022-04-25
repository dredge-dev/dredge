package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHome(t *testing.T) {
	tests := map[string]struct {
		runtime Runtime
		home    string
	}{
		"default home": {
			runtime: Runtime{Home: ""},
			home:    DEFAULT_HOME,
		},
		"custom home": {
			runtime: Runtime{Home: "/my-home"},
			home:    "/my-home",
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		home := test.runtime.GetHome()
		assert.Equal(t, test.home, home)
	}
}

func TestGetValue(t *testing.T) {
	input := Input{
		Name:        "city",
		Description: "city",
		Type:        "select",
		Values:      []string{"Brussels", "Barcelona"},
	}
	assert.Equal(t, true, input.HasValue("Brussels"))
	assert.Equal(t, true, input.HasValue("Barcelona"))
	assert.Equal(t, false, input.HasValue("London"))
}
