package yetenv

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	environmentVariableValueProduction = "production"
	environmentVariableValueStaging    = "staging"
	environmentVariableValueTest       = "test"

	defaultDevelopConfigFile    = "cfg.dev"
	defaultTestConfigFile       = "cfg.test"
	defaultStagingConfigFile    = "cfg.staging"
	defaultProductionConfigFile = "cfg.prod"
	defaultCustomConfigFile     = "cfg"
)

// Environment defines the environment of an application (e.g. Develop, Staging, Production, etc.)
type Environment string

const (
	Production Environment = "production"
	Staging    Environment = "staging"
	Test       Environment = "test"
	Develop    Environment = "develop"
	Custom     Environment = "custom"
)

// ConfigFileExtension represents a possible config file extension which is usable by the config loader.
type ConfigFileExtension string

const (
	YAML   ConfigFileExtension = ".yaml"
	JSON   ConfigFileExtension = ".json"
	TOML   ConfigFileExtension = ".toml"
	DOTENV ConfigFileExtension = ".env"
)

// ConditionalLoadFunc allows to define a condition for a file to be loaded by the config loader.
type ConditionalLoadFunc func(configLoader *ConfigLoader, currentEnvironment Environment) bool

// DefaultConditionForDevelopEnvironment retuns true when the current environment is Develop otherwise false.
var DefaultConditionForDevelopEnvironment = func(configLoader *ConfigLoader, currentEnvironment Environment) bool {
	return currentEnvironment == Develop
}

// DefaultConditionForDevelopEnvironment retuns true when the current environment is Develop otherwise false.
var DefaultConditionForTestEnvironment = func(configLoader *ConfigLoader, currentEnvironment Environment) bool {
	return currentEnvironment == Test
}

// DefaultConditionForStagingEnvironment returns true when the current environment is Staging otherwise false.
var DefaultConditionForStagingEnvironment = func(configLoader *ConfigLoader, currentEnvironment Environment) bool {
	return currentEnvironment == Staging
}

// DefaultConditionForProductionEnvironment returns true when the current environment is Production otherweise false.
var DefaultConditionForProductionEnvironment = func(configLoader *ConfigLoader, currentEnvironment Environment) bool {
	return currentEnvironment == Production
}

// LoadBehavior is used to define the load behavior of the config loader.
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
	case environmentVariableValueTest:
		return Test
	}

	return Develop
}

// ConfigLoader loads configuration values from files and the OS environment into a configuration struct.
// It uses the builder pattern and needs to be extecuted by the finishing method.
type ConfigLoader struct {
	LoadPath      string
	FileExtension ConfigFileExtension
	ConfigFiles   map[Environment]string
	Environment   Environment
	LoadBehavior  LoadBehavior
	loadOrder     []loadOrderItem
}

// NewConfigLoader initializes a new ConfigLoader builder.
func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{
		LoadPath:      "./",
		FileExtension: DOTENV,
		ConfigFiles: map[Environment]string{
			Develop:    defaultDevelopConfigFile,
			Test:       defaultTestConfigFile,
			Staging:    defaultStagingConfigFile,
			Production: defaultProductionConfigFile,
			Custom:     defaultCustomConfigFile,
		},
		Environment:  "",
		LoadBehavior: LoadBehaviorUnknown,
		loadOrder:    []loadOrderItem{},
	}
}

// UseLoadPath can be used to change the load path for the default load behavior. It defaults to "./".
func (c *ConfigLoader) UseLoadPath(path string) *ConfigLoader {
	c.LoadPath = path
	return c
}

// UseFileProcessor can be used to change the file processor for config files when using the default
// load behavior. It defaults to "DOTENV".
func (c *ConfigLoader) UseFileProcessor(extension ConfigFileExtension) *ConfigLoader {
	c.FileExtension = extension
	return c
}

// UseFileNameForEnvironment can be used to change the config file name for a specific environment.
// Default file names are:
// Develop     -> 'cfg.dev'
// Develop     -> 'cfg.test'
// Staging     -> 'cfg.staging'
// Production  -> 'cfg.prod'
// Custom      -> 'cfg' or '.env'
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

// UseEnvironment can be used to change the current environment value of the ConfigLoader.
// Defaults to the value of the ENVIRONMENT variable.
func (c *ConfigLoader) UseEnvironment(environment Environment) *ConfigLoader {
	c.Environment = environment
	return c
}

// UseLoadBehavior can be used to set a LoadBehavior to a specific value.
func (c *ConfigLoader) UseLoadBehavior(behavior LoadBehavior) *ConfigLoader {
	c.LoadBehavior = behavior
	return c
}

// UseDefaultLoadBehavior sets load behavior to LoadBehaviorDefault.
func (c *ConfigLoader) UseDefaultLoadBehavior() *ConfigLoader {
	c.UseLoadBehavior(LoadBehaviorDefault)
	return c
}

// UseCustomLoadBehavior sets load behavior to LoadBehaviorCustom.
func (c *ConfigLoader) UseCustomLoadBehavior() *ConfigLoader {
	c.UseLoadBehavior(LoadBehaviorCustom)
	return c
}

// LoadFromFileForEnvironment can be used to reuse environmental load logic for a custom load behavior.
// For example: LoadFromFileForEnvironment(Develop) will behave the same as in the default load behavior.
func (c *ConfigLoader) LoadFromFileForEnvironment(environment Environment) *ConfigLoader {
	configFileName := c.ConfigFiles[environment]
	if c.FileExtension == DOTENV && configFileName == defaultCustomConfigFile {
		configFileName = ""
	}

	fullFilePath := c.composeFilePath(c.LoadPath, configFileName, c.FileExtension)

	switch environment {
	case Develop:
		c.LoadFromConditionalFile(fullFilePath, DefaultConditionForDevelopEnvironment)
	case Test:
		c.LoadFromConditionalFile(fullFilePath, DefaultConditionForTestEnvironment)
	case Staging:
		c.LoadFromConditionalFile(fullFilePath, DefaultConditionForStagingEnvironment)
	case Production:
		c.LoadFromConditionalFile(fullFilePath, DefaultConditionForProductionEnvironment)
	case Custom:
		c.LoadFromFile(fullFilePath)
	}

	return c
}

// LoadFromFile can be used to load a specific config file when using custom load behavior.
// It will not use the LoadPath, so the full path to config file should be provided.
func (c *ConfigLoader) LoadFromFile(filePath string) *ConfigLoader {
	c.loadOrder = append(c.loadOrder, loadOrderItem{
		file:          filePath,
		conditionFunc: nil,
	})

	return c
}

// LoadFromConditionalFile can be used to load a config file only when the condition of the conditionFunc is met.
// It will not use the LoadPath, so the full path to config file should be provided.
func (c *ConfigLoader) LoadFromConditionalFile(filePath string, conditionFunc ConditionalLoadFunc) *ConfigLoader {
	c.loadOrder = append(c.loadOrder, loadOrderItem{
		file:          filePath,
		conditionFunc: conditionFunc,
	})

	return c
}

// LoadInto will finish the ConfigLoader and execute the load process. The provided config struct should be a pointer.
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
	c.LoadFromFileForEnvironment(Test)
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
