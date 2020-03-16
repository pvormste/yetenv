package injector

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/pvormste/yeterr"

	"github.com/pvormste/yetenv/dotenv"
)

const (
	ErrorFlagFailedTypeParsing yeterr.ErrorFlag = "failed_type_parsing"
	ErrorFlagUnhandledType     yeterr.ErrorFlag = "unhandled_type"

	ErrorMetadataKeyKind     = "kind"
	ErrorMetadataKeyEnvValue = "envValue"
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

func (e *EnvInjector) setStructFieldValue(field reflect.Value, envValue string) bool {
	if len(envValue) == 0 {
		return true
	}

	var (
		err         error
		errFlag     yeterr.ErrorFlag
		parsedBool  bool
		parsedFloat float64
		parsedInt   int64
		parsedUInt  uint64
	)

	if field.CanAddr() {
		field = field.Addr()
	}

	fieldValue := field.Elem()

	switch fieldValue.Kind() {
	case reflect.Bool:
		parsedBool, err = strconv.ParseBool(envValue)
		if err == nil {
			fieldValue.SetBool(parsedBool)
		}
	case reflect.Float32, reflect.Float64:
		parsedFloat, err = strconv.ParseFloat(envValue, fieldValue.Type().Bits())
		if err == nil {
			fieldValue.SetFloat(parsedFloat)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsedInt, err = strconv.ParseInt(envValue, 10, fieldValue.Type().Bits())
		if err == nil {
			fieldValue.SetInt(parsedInt)
		}
	case reflect.String:
		fieldValue.SetString(envValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsedUInt, err = strconv.ParseUint(envValue, 10, fieldValue.Type().Bits())
		if err == nil {
			fieldValue.SetUint(parsedUInt)
		}
	default:
		errFlag = ErrorFlagUnhandledType
		err = errors.New(fmt.Sprintf("the type of kind '%s' is not handled", fieldValue.Type().Kind().String()))
	}

	if err != nil {
		errMetadata := yeterr.ErrorMetadata{
			ErrorMetadataKeyKind:     fieldValue.Type().Kind().String(),
			ErrorMetadataKeyEnvValue: envValue,
		}

		if errFlag == "" {
			errFlag = ErrorFlagFailedTypeParsing
		}

		e.OccurredErrors.AddFlaggedError(err, errMetadata, errFlag)

		return false
	}

	return true
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
