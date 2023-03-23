package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddInput(t *testing.T) {
	providers, err := GetProviders()
	assert.Nil(t, err)
	assert.Equal(t, 4, len(providers))
}
