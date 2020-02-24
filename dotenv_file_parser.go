package yetenv

import (
	"io/ioutil"

	"github.com/pvormste/yeterr"
)

var treatEnvFileNotFoundAsFatalError = false

const (
	flagEnvFileNotFound yeterr.ErrorFlag = "envFileNotFound"

	metadataFilePathKey string = "file_path"
)

type dotenvFileParser struct {
	occurredErrors yeterr.Collection
}

func newDotenvFileParser() dotenvFileParser {
	return dotenvFileParser{
		occurredErrors: yeterr.NewErrorCollection(),
	}
}

func (p *dotenvFileParser) readBytesFromFile(pathToFile string) (content []byte, ok bool) {
	content, err := ioutil.ReadFile(pathToFile)
	if err == nil {
		return content, true
	}

	errMetadata := yeterr.ErrorMetadata{
		metadataFilePathKey: pathToFile,
	}

	if treatEnvFileNotFoundAsFatalError {
		p.occurredErrors.AddFlaggedFatalError(err, errMetadata, flagEnvFileNotFound)
		return nil, false
	}

	p.occurredErrors.AddFlaggedError(err, errMetadata, flagEnvFileNotFound)
	return nil, false
}

func (p *dotenvFileParser) parseFromBytes(content []byte) {

}
