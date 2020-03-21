package yetenv

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testConfig struct {
	Develop    bool   `yaml:"develop" env:"DEVELOP"`
	Staging    bool   `yaml:"staging" env:"STAGING"`
	Production bool   `yaml:"prod" env:"PROD"`
	Custom     bool   `yaml:"custom" env:"CUSTOM"`
	LastFile   string `yaml:"last_file" env:"LAST_FILE"`
}

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

func TestConfigLoader_UseLoadPath(t *testing.T) {
	configLoader := NewConfigLoader()
	require.Equal(t, "./", configLoader.LoadPath)

	configLoader.UseLoadPath("./testdata")
	assert.Equal(t, "./testdata", configLoader.LoadPath)
}

func TestConfigLoader_UseFileProcessor(t *testing.T) {
	configLoader := NewConfigLoader()
	require.Equal(t, DOTENV, configLoader.FileExtension)

	configLoader.UseFileProcessor(YAML)
	require.Equal(t, YAML, configLoader.FileExtension)
}

func TestConfigLoader_UseFileNameForEnvironment(t *testing.T) {
	configLoader := NewConfigLoader()
	require.Equal(t, defaultDevelopConfigFile, configLoader.ConfigFiles[Develop])
	require.Equal(t, defaultStagingConfigFile, configLoader.ConfigFiles[Staging])
	require.Equal(t, defaultProductionConfigFile, configLoader.ConfigFiles[Production])
	require.Equal(t, defaultCustomConfigFile, configLoader.ConfigFiles[Custom])

	configLoader.
		UseFileNameForEnvironment(Develop, "my.dev.yaml").
		UseFileNameForEnvironment(Staging, "my.staging.env").
		UseFileNameForEnvironment(Production, "my.prod.toml").
		UseFileNameForEnvironment(Custom, "my.json")

	assert.Equal(t, "my.dev", configLoader.ConfigFiles[Develop])
	assert.Equal(t, "my.staging", configLoader.ConfigFiles[Staging])
	assert.Equal(t, "my.prod", configLoader.ConfigFiles[Production])
	assert.Equal(t, "my", configLoader.ConfigFiles[Custom])

	configLoader.UseFileNameForEnvironment(Develop, "your.dev")
	assert.Equal(t, "your.dev", configLoader.ConfigFiles[Develop])
}

func TestConfigLoader_UseCustomLoadBehavior(t *testing.T) {
	configLoader := NewConfigLoader()
	require.False(t, configLoader.CustomLoadBehavior)

	configLoader.UseCustomLoadBehavior()
	assert.True(t, configLoader.CustomLoadBehavior)
}

func TestConfigLoader_LoadFromFileForEnvironment(t *testing.T) {
	configLoader := NewConfigLoader()
	require.Equal(t, len(configLoader.loadOrder), 0)

	configLoader.LoadFromFileForEnvironment(Develop).
		LoadFromFileForEnvironment(Staging).
		LoadFromFileForEnvironment(Production).
		LoadFromFileForEnvironment(Custom)

	loadOrderItemDevelop := loadOrderItem{
		file:          "cfg.dev.env",
		conditionFunc: DefaultConditionForDevelopEnvironment,
	}

	loadOrderItemStaging := loadOrderItem{
		file:          "cfg.staging.env",
		conditionFunc: DefaultConditionForStagingEnvironment,
	}

	loadOrderItemProduction := loadOrderItem{
		file:          "cfg.prod.env",
		conditionFunc: DefaultConditionForProductionEnvironment,
	}

	loadOrderItemCustom := loadOrderItem{
		file:          ".env",
		conditionFunc: nil,
	}

	assert.Equal(t, loadOrderItemDevelop.file, configLoader.loadOrder[0].file)
	assert.NotNil(t, configLoader.loadOrder[0].conditionFunc)

	assert.Equal(t, loadOrderItemStaging.file, configLoader.loadOrder[1].file)
	assert.NotNil(t, configLoader.loadOrder[1].conditionFunc)

	assert.Equal(t, loadOrderItemProduction.file, configLoader.loadOrder[2].file)
	assert.NotNil(t, configLoader.loadOrder[2].conditionFunc)

	assert.Equal(t, loadOrderItemCustom.file, configLoader.loadOrder[3].file)
	assert.Nil(t, configLoader.loadOrder[3].conditionFunc)
}

func TestConfigLoader_LoadFromConditionalFile(t *testing.T) {
	configLoader := NewConfigLoader()
	require.Equal(t, len(configLoader.loadOrder), 0)

	configLoader.LoadFromConditionalFile("./cfg.dev.env", DefaultConditionForDevelopEnvironment)

	loadOrderItemDevelop := loadOrderItem{
		file:          "./cfg.dev.env",
		conditionFunc: DefaultConditionForDevelopEnvironment,
	}

	assert.Equal(t, configLoader.loadOrder[0].file, loadOrderItemDevelop.file)
	assert.NotNil(t, configLoader.loadOrder[0].conditionFunc)
}

func TestNewConfigLoader(t *testing.T) {
	type Cfg struct {
	}

	c := Cfg{}

	_ = NewConfigLoader().
		UseLoadPath("./").
		UseFileProcessor(YAML).
		UseFileNameForEnvironment(Production, "prod.yaml").
		UseCustomLoadBehavior().
		LoadFromFileForEnvironment(Develop).
		LoadFromFile("./bla.env").
		LoadInto(&c)
}
