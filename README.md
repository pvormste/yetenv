[![GoDoc](https://godoc.org/github.com/pvormste/yetwebutils/yetenv?status.svg)](https://godoc.org/github.com/pvormste/yetwebutils/yetenv)

# yetenv

`yetenv` is small util package which helps to determine on which environment the application is running. It reads the `ENVIRONMENT` variable.

The `ENVIRONMENT` values are not case-sensitives.

| `ENVIRONMENT` value | Constant |
| ------------------- | -------- |
| `production` | yetenv.Production |
| `staging` | yetenv.Staging |
| any other value | yetenv.Develop |

## Usage

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