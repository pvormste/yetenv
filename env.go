package yetenv

import (
	"os"
	"strings"
)

const (
	environmentVariableName            = "ENVIRONMENT"
	environmentVariableValueProduction = "production"
	environmentVariableValueStaging    = "staging"
)

type Environment string

const (
	Production Environment = "production"
	Staging    Environment = "staging"
	Develop    Environment = "develop"
)

// GetEnvironment returns the current Environment value depending on the OS environment
// value of the variable ENVIRONMENT.
func GetEnvironment() Environment {
	envRaw := os.Getenv(environmentVariableName)
	env := strings.ToLower(envRaw)

	switch env {
	case environmentVariableValueProduction:
		return Production
	case environmentVariableValueStaging:
		return Staging
	}

	return Develop
}