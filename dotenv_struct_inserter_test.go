package yetenv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStructEnvInserter_DetermineEnvValue(t *testing.T) {
	t.Run("when struct tag is NOT set", func(t *testing.T) {
		t.Run("should return empty string when value for field cant be found", func(t *testing.T) {
			inserter := newStructEnvInserter()
			envValue := inserter.getEnvValue(EnvVariables{"NO_FIELD": "value"}, "field", "")
			assert.Equal(t, "", envValue)
		})

		t.Run("should load from variables map by uppercase field name", func(t *testing.T) {
			inserter := newStructEnvInserter()
			envValue := inserter.getEnvValue(EnvVariables{"FIELD": "value"}, "field", "")
			assert.Equal(t, "value", envValue)
		})

		t.Run("should load from OS by uppercase field name", func(t *testing.T) {
			err := os.Setenv("FIELD", "value")
			require.NoError(t, err)

			inserter := newStructEnvInserter()
			envValue := inserter.getEnvValue(EnvVariables{"FIELD": "unused value"}, "field", "")
			assert.Equal(t, "value", envValue)
		})
	})

	t.Run("when struct tag is set", func(t *testing.T) {
		t.Run("should return empty string when value for field cant be found", func(t *testing.T) {
			inserter := newStructEnvInserter()
			envValue := inserter.getEnvValue(EnvVariables{"NO_FIELD": "value"}, "field", "ENV_NAME")
			assert.Equal(t, "", envValue)
		})

		t.Run("should load from variables map by struct tag", func(t *testing.T) {
			inserter := newStructEnvInserter()
			envValue := inserter.getEnvValue(EnvVariables{"ENV_NAME": "value"}, "field", "ENV_NAME")
			assert.Equal(t, "value", envValue)
		})

		t.Run("should load from OS by struct tag", func(t *testing.T) {
			err := os.Setenv("ENV_NAME", "value")
			require.NoError(t, err)

			inserter := newStructEnvInserter()
			envValue := inserter.getEnvValue(EnvVariables{"ENV_NAME": "unused value"}, "field", "ENV_NAME")
			assert.Equal(t, "value", envValue)
		})
	})
}
