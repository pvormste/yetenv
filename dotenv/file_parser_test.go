package dotenv

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	pathToNonExistingEnvFile = "./testdata/.env.not-existing"
	pathToEnvMinimal         = "./testdata/.env.minimal"
	pathToEnvValid           = "./testdata/.env.valid"
)

func TestDotenvFileParser_parse(t *testing.T) {
	t.Run("should return not ok when file cant be read", func(t *testing.T) {
		parser := NewFileParser()
		vars, ok := parser.Parse(pathToNonExistingEnvFile)

		assert.False(t, ok)
		assert.Nil(t, vars)
	})

	t.Run("should successfully parse a dotenv file", func(t *testing.T) {
		parser := NewFileParser()
		vars, ok := parser.Parse(pathToEnvValid)

		assert.True(t, ok)
		assert.Equal(t, "world", vars["HELLO"])
		assert.Equal(t, "value", vars["MY_VAR"])
		assert.Equal(t, "your value", vars["YOUR_VAR"])
	})
}

func TestDotenvFileParser_readBytesFromFile(t *testing.T) {
	t.Run("should add an error to occurredErrors when .env file cant be found but not a fatalError", func(t *testing.T) {
		treatEnvFileNotFoundAsFatalError = false

		parser := NewFileParser()
		content, ok := parser.readBytesFromFile(pathToNonExistingEnvFile)

		assert.Nil(t, content)
		assert.False(t, ok)
		assert.Equal(t, 1, parser.occurredErrors.Count())
		assert.False(t, parser.occurredErrors.HasFatalError())
		assert.Equal(t, pathToNonExistingEnvFile, parser.occurredErrors.FirstError().Metadata[metadataFilePathKey])
		assert.Equal(t, flagEnvFileNotFound, parser.occurredErrors.FirstError().Flag)
	})

	t.Run("should add an fatal error to occurredErrors when .env file cant be found", func(t *testing.T) {
		treatEnvFileNotFoundAsFatalError = true

		parser := NewFileParser()
		content, ok := parser.readBytesFromFile(pathToNonExistingEnvFile)

		assert.Nil(t, content)
		assert.False(t, ok)
		assert.Equal(t, 1, parser.occurredErrors.Count())
		assert.True(t, parser.occurredErrors.HasFatalError())
		assert.Equal(t, pathToNonExistingEnvFile, parser.occurredErrors.FirstError().Metadata[metadataFilePathKey])
		assert.Equal(t, flagEnvFileNotFound, parser.occurredErrors.FirstError().Flag)
	})

	t.Run("should read content from .env file successfully", func(t *testing.T) {
		parser := NewFileParser()
		content, ok := parser.readBytesFromFile(pathToEnvMinimal)

		assert.Equal(t, []byte("HELLO=WORLD"), content)
		assert.True(t, ok)
		assert.Equal(t, 0, parser.occurredErrors.Count())
	})
}

func TestIsLineValid(t *testing.T) {
	t.Run("should return true for valid cases", func(t *testing.T) {
		lineCases := []string{
			`VARIABLE=value`,
			`VARIABLE="many values"`,
			`VARIABLE="many_values"`,
			`VARIABLE007="value"`,
			`VARIABLE=many values`,
			`VAR_IABLE=value`,
			`VAR_IAB_LE="value"`,
			`VARIABLE="value"       `,
			`      VARIABLE="value"`,
			`		VARIABLE=value`,
			`export VARIABLE="value"`,
			`    export VARIABLE="value"`,
		}

		for _, lineCase := range lineCases {
			t.Run(lineCase, func(t *testing.T) {
				parser := NewFileParser()
				isValid := parser.isLineValid(lineCase)

				assert.True(t, isValid)
			})
		}
	})

	t.Run("should return false for invalid cases", func(t *testing.T) {
		lineCases := []string{
			`VARIABLE value`,
			`VAR_IABLE value`,
			`=value`,
			`="value"`,
		}

		for _, lineCase := range lineCases {
			t.Run(lineCase, func(t *testing.T) {
				parser := NewFileParser()
				isValid := parser.isLineValid(lineCase)

				assert.False(t, isValid)
			})
		}
	})
}

func TestSanitizeLine(t *testing.T) {
	linesToSanitize := []string{
		`   VAR_IAB_LE="value"   `,
		`   export VAR_IAB_LE="value"   `,
	}

	expectedSanitizedLines := []string{
		`VAR_IAB_LE="value"`,
		`VAR_IAB_LE="value"`,
	}

	for i := 0; i < len(linesToSanitize); i++ {
		t.Run(fmt.Sprintf("should successfully sanitize line: %s", linesToSanitize[i]), func(t *testing.T) {
			parser := NewFileParser()
			sanitizedLine := parser.sanitizeLine(linesToSanitize[i])

			assert.Equal(t, expectedSanitizedLines[i], sanitizedLine)
		})
	}
}

func TestParseSanitizedLine(t *testing.T) {
	t.Run("should parse successfully sanitized line:", func(t *testing.T) {
		sanitizedLines := []string{
			`VAR_IAB_LE="value"`,
			`VARIABLE=value`,
			`VARIABLE10="value=123"`,
		}

		expectedVariable := []string{
			"VAR_IAB_LE",
			"VARIABLE",
			"VARIABLE10",
		}

		expectedValue := []string{
			"value",
			"value",
			"value=123",
		}

		for i, sanitizedLine := range sanitizedLines {
			t.Run(sanitizedLine, func(t *testing.T) {
				parser := NewFileParser()
				variable, value := parser.parseSanitizedLine(sanitizedLine)

				assert.Equal(t, expectedVariable[i], variable)
				assert.Equal(t, expectedValue[i], value)
			})
		}
	})
}

func TestParseFromBytes(t *testing.T) {
	t.Run("should return empty Variables when content is empty", func(t *testing.T) {
		content := []byte("")

		parser := NewFileParser()
		variables := parser.parseFromBytes(content)

		assert.Equal(t, 0, variables.Count())
		assert.Equal(t, 0, parser.occurredErrors.Count())
	})

	t.Run("should return empty Variables when content does not contain valid variables", func(t *testing.T) {
		content := []byte("VARIABLE1 value1\nVARIABLE2 value2")

		parser := NewFileParser()
		variables := parser.parseFromBytes(content)

		assert.Equal(t, 0, variables.Count())
		assert.Equal(t, 2, parser.occurredErrors.Count())
	})

	t.Run("should successfully parse dotenv variables and values", func(t *testing.T) {
		content := []byte("VARIABLE1=value1\nexport VARIABLE2=value2")

		parser := NewFileParser()
		variables := parser.parseFromBytes(content)

		assert.Equal(t, Variables{"VARIABLE1": "value1", "VARIABLE2": "value2"}, variables)
		assert.Equal(t, 0, parser.occurredErrors.Count())
	})
}
