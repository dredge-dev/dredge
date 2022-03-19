package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHome(t *testing.T) {
	customHome := "/my-home"

	tests := map[string]struct {
		runtime Runtime
		home    string
	}{
		"default home": {
			runtime: Runtime{Home: nil},
			home:    DEFAULT_HOME,
		},
		"custom home": {
			runtime: Runtime{Home: &customHome},
			home:    customHome,
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
