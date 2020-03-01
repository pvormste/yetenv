package yetenv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvVariables_Count(t *testing.T) {
	t.Run("should return 0 when map is empty", func(t *testing.T) {
		variables := EnvVariables{}
		assert.Equal(t, 0, variables.Count())
	})

	t.Run("should return the correct size when map is not empty", func(t *testing.T) {
		variables := EnvVariables{
			"VARIABLE1": "value1",
			"VARIABLE2": "value2",
		}

		assert.Equal(t, 2, variables.Count())
	})
}
