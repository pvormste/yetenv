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

	ErrorMetadataKeyKind       = "kind"
	ErrorMetadataKeyEnvValue   = "envValue"
	ErrorMetadataKeyFieldName  = "field_name"
	ErrorMetadataKeyStructName = "struct_name"

	prefixExtensionTemplate     = "%s_"
	prefixWithFieldNameTemplate = "%s%s"

	structTagEnv = "env"
)

var errRequiredPointer = errors.New("struct should be a pointer (*struct)")

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
	targetPtr := reflect.ValueOf(target)
	if targetPtr.Kind() != reflect.Ptr {
		return errRequiredPointer
	}

	var err error

	switch targetPtr.Elem().Kind() {
	case reflect.Struct:
		err = e.injectIntoStruct(targetPtr, variables, "")
	default:
	}

	if err != nil {
		return err
	}

	return nil
}

func (e *EnvInjector) injectIntoStruct(structPtr reflect.Value, variables dotenv.Variables, prefix string) error {
	structValue := structPtr.Elem()
	structName := structValue.Type().Name()
	for i := 0; i < structValue.Type().NumField(); i++ {
		field := structValue.Type().Field(i)
		fieldPtr := structValue.Field(i).Addr()

		var err error

		switch field.Type.Kind() {
		case reflect.Struct:
			newPrefix := e.extendPrefix(prefix, field.Name)
			err = e.injectIntoStruct(fieldPtr, variables, newPrefix)
		default:
			prefixedFieldName := e.appendFieldNameToPrefix(prefix, field.Name)
			envValue := e.getEnvValue(variables, prefixedFieldName, field.Tag.Get(structTagEnv))

			e.setStructFieldValue(fieldPtr, envValue, prefixedFieldName, structName)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (e *EnvInjector) setStructFieldValue(field reflect.Value, envValue string, fieldName string, structName string) bool {
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
			ErrorMetadataKeyStructName: structName,
			ErrorMetadataKeyFieldName:  fieldName,
			ErrorMetadataKeyKind:       fieldValue.Type().Kind().String(),
			ErrorMetadataKeyEnvValue:   envValue,
		}

		if errFlag == "" {
			errFlag = ErrorFlagFailedTypeParsing
		}

		e.OccurredErrors.AddFlaggedError(err, errMetadata, errFlag)

		return false
	}

	return true
}

func (e *EnvInjector) getEnvValue(variables dotenv.Variables, fieldNameWithPrefix string, structTagValue string) string {
	actualEnvName := structTagValue
	if len(actualEnvName) == 0 {
		actualEnvName = strings.ToUpper(fieldNameWithPrefix)
	}

	actualEnvValue := variables[actualEnvName]
	envValueFromOS := os.Getenv(actualEnvName)
	if len(envValueFromOS) > 0 {
		actualEnvValue = envValueFromOS
	}

	return actualEnvValue
}

func (e *EnvInjector) appendFieldNameToPrefix(prefix string, fieldName string) string {
	return fmt.Sprintf(prefixWithFieldNameTemplate, prefix, strings.ToUpper(fieldName))
}

func (e *EnvInjector) extendPrefix(currentPrefix string, fieldName string) string {
	prefixedFieldName := e.appendFieldNameToPrefix(currentPrefix, fieldName)
	return fmt.Sprintf(prefixExtensionTemplate, prefixedFieldName)
}
