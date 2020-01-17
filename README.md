[![GoDoc](https://godoc.org/github.com/pvormste/yetenv?status.svg)](https://godoc.org/github.com/pvormste/yetenv) ![](https://github.com/pvormste/yetenv/workflows/lint/badge.svg?branch=master) ![](https://github.com/pvormste/yetenv/workflows/tests/badge.svg?branch=master)

# yetenv

`yetenv` is small util package which helps to determine on which environment the application is running. It reads the `ENVIRONMENT` variable.

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

```bash
$ ENVIRONMENT="production" go run main.go   # yetenv.Production
$ ENVIRONMENT="staging" go run main.go      # yetenv.Staging
$ ENVIRONMENT="test" go run main.go         # yetenv.Develop
```

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
