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
}