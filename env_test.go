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

func TestConfigLoader_UseEnvironment(t *testing.T) {
	configLoader := NewConfigLoader()
	require.Equal(t, Environment(""), configLoader.Environment)

	configLoader.UseEnvironment(Staging)
	assert.Equal(t, Staging, configLoader.Environment)
}

func TestConfigLoader_UseLoadBehavior(t *testing.T) {
	configLoader := NewConfigLoader()
	require.Equal(t, LoadBehaviorUnknown, configLoader.LoadBehavior)

	configLoader.UseLoadBehavior(LoadBehaviorDefault)
	assert.Equal(t, LoadBehaviorDefault, configLoader.LoadBehavior)

	configLoader.UseLoadBehavior(LoadBehaviorCustom)
	assert.Equal(t, LoadBehaviorCustom, configLoader.LoadBehavior)
}

func TestConfigLoader_UseDefaultLoadBehavior(t *testing.T) {
	configLoader := NewConfigLoader()
	require.Equal(t, LoadBehaviorUnknown, configLoader.LoadBehavior)

	configLoader.UseDefaultLoadBehavior()
	assert.Equal(t, LoadBehaviorDefault, configLoader.LoadBehavior)
}

func TestConfigLoader_UseCustomLoadBehavior(t *testing.T) {
	configLoader := NewConfigLoader()
	require.Equal(t, LoadBehaviorUnknown, configLoader.LoadBehavior)

	configLoader.UseCustomLoadBehavior()
	assert.Equal(t, LoadBehaviorCustom, configLoader.LoadBehavior)
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

	assert.Equal(t, len(configLoader.loadOrder), 4)

	assert.Equal(t, loadOrderItemDevelop.file, configLoader.loadOrder[0].file)
	assert.NotNil(t, configLoader.loadOrder[0].conditionFunc)

	assert.Equal(t, loadOrderItemStaging.file, configLoader.loadOrder[1].file)
	assert.NotNil(t, configLoader.loadOrder[1].conditionFunc)

	assert.Equal(t, loadOrderItemProduction.file, configLoader.loadOrder[2].file)
	assert.NotNil(t, configLoader.loadOrder[2].conditionFunc)

	assert.Equal(t, loadOrderItemCustom.file, configLoader.loadOrder[3].file)
	assert.Nil(t, configLoader.loadOrder[3].conditionFunc)
}

func TestConfigLoader_LoadFromFile(t *testing.T) {
	configLoader := NewConfigLoader()
	require.Equal(t, len(configLoader.loadOrder), 0)

	configLoader.LoadFromFile("custom.env")

	customLoadOrderItem := loadOrderItem{
		file:          "custom.env",
		conditionFunc: nil,
	}

	assert.Len(t, configLoader.loadOrder, 1)
	assert.Equal(t, customLoadOrderItem.file, configLoader.loadOrder[0].file)
	assert.Nil(t, configLoader.loadOrder[0].conditionFunc)
}

func TestConfigLoader_LoadFromConditionalFile(t *testing.T) {
	configLoader := NewConfigLoader()
	require.Equal(t, len(configLoader.loadOrder), 0)

	configLoader.LoadFromConditionalFile("./cfg.dev.env", DefaultConditionForDevelopEnvironment)

	loadOrderItemDevelop := loadOrderItem{
		file:          "./cfg.dev.env",
		conditionFunc: DefaultConditionForDevelopEnvironment,
	}

	assert.Equal(t, len(configLoader.loadOrder), 1)
	assert.Equal(t, configLoader.loadOrder[0].file, loadOrderItemDevelop.file)
	assert.NotNil(t, configLoader.loadOrder[0].conditionFunc)
}

func TestConfigLoader_LoadInto(t *testing.T) {
	t.Run("should return error if load behavior is not set", func(t *testing.T) {
		resetEnv()

		c := testConfig{}
		err := NewConfigLoader().LoadInto(&c)

		assert.Error(t, err)
		assert.Equal(t, ErrUnknownLoadBehavior, err)
	})

	t.Run("default behavior", func(t *testing.T) {
		t.Run("should load dev and custom config when env=development", func(t *testing.T) {
			resetEnv()

			c := testConfig{}
			err := NewConfigLoader().
				UseLoadPath("./testdata").
				UseEnvironment(Develop).
				UseDefaultLoadBehavior().
				LoadInto(&c)

			expectedConfig := testConfig{
				Develop:    true,
				Staging:    false,
				Production: false,
				Custom:     true,
				LastFile:   "custom",
			}

			assert.NoError(t, err)
			assert.Equal(t, expectedConfig, c)
		})

		t.Run("should load staging and custom config when env=staging", func(t *testing.T) {
			resetEnv()

			c := testConfig{}
			err := NewConfigLoader().
				UseLoadPath("./testdata").
				UseEnvironment(Staging).
				UseDefaultLoadBehavior().
				LoadInto(&c)

			expectedConfig := testConfig{
				Develop:    false,
				Staging:    true,
				Production: false,
				Custom:     true,
				LastFile:   "custom",
			}

			assert.NoError(t, err)
			assert.Equal(t, expectedConfig, c)
		})

		t.Run("should load production and custom config when env=production", func(t *testing.T) {
			resetEnv()

			c := testConfig{}
			err := NewConfigLoader().
				UseLoadPath("./testdata").
				UseEnvironment(Production).
				UseDefaultLoadBehavior().
				LoadInto(&c)

			expectedConfig := testConfig{
				Develop:    false,
				Staging:    false,
				Production: true,
				Custom:     true,
				LastFile:   "custom",
			}

			assert.NoError(t, err)
			assert.Equal(t, expectedConfig, c)
		})

		t.Run("should load develop and custom config when env=production and file processor is yaml", func(t *testing.T) {
			resetEnv()

			c := testConfig{}
			err := NewConfigLoader().
				UseLoadPath("./testdata").
				UseEnvironment(Develop).
				UseFileProcessor(YAML).
				UseDefaultLoadBehavior().
				LoadInto(&c)

			expectedConfig := testConfig{
				Develop:    true,
				Staging:    false,
				Production: false,
				Custom:     true,
				LastFile:   "custom",
			}

			assert.NoError(t, err)
			assert.Equal(t, expectedConfig, c)
		})
	})
}

func resetEnv() {
	_ = os.Unsetenv("DEVELOP")
	_ = os.Unsetenv("STAGING")
	_ = os.Unsetenv("PROD")
	_ = os.Unsetenv("CUSTOM")
	_ = os.Unsetenv("LAST_FILE")
}

func TestNewConfigLoader(t *testing.T) {
	type Cfg struct {
	}

	c := Cfg{}

	_ = NewConfigLoader().
		UseLoadPath("./").
		UseFileProcessor(YAML).
		UseFileNameForEnvironment(Production, "prod.yaml").
		UseLoadBehavior(LoadBehaviorCustom).
		LoadFromFileForEnvironment(Develop).
		LoadFromFile("./bla.env").
		LoadInto(&c)
}
