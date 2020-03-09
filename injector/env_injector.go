package injector

import (
	"errors"
	"os"
	"reflect"
	"strings"

	"github.com/pvormste/yeterr"

	"github.com/pvormste/yetenv/dotenv"
)

var errRequiredPointer = errors.New("struct should be a pointer (*struct)")

type injectionInput struct {
	target          interface{}
	reflectionValue reflect.Value
	reflectionType  reflect.Type
	reflectionKind  reflect.Kind
}

type TypeInjector interface {
	InjectVariables(input interface{}, variables dotenv.Variables)
	OccurredErrors() yeterr.Collection
}

type EnvInjector struct {
	OccurredErrors yeterr.Collection
}

func NewEnvInjector() EnvInjector {
	return EnvInjector{
		OccurredErrors: yeterr.NewErrorCollection(),
	}
}

func (e *EnvInjector) InjectVariables(target interface{}, variables dotenv.Variables) error {
	targetValue := reflect.ValueOf(target)
	targetType := targetValue.Type()
	targetKind := targetValue.Kind()

	input := injectionInput{
		target:          target,
		reflectionValue: targetValue,
		reflectionType:  targetType,
		reflectionKind:  targetKind,
	}

	var err error

	switch targetKind {
	case reflect.Struct:
		err = e.injectIntoStruct(input, variables, "")
	default:
	}

	if err != nil {
		return err
	}

	return nil
}

func (e *EnvInjector) injectIntoStruct(input injectionInput, variables dotenv.Variables, prefix string) error {
	if input.reflectionKind != reflect.Ptr {
		return errRequiredPointer
	}

	dereferencedPtr := input.reflectionValue.Elem()
	for i := 0; i < dereferencedPtr.NumField(); i++ {
		/*field := inputType.Field(i)
		fieldName := field.Name
		fieldType := field.Type
		envName := field.Tag.Get(envStructField)*/
	}

	return nil
}

func (e *EnvInjector) setStructValueByType(structName string, fieldName string, fieldType reflect.Value, envValue string) {

}

func (e *EnvInjector) getEnvValue(variables dotenv.Variables, fieldName string, structTagValue string) string {
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
