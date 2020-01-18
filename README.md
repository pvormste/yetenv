[![GoDoc](https://godoc.org/github.com/pvormste/yetenv?status.svg)](https://godoc.org/github.com/pvormste/yetenv) ![](https://github.com/pvormste/yetenv/workflows/lint/badge.svg?branch=master) ![](https://github.com/pvormste/yetenv/workflows/tests/badge.svg?branch=master)

# yetenv

`yetenv` is small util package which helps to determine on which environment the application is running. 
It reads the a specific environment variable (Default: `ENVIRONMENT`). See examples below.

The `ENVIRONMENT` values are not case-sensitives.

| `ENVIRONMENT` value | Constant |
| ------------------- | -------- |
| `production` | yetenv.Production |
| `staging` | yetenv.Staging |
| any other value | yetenv.Develop |

## Install

```bash
go get -u github.com/pvormste/yetenv
```

## Usage

### Using the defaults

Shell:
```bash
$ ENVIRONMENT="production" go run main.go   # yetenv.Production
$ ENVIRONMENT="staging" go run main.go      # yetenv.Staging
$ ENVIRONMENT="test" go run main.go         # yetenv.Develop
```

Go Code:
```go
environment := yetenv.GetEnvironment()

switch environment {
case yetenv.Production:
    // Do something in production environment
case yetenv.Staging:
    // Do something in staging environment
case yetenv.Develop:
    // Do something in develop environment
}
```

### Changing the default variable name

Shell:
```bash
$ APP_ENV="production" go run main.go   # yetenv.Production
$ APP_ENV="staging" go run main.go      # yetenv.Staging
$ APP_ENV="test" go run main.go         # yetenv.Develop
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
case yetenv.Develop:
    // Do something in develop environment
}
```

### Reading temporary from a custom variable name

Shell:
```bash
$ APP_ENV="production" go run main.go   # yetenv.Production
$ APP_ENV="staging" go run main.go      # yetenv.Staging
$ APP_ENV="test" go run main.go         # yetenv.Develop
```

Go Code:
```go
environment := yetenv.GetEnvironmentFromVariable("APP_ENV")

switch environment {
case yetenv.Production:
    // Do something in production environment
case yetenv.Staging:
    // Do something in staging environment
case yetenv.Develop:
    // Do something in develop environment
}
```