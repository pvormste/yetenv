package yetenv

import (
	"io/ioutil"
	"regexp"

	"github.com/pvormste/yeterr"
)

var treatEnvFileNotFoundAsFatalError = false

const (
	flagEnvFileNotFound yeterr.ErrorFlag = "envFileNotFound"

	metadataFilePathKey string = "file_path"

	dotenvLineRegex string = `^(\s)*(export)?(\s)*([a-zA-Z_0-9])+=(")?(.)*(")?$`
)

type dotenvFileParser struct {
	occurredErrors yeterr.Collection
	lineRegEx      *regexp.Regexp
}

func newDotenvFileParser() dotenvFileParser {
	return dotenvFileParser{
		occurredErrors: yeterr.NewErrorCollection(),
		lineRegEx:      regexp.MustCompile(dotenvLineRegex),
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

func (p *dotenvFileParser) isLineValid(line string) bool {
	return p.lineRegEx.MatchString(line)
}

func (p *dotenvFileParser) sanitizeLine(line string) string {
	return ""
}

func (p *dotenvFileParser) parseSanitizedLine(sanitizedLine string) (variable string, value string) {
	return "", ""
}

func (p *dotenvFileParser) parseFromBytes(content []byte) {

}
