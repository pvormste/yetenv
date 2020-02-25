package yetenv

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	pathToNonExistingEnvFile = "./testdata/.env.not-existing"
	pathToEnvMinimal         = "./testdata/.env.minimal"
)

func TestDotenvFileParser_readBytesFromFile(t *testing.T) {
	t.Run("should add an error to occurredErrors when .env file cant be found but not a fatalError", func(t *testing.T) {
		treatEnvFileNotFoundAsFatalError = false

		parser := newDotenvFileParser()
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

		parser := newDotenvFileParser()
		content, ok := parser.readBytesFromFile(pathToNonExistingEnvFile)

		assert.Nil(t, content)
		assert.False(t, ok)
		assert.Equal(t, 1, parser.occurredErrors.Count())
		assert.True(t, parser.occurredErrors.HasFatalError())
		assert.Equal(t, pathToNonExistingEnvFile, parser.occurredErrors.FirstError().Metadata[metadataFilePathKey])
		assert.Equal(t, flagEnvFileNotFound, parser.occurredErrors.FirstError().Flag)
	})

	t.Run("should read content from .env file successfully", func(t *testing.T) {
		parser := newDotenvFileParser()
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
				parser := newDotenvFileParser()
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
				parser := newDotenvFileParser()
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
			parser := newDotenvFileParser()
			sanitizedLine := parser.sanitizeLine(linesToSanitize[i])

			assert.Equal(t, expectedSanitizedLines[i], sanitizedLine)
		})
	}
}
