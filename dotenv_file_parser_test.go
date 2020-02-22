package yetenv

import (
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
