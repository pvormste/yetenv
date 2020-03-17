package injector

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pvormste/yetenv/dotenv"
)

type testStruct struct {
	BoolValue     bool
	Float32Value  float32
	Float64Value  float64
	IntValue      int
	Int8Value     int8
	Int16Value    int16
	Int32Value    int32
	Int64Value    int64
	StringValue   string
	UIntValue     uint
	UInt8Value    uint8
	UInt16Value   uint16
	UInt32Value   uint32
	UInt64Value   uint64
	UnhandledType *bool
}

func (ts *testStruct) fieldValueByName(fieldName string) reflect.Value {
	testStructPtr := reflect.ValueOf(ts)
	testStructValue := testStructPtr.Elem()
	return testStructValue.FieldByName(fieldName)
}

func TestEnvInjector_setStructFieldValue(t *testing.T) {
	t.Run("bool", func(t *testing.T) {
		t.Run("should set successfully a bool value", func(t *testing.T) {
			instance := &testStruct{}
			field := instance.fieldValueByName("BoolValue")

			injector := NewEnvInjector()
			success := injector.setStructFieldValue(field, "true", "BoolValue", "testStruct")

			assert.True(t, success)
			assert.Equal(t, true, instance.BoolValue)
		})
	})

	t.Run("float", func(t *testing.T) {
		t.Run("should set successfully a float32 value", func(t *testing.T) {
			instance := &testStruct{}
			field := instance.fieldValueByName("Float32Value")

			injector := NewEnvInjector()
			success := injector.setStructFieldValue(field, "3.3", "Float32Value", "testStruct")

			assert.True(t, success)
			assert.Equal(t, float32(3.3), instance.Float32Value)
		})

		t.Run("should set successfully a float64 value", func(t *testing.T) {
			instance := &testStruct{}
			field := instance.fieldValueByName("Float64Value")

			injector := NewEnvInjector()
			success := injector.setStructFieldValue(field, "3.3", "Float64Value", "testStruct")

			assert.True(t, success)
			assert.Equal(t, float64(3.3), instance.Float64Value)
		})
	})

	t.Run("int", func(t *testing.T) {
		t.Run("should set successfully a int value", func(t *testing.T) {
			instance := &testStruct{}
			field := instance.fieldValueByName("IntValue")

			injector := NewEnvInjector()
			success := injector.setStructFieldValue(field, "1", "IntValue", "testStruct")

			assert.True(t, success)
			assert.Equal(t, int(1), instance.IntValue)
		})

		t.Run("should set successfully a int8 value", func(t *testing.T) {
			instance := &testStruct{}
			field := instance.fieldValueByName("Int8Value")

			injector := NewEnvInjector()
			success := injector.setStructFieldValue(field, "1", "Int8Value", "testStruct")

			assert.True(t, success)
			assert.Equal(t, int8(1), instance.Int8Value)
		})

		t.Run("should set successfully a int16 value", func(t *testing.T) {
			instance := &testStruct{}
			field := instance.fieldValueByName("Int16Value")

			injector := NewEnvInjector()
			success := injector.setStructFieldValue(field, "1", "Int16Value", "testStruct")

			assert.True(t, success)
			assert.Equal(t, int16(1), instance.Int16Value)
		})

		t.Run("should set successfully a int32 value", func(t *testing.T) {
			instance := &testStruct{}
			field := instance.fieldValueByName("Int32Value")

			injector := NewEnvInjector()
			success := injector.setStructFieldValue(field, "1", "Int32Value", "testStruct")

			assert.True(t, success)
			assert.Equal(t, int32(1), instance.Int32Value)
		})

		t.Run("should set successfully a int64 value", func(t *testing.T) {
			instance := &testStruct{}
			field := instance.fieldValueByName("Int64Value")

			injector := NewEnvInjector()
			success := injector.setStructFieldValue(field, "1", "Int64Value", "testStruct")

			assert.True(t, success)
			assert.Equal(t, int64(1), instance.Int64Value)
		})
	})

	t.Run("string", func(t *testing.T) {
		t.Run("should set successfully a string value", func(t *testing.T) {
			instance := &testStruct{}
			field := instance.fieldValueByName("StringValue")

			injector := NewEnvInjector()
			success := injector.setStructFieldValue(field, "myString", "StringValue", "testStruct")

			assert.True(t, success)
			assert.Equal(t, "myString", instance.StringValue)
		})
	})

	t.Run("uint", func(t *testing.T) {
		t.Run("should set successfully a uint value", func(t *testing.T) {
			instance := &testStruct{}
			field := instance.fieldValueByName("UIntValue")

			injector := NewEnvInjector()
			success := injector.setStructFieldValue(field, "1", "UIntValue", "testStruct")

			assert.True(t, success)
			assert.Equal(t, uint(1), instance.UIntValue)
		})

		t.Run("should set successfully a uint8 value", func(t *testing.T) {
			instance := &testStruct{}
			field := instance.fieldValueByName("UInt8Value")

			injector := NewEnvInjector()
			success := injector.setStructFieldValue(field, "1", "UInt8Value", "testStruct")

			assert.True(t, success)
			assert.Equal(t, uint8(1), instance.UInt8Value)
		})

		t.Run("should set successfully a uint16 value", func(t *testing.T) {
			instance := &testStruct{}
			field := instance.fieldValueByName("UInt16Value")

			injector := NewEnvInjector()
			success := injector.setStructFieldValue(field, "1", "UInt16Value", "testStruct")

			assert.True(t, success)
			assert.Equal(t, uint16(1), instance.UInt16Value)
		})

		t.Run("should set successfully a uint32 value", func(t *testing.T) {
			instance := &testStruct{}
			field := instance.fieldValueByName("UInt32Value")

			injector := NewEnvInjector()
			success := injector.setStructFieldValue(field, "1", "UInt32Value", "testStruct")

			assert.True(t, success)
			assert.Equal(t, uint32(1), instance.UInt32Value)
		})

		t.Run("should set successfully a uint64 value", func(t *testing.T) {
			instance := &testStruct{}
			field := instance.fieldValueByName("UInt64Value")

			injector := NewEnvInjector()
			success := injector.setStructFieldValue(field, "1", "UInt64Value", "testStruct")

			assert.True(t, success)
			assert.Equal(t, uint64(1), instance.UInt64Value)
		})
	})

	t.Run("error", func(t *testing.T) {
		t.Run("should file an error on unhandled type", func(t *testing.T) {
			instance := &testStruct{}
			field := instance.fieldValueByName("UnhandledType")

			injector := NewEnvInjector()
			success := injector.setStructFieldValue(field, "true", "UnhandledType", "testStruct")

			assert.False(t, success)
			assert.Equal(t, ErrorFlagUnhandledType, injector.OccurredErrors.FirstError().Flag)
		})

		t.Run("should file an error on failed type conversion", func(t *testing.T) {
			instance := &testStruct{}
			field := instance.fieldValueByName("BoolValue")

			injector := NewEnvInjector()
			success := injector.setStructFieldValue(field, "NotABool", "BoolValue", "testStruct")

			assert.False(t, success)
			assert.Equal(t, ErrorFlagFailedTypeParsing, injector.OccurredErrors.FirstError().Flag)
		})
	})

}

