package yetenv

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEnvironment(t *testing.T) {
	t.Run("should return environment 'production'", func(t *testing.T) {
		values := []string{
			"production",
			"PRODUCTION",
		}

		for _, value := range values {
			value := value

			t.Run(fmt.Sprintf("ENVIRONMENT=%s", value), func(t *testing.T) {
				assert := assert.New(t)

				err := os.Setenv("ENVIRONMENT", value)
				require.NoError(t, err)

				actualEnv := GetEnvironment()
				assert.Equal(Production, actualEnv)
			})
		}
	})

	t.Run("should return environemnt 'staging'", func(t *testing.T) {
		values := []string{
			"staging",
			"STAGING",
		}

		for _, value := range values {
			value := value

			t.Run(fmt.Sprintf("ENVIRONMENT=%s", value), func(t *testing.T) {
				assert := assert.New(t)

				err := os.Setenv("ENVIRONMENT", value)
				require.NoError(t, err)

				actualEnv := GetEnvironment()
				assert.Equal(Staging, actualEnv)
			})
		}
	})

	t.Run("should return environemnt 'develop' for any other value of ENVIRONMENT", func(t *testing.T) {
		values := []string{
			"develop",
			"DEVELOP",
			"ANY",
		}

		for _, value := range values {
			value := value

			t.Run(fmt.Sprintf("ENVIRONMENT=%s", value), func(t *testing.T) {
				assert := assert.New(t)

				err := os.Setenv("ENVIRONMENT", value)
				require.NoError(t, err)

				actualEnv := GetEnvironment()
				assert.Equal(Develop, actualEnv)
			})
		}
	})

	t.Run("should return environemnt 'production' when changing DefaultVariableName", func(t *testing.T) {
		values := []string{
			"production",
			"PRODUCTION",
		}

		DefaultVariableName = "CHANGED_ENV"

		for _, value := range values {
			value := value

			t.Run(fmt.Sprintf("CHANGED_ENV=%s", value), func(t *testing.T) {
				assert := assert.New(t)

				err := os.Setenv("CHANGED_ENV", value)
				require.NoError(t, err)

				actualEnv := GetEnvironment()
				assert.Equal(Production, actualEnv)
			})
		}
	})
}

func TestGetEnvironmentFromVariable(t *testing.T) {
	customVariableName := "CUSTOM_ENV"

	t.Run("should return environment 'production' from variable 'CUSTOM_ENV'", func(t *testing.T) {
		values := []string{
			"PRODUCTION",
			"production",
		}

		for _, value := range values {
			value := value

			t.Run(fmt.Sprintf("%s=%s", customVariableName, value), func(t *testing.T) {
				assert := assert.New(t)

				err := os.Setenv(customVariableName, value)
				require.NoError(t, err)

				actualEnv := GetEnvironmentFromVariable(customVariableName)
				assert.Equal(Production, actualEnv)
			})
		}
	})

	t.Run("should return environment 'staging' from variable 'CUSTOM_ENV'", func(t *testing.T) {
		values := []string{
			"STAGING",
			"staging",
		}

		for _, value := range values {
			value := value

			t.Run(fmt.Sprintf("%s=%s", customVariableName, value), func(t *testing.T) {
				assert := assert.New(t)

				err := os.Setenv(customVariableName, value)
				require.NoError(t, err)

				actualEnv := GetEnvironmentFromVariable(customVariableName)
				assert.Equal(Staging, actualEnv)
			})
		}
	})

	t.Run("should return environment 'develop' from variable 'CUSTOM_ENV' from any other value", func(t *testing.T) {
		values := []string{
			"develop",
			"DEVELOP",
			"ANY",
		}

		for _, value := range values {
			value := value

			t.Run(fmt.Sprintf("%s=%s", customVariableName, value), func(t *testing.T) {
				assert := assert.New(t)

				err := os.Setenv(customVariableName, value)
				require.NoError(t, err)

				actualEnv := GetEnvironmentFromVariable(customVariableName)
				assert.Equal(Develop, actualEnv)
			})
		}
	})
}