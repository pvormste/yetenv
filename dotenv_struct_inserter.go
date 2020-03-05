package yetenv

import (
	"os"
	"reflect"
	"strings"

	"github.com/pvormste/yeterr"
)

const (
	envStructField = "env"
)

type structEnvInserter struct {
	occurredErrors yeterr.Collection
}

func newStructEnvInserter() structEnvInserter {
	return structEnvInserter{
		occurredErrors: yeterr.NewErrorCollection(),
	}
}

func (sp *structEnvInserter) insertVariables(input interface{}, variables EnvVariables) {
	inputType := reflect.TypeOf(input)
	kind := inputType.Kind()

	switch kind {
	case reflect.Struct:
		for i := 0; i < inputType.NumField(); i++ {
			/*field := inputType.Field(i)
			fieldName := field.Name
			fieldType := field.Type
			envName := field.Tag.Get(envStructField)*/

		}
	default:
	}
}

func (sp *structEnvInserter) getEnvValue(variables EnvVariables, fieldName string, structTagValue string) string {
	actualEnvName := structTagValue
	if len(actualEnvName) == 0 {
		actualEnvName = strings.ToUpper(fieldName)
	}

	actualEnvValue := variables[actualEnvName]
	envValueFromOS := os.Getenv(actualEnvName)
	if len(envValueFromOS) > 0 {
		actualEnvValue = envValueFromOS
	}

	return actualEnvValue
}