func TestEnvInjector_getEnvValue(t *testing.T) {
	t.Run("when struct tag is NOT set", func(t *testing.T) {
		t.Run("should return empty string when value for field cant be found", func(t *testing.T) {
			injector := NewEnvInjector()
			envValue := injector.getEnvValue(dotenv.Variables{"NO_FIELD": "value"}, "field", "")
			assert.Equal(t, "", envValue)
		})

		t.Run("should load from variables map by uppercase field name", func(t *testing.T) {
			injector := NewEnvInjector()
			envValue := injector.getEnvValue(dotenv.Variables{"FIELD": "value"}, "field", "")
			assert.Equal(t, "value", envValue)
		})

		t.Run("should load from OS by uppercase field name", func(t *testing.T) {
			err := os.Setenv("FIELD", "value")
			require.NoError(t, err)

			injector := NewEnvInjector()
			envValue := injector.getEnvValue(dotenv.Variables{"FIELD": "unused value"}, "field", "")
			assert.Equal(t, "value", envValue)
		})
	})

	t.Run("when struct tag is set", func(t *testing.T) {
		t.Run("should return empty string when value for field cant be found", func(t *testing.T) {
			injector := NewEnvInjector()
			envValue := injector.getEnvValue(dotenv.Variables{"NO_FIELD": "value"}, "field", "ENV_NAME")
			assert.Equal(t, "", envValue)
		})

		t.Run("should load from variables map by struct tag", func(t *testing.T) {
			injector := NewEnvInjector()
			envValue := injector.getEnvValue(dotenv.Variables{"ENV_NAME": "value"}, "field", "ENV_NAME")
			assert.Equal(t, "value", envValue)
		})

		t.Run("should load from OS by struct tag", func(t *testing.T) {
			err := os.Setenv("ENV_NAME", "value")
			require.NoError(t, err)

			injector := NewEnvInjector()
			envValue := injector.getEnvValue(dotenv.Variables{"ENV_NAME": "unused value"}, "field", "ENV_NAME")
			assert.Equal(t, "value", envValue)
		})
	})
}
