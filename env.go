package yetenv

import (
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"path/filepath"
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

var DefaultConditionForDevelopEnvironment = func(configLoader *ConfigLoader, currentEnvironment Environment) bool {
	return currentEnvironment == Develop
}

var DefaultConditionForStagingEnvironment = func(configLoader *ConfigLoader, currentEnvironment Environment) bool {
	return currentEnvironment == Staging
}

var DefaultConditionForProductionEnvironment = func(configLoader *ConfigLoader, currentEnvironment Environment) bool {
	return currentEnvironment == Production
}

type LoadBehavior int

const (
	LoadBehaviorUnknown LoadBehavior = iota
	LoadBehaviorDefault
	LoadBehaviorCustom
)

// DefaultVariableName defines the default name of the environment variable.
var DefaultVariableName = "ENVIRONMENT"

var (
	ErrUnknownLoadBehavior = errors.New("load behavior is unknown - only default or custom load behavior is allowed")
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

type ConfigLoader struct {
	LoadPath      string
	FileExtension ConfigFileExtension
	ConfigFiles   map[Environment]string
	Environment   Environment
	LoadBehavior  LoadBehavior
	loadOrder     []loadOrderItem
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
		Environment:  "",
		LoadBehavior: LoadBehaviorUnknown,
		loadOrder:    []loadOrderItem{},
	}
}

func (c *ConfigLoader) UseLoadPath(path string) *ConfigLoader {
	c.LoadPath = path
	return c
}

func (c *ConfigLoader) UseFileProcessor(extension ConfigFileExtension) *ConfigLoader {
	c.FileExtension = extension
	return c
}

func (c *ConfigLoader) UseFileNameForEnvironment(environment Environment, fileName string) *ConfigLoader {
	fileExtensions := []ConfigFileExtension{DOTENV, YAML, TOML, JSON}
	for _, fileExtension := range fileExtensions {
		if strings.HasSuffix(fileName, string(fileExtension)) {
			fileName = strings.TrimRight(fileName, string(fileExtension))
			break
		}
	}

	c.ConfigFiles[environment] = fileName

	return c
}

func (c *ConfigLoader) UseEnvironment(environment Environment) *ConfigLoader {
	c.Environment = environment
	return c
}

func (c *ConfigLoader) UseLoadBehavior(behavior LoadBehavior) *ConfigLoader {
	c.LoadBehavior = behavior
	return c
}

func (c *ConfigLoader) UseDefaultLoadBehavior() *ConfigLoader {
	c.UseLoadBehavior(LoadBehaviorDefault)
	return c
}

func (c *ConfigLoader) UseCustomLoadBehavior() *ConfigLoader {
	c.UseLoadBehavior(LoadBehaviorCustom)
	return c
}

func (c *ConfigLoader) LoadFromFileForEnvironment(environment Environment) *ConfigLoader {
	configFileName := c.ConfigFiles[environment]
	if c.FileExtension == DOTENV && configFileName == defaultCustomConfigFile {
		configFileName = ""
	}

	fullFilePath := c.composeFilePath(c.LoadPath, configFileName, c.FileExtension)

	switch environment {
	case Develop:
		c.LoadFromConditionalFile(fullFilePath, DefaultConditionForDevelopEnvironment)
	case Staging:
		c.LoadFromConditionalFile(fullFilePath, DefaultConditionForStagingEnvironment)
	case Production:
		c.LoadFromConditionalFile(fullFilePath, DefaultConditionForProductionEnvironment)
	case Custom:
		c.LoadFromFile(fullFilePath)
	}

	return c
}

func (c *ConfigLoader) LoadFromFile(filePath string) *ConfigLoader {
	c.loadOrder = append(c.loadOrder, loadOrderItem{
		file:          filePath,
		conditionFunc: nil,
	})

	return c
}

func (c *ConfigLoader) LoadFromConditionalFile(filePath string, conditionFunc ConditionalLoadFunc) *ConfigLoader {
	c.loadOrder = append(c.loadOrder, loadOrderItem{
		file:          filePath,
		conditionFunc: conditionFunc,
	})

	return c
}

func (c *ConfigLoader) LoadInto(cfg interface{}) error {
	switch c.LoadBehavior {
	case LoadBehaviorUnknown:
		return ErrUnknownLoadBehavior
	case LoadBehaviorDefault:
		c.setupDefaultLoadBehavior()
	}

	if c.Environment == "" {
		c.Environment = GetEnvironment()
	}

	for _, loadItem := range c.loadOrder {
		canLoadFile := true
		if loadItem.conditionFunc != nil {
			canLoadFile = loadItem.conditionFunc(c, c.Environment)
		}

		if !canLoadFile {
			continue
		}

		err := c.loadConfigFromFile(loadItem.file, cfg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ConfigLoader) loadConfigFromFile(file string, cfg interface{}) error {
	if fileExists(file) {
		return cleanenv.ReadConfig(file, cfg)
	}

	return nil
}

func (c *ConfigLoader) composeFilePath(loadPath string, fileName string, fileExtension ConfigFileExtension) string {
	fileName = strings.TrimRight(fileName, ".")
	fullFileName := fmt.Sprintf("%s%s", fileName, string(fileExtension))

	return filepath.Join(loadPath, fullFileName)
}

func (c *ConfigLoader) setupDefaultLoadBehavior() {
	c.loadOrder = []loadOrderItem{}

	c.LoadFromFileForEnvironment(Develop)
	c.LoadFromFileForEnvironment(Staging)
	c.LoadFromFileForEnvironment(Production)
	c.LoadFromFileForEnvironment(Custom)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

type loadOrderItem struct {
	file          string
	conditionFunc ConditionalLoadFunc
}
