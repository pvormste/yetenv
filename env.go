package yetenv

import (
	"os"
	"strings"

	"github.com/pvormste/yeterr"
)

// DefaultVariableName defines the default name of the environment variable.
var DefaultVariableName = "ENVIRONMENT"

const (
	environmentVariableValueProduction = "production"
	environmentVariableValueStaging    = "staging"
)

type Environment string

const (
	Production Environment = "production"
	Staging    Environment = "staging"
	Develop    Environment = "develop"

	DefaultEnvFileCustom     string = "./.env"
	DefaultEnvFileDevelop    string = "./.env.develop"
	DefaultEnvFileStaging    string = "./.env.staging"
	DefaultEnvFileProduction string = "./.env.production"
)

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

type EnvLoader struct {
	OccurredErrors     yeterr.Collection
	CurrentEnvironment Environment
	Variables          EnvVariables
	EnvFileCustom      string
	EnvFileDevelop     string
	EnvFileStaging     string
	EnvFileProduction  string
}

func NewEnvLoader() EnvLoader {
	return EnvLoader{
		OccurredErrors:     yeterr.NewErrorCollection(),
		CurrentEnvironment: GetEnvironment(),
		Variables:          EnvVariables{},
		EnvFileCustom:      DefaultEnvFileCustom,
		EnvFileDevelop:     DefaultEnvFileDevelop,
		EnvFileStaging:     DefaultEnvFileStaging,
		EnvFileProduction:  DefaultEnvFileProduction,
	}
}
