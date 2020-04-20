[![GoDoc](https://godoc.org/github.com/pvormste/yetenv?status.svg)](https://godoc.org/github.com/pvormste/yetenv) ![](https://github.com/pvormste/yetenv/workflows/lint/badge.svg?branch=master) ![](https://github.com/pvormste/yetenv/workflows/tests/badge.svg?branch=master)

# yetenv

`yetenv` is small util package which helps to detect on which environment the application is running and can
load a configuration from a file into a configuration struct. 

## Install

```bash
go get -u github.com/pvormste/yetenv
```

## Usage

1. [Environment detection](https://github.com/pvormste/yetenv#environment-detection)
2. [Config Loader](https://github.com/pvormste/yetenv#config-loader)

### Environment detection

Environment detection reads from a specific environment variable (Default: `ENVIRONMENT`). See examples below.

The `ENVIRONMENT` values are not case-sensitives.

| `ENVIRONMENT` value | Constant |
| ------------------- | -------- |
| `production` | yetenv.Production |
| `staging` | yetenv.Staging |
| `test` | yetenv.Test |
| any other value | yetenv.Develop |

#### Using the defaults

Shell:
```bash
$ ENVIRONMENT="production" go run main.go  # yetenv.Production
$ ENVIRONMENT="staging" go run main.go     # yetenv.Staging
$ ENVIRONMENT="test" go run main.go        # yetenv.Test
$ ENVIRONMENT="local" go run main.go       # yetenv.Develop
```

Go Code:
```go
environment := yetenv.GetEnvironment()

switch environment {
case yetenv.Production:
    // Do something in production environment
case yetenv.Staging:
    // Do something in staging environment
case yetenv.Test:
    // Do something in test environment
case yetenv.Develop:
    // Do something in develop environment
}
```

#### Changing the default variable name

Shell:
```bash
$ APP_ENV="production" go run main.go   # yetenv.Production
$ APP_ENV="staging" go run main.go      # yetenv.Staging
$ APP_ENV="test" go run main.go         # yetenv.Test
$ APP_ENV="local" go run main.go         # yetenv.Develop
```

Go Code:
```go
yetenv.DefaultVariableName = "APP_ENV"
environment := yetenv.GetEnvironment()

switch environment {
case yetenv.Production:
    // Do something in production environment
case yetenv.Staging:
    // Do something in staging environment
case yetenv.Test:
    // Do something in test environment
case yetenv.Develop:
    // Do something in develop environment
}
```

#### Reading temporary from a custom variable name

Shell:
```bash
$ APP_ENV="production" go run main.go   # yetenv.Production
$ APP_ENV="staging" go run main.go      # yetenv.Staging
$ APP_ENV="test" go run main.go         # yetenv.Test
$ APP_ENV="local" go run main.go        # yetenv.Develop
```

Go Code:
```go
environment := yetenv.GetEnvironmentFromVariable("APP_ENV")

switch environment {
case yetenv.Production:
    // Do something in production environment
case yetenv.Staging:
    // Do something in staging environment
case yetenv.Test:
    // Do something in staging environment
case yetenv.Develop:
    // Do something in develop environment
}
```

### Config Loader
The config loader is able to load a configuration from one or more files into a configuration struct. To achieve this, 
it uses the [cleanenv package](https://github.com/ilyakaznacheev/cleanenv) under the hood - so please have a look into their
documentation to get familar about the usage of configuration structs.

It also comes with a set of defaults, so it can be used with zero configuration out of the box. Nevertheless it is possible to
customize the behavior of the ConfigLoader, so it doesn't get in your way.

#### Default Load Behavior
The default load behavior works like this depending on the environment:
 - For `Develop`: Load from `./cfg.dev.env` and overwrite it by `./.env` and OS environment values.
 - For `Test`: Load from `./cfg.test.env` and overwrite it by `./.env` and OS environment values.
 - For `Staging`: Load from `./cfg.staging.env` and overwrite it by `./.env` and OS environment values.
 - For `Production`: Load from `./cfg.prod.env` and overwrite it by `./.env` and OS environment values.
 
 ```go
c := Config{}
err := yetenv.NewConfigLoader().
    UseDefaultLoadBehavior().
    LoadInto(&c)
```

##### Change load path
 - For `Develop`: Load from `./config/cfg.dev.env` and overwrite it by `./config/.env` and OS environment values.
 - For `Test`: Load from `./config/cfg.test.env` and overwrite it by `./config/.env` and OS environment values.
 - For `Staging`: Load from `./config/cfg.staging.env` and overwrite it by `./config/.env` and OS environment values.
 - For `Production`: Load from `./config/cfg.prod.env` and overwrite it by `./config/.env` and OS environment values.
 
 ```go
c := Config{}
err := yetenv.NewConfigLoader().
    UseLoadPath("./config").
    UseDefaultLoadBehavior().
    LoadInto(&c)
```

##### Change file processor
 - For `Develop`: Load from `./cfg.dev.yaml` and overwrite it by `./cfg.yaml` and OS environment values.
 - For `Test`: Load from `./cfg.test.yaml` and overwrite it by `./cfg.yaml` and OS environment values.
 - For `Staging`: Load from `./cfg.staging.yaml` and overwrite it by `./cfg.yaml` and OS environment values.
 - For `Production`: Load from `./cfg.prod.yaml` and overwrite it by `./cfg.yaml` and OS environment values.
 
 ```go
c := Config{}
err := yetenv.NewConfigLoader().
    UseFileProcessor(yetenv.YAML).
    UseDefaultLoadBehavior().
    LoadInto(&c)
```

##### Change file name for a specific environment
The default load behavior works like this depending on the environment:
 - For `Develop`: Load from `./cfg.local.env` and overwrite it by `./.env` and OS environment values.
 - For `Test`: Load from `./cfg.testing.env` and overwrite it by `./.env` and OS environment values.
 - For `Staging`: Load from `./cfg.qa.env` and overwrite it by `./.env` and OS environment values.
 - For `Production`: Load from `./cfg.production.env` and overwrite it by `./.env` and OS environment values.
 
 ```go
c := Config{}
err := yetenv.NewConfigLoader().
    UseFileNameForEnvironment(yetenv.Develop, "cfg.local").
    UseFileNameForEnvironment(yetenv.Test, "cfg.testing").
    UseFileNameForEnvironment(yetenv.Staging, "cfg.qa").
    UseFileNameForEnvironment(yetenv.Production, "cfg.production").
    UseDefaultLoadBehavior().
    LoadInto(&c)
```

##### Inject environment
By default the ConfigLoader will use `yetenv.GetEnvironment()` to detect the current environment. If you customize the 
environment detection you can inject it this way:

 ```go
customEnvironment := yetenv.Staging
c := Config{}
err := yetenv.NewConfigLoader().
    UseEnvironment(customEnvironment).
    UseDefaultLoadBehavior().
    LoadInto(&c)
```

#### Custom Load Behavior
The custom load behavior is highly customizable - but keep in mind that some settings only applies when using specific
setting methods.

| Setting | Applies When using | 
| ------- | ------------------ |
| UseLoadPath() | `LoadFromFileForEnvironment()` |
| UseFileProcessor() | `LoadFromFileForEnvironment()` |
| UseFileNameForEnvironment() | `LoadFromFileForEnvironment()` |
| UseEnvironment() | `LoadFromFileForEnvironment()` or `LoadFromConditionalFile()` |

##### Example for a custom load behavior
 - For `Develop` and `Staging` and `Test`: Load from `./cfg.base.env` and overwrite by OS environment values.
 - For `Production`: Load from `./config/cfg.prod.yaml` and overwrite by OS environment values.
 - For `All Environments`: Load from `./other-config/cfg.toml` and  overwrite by OS environment values.

 ```go
developAndStagingCondition := func(configLoader *ConfigLoader, currentEnvironment Environment) bool {
    return currentEnvironment == yetenv.Develop || currentEnvironment == yetenv.Staging || currentEnvironment == yetenv.Test
}

customEnvironment := yetenv.Staging
c := Config{}
err := yetenv.NewConfigLoader().
    UseLoadPath("./config").
    UseFileProcessor(yetenv.YAML).
    UseEnvironment(customEnvironment).
    UseCustomLoadBehavior().
    LoadFromConditionalFile("./cfg.base.env", developAndStagingCondition).
    LoadFromFileForEnvironment(yetenv.Production).
    LoadFromFile("./other-config/cfg.toml").
    LoadInto(&c)
```

As you can see: you can do very complicated things with it - I personally would recommend to keep it simple :-).