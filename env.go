package yetenv

import (
	"os"
	"strings"
)

const (
	environmentVariableValueProduction = "production"
	environmentVariableValueStaging    = "staging"

	defaultDevelopConfigFile    = "cfg.dev"
	defaultStagingConfigFile    = "cfg.staging"
	defaultProductionConfigFile = "cfg.prod"
	defaultCustomConfigFile     = "cfg"
)

type Environment string

const (
	Production Environment = "production"
	Staging    Environment = "staging"
	Develop    Environment = "develop"
	Custom     Environment = "custom"
)

type ConfigFileExtension string

const (
	YAML   ConfigFileExtension = ".yaml"
	JSON   ConfigFileExtension = ".json"
	TOML   ConfigFileExtension = ".toml"
	DOTENV ConfigFileExtension = ".env"
)

type ConditionalLoadFunc func(configLoader *ConfigLoader, currentEnvironment Environment) bool

// DefaultVariableName defines the default name of the environment variable.
var DefaultVariableName = "ENVIRONMENT"

// GetEnvironment returns the current Environment value depending on the OS environment
// value of the variable defined by DefaultVariableName.
func GetEnvironment() Environment {
	envRaw := os.Getenv(DefaultVariableName)
	return environmentFromVariableValue(envRaw)
}

// GetEnvironmentFromVariable returns the current Environment value depending on the OS environment
// value of the variable provided by the parameter.
func GetEnvironmentFromVariable(variableName string) Environment {
	envRaw := os.Getenv(variableName)
	return environmentFromVariableValue(envRaw)
}

func environmentFromVariableValue(variableValue string) Environment {
	env := strings.ToLower(variableValue)

	switch env {
	case environmentVariableValueProduction:
		return Production
	case environmentVariableValueStaging:
		return Staging
	}

	return Develop
}

type ConfigLoader struct {
	LoadPath           string
	FileExtension      ConfigFileExtension
	ConfigFiles        map[Environment]string
	CustomLoadBehavior bool
}

func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{
		LoadPath:      "./",
		FileExtension: DOTENV,
		ConfigFiles: map[Environment]string{
			Develop:    defaultDevelopConfigFile,
			Staging:    defaultStagingConfigFile,
			Production: defaultProductionConfigFile,
			Custom:     defaultCustomConfigFile,
		},
		CustomLoadBehavior: false,
	}
}

func (c *ConfigLoader) UseLoadPath(path string) *ConfigLoader {
	return c
}

func (c *ConfigLoader) UseFileProcessor(extension ConfigFileExtension) *ConfigLoader {
	return c
}

func (c *ConfigLoader) UseFileNameForEnvironment(environment Environment, fileName string) *ConfigLoader {
	return c
}

func (c *ConfigLoader) UseCustomLoadBehavior() *ConfigLoader {
	return nil
}

func (c *ConfigLoader) AddFileForEnvironment(environment Environment) *ConfigLoader {
	return nil
}

func (c *ConfigLoader) AddFile(filePath string) *ConfigLoader {
	return nil
}

func (c *ConfigLoader) AddConditionalFile(filePath string, loadFunc ConditionalLoadFunc) *ConfigLoader {
	return nil
}

func (c *ConfigLoader) LoadInto(cfg interface{}) error {
	return nil
}
