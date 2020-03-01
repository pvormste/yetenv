package yetenv

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/pvormste/yeterr"
)

var treatEnvFileNotFoundAsFatalError = false

const (
	flagEnvFileNotFound yeterr.ErrorFlag = "envFileNotFound"
	flagInvalidLine     yeterr.ErrorFlag = "invalidLine"

	metadataFilePathKey string = "file_path"

	dotenvLineRegex   string = `^(\s)*(export)?(\s)*([a-zA-Z_0-9])+=(")?(.)*(")?$`
	dotenvExportRegex string = `^(\s)*(export)?(\s)*`
)

var errInvalidDotenvLine = errors.New("line is not a valid dotenv assignent")

type dotenvFileParser struct {
	occurredErrors yeterr.Collection
	lineRegEx      *regexp.Regexp
	exportRegEx    *regexp.Regexp
}

func newDotenvFileParser() dotenvFileParser {
	return dotenvFileParser{
		occurredErrors: yeterr.NewErrorCollection(),
		lineRegEx:      regexp.MustCompile(dotenvLineRegex),
		exportRegEx:    regexp.MustCompile(dotenvExportRegex),
	}
}

func (p *dotenvFileParser) parse(pathToFile string) (variables dotenvVariables, ok bool) {
	fileContent, ok := p.readBytesFromFile(pathToFile)
	if !ok {
		return nil, false
	}

	variables = p.parseFromBytes(fileContent)
	return variables, true
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
	if p.exportRegEx.MatchString(line) {
		splittedLine := p.exportRegEx.Split(line, 2)
		if len(splittedLine) > 1 {
			line = splittedLine[1]
		}
	}

	return strings.Trim(line, " ")
}

func (p *dotenvFileParser) parseSanitizedLine(sanitizedLine string) (variable string, value string) {
	splittedLine := strings.SplitN(sanitizedLine, "=", 2)
	variable = splittedLine[0]

	if len(splittedLine) > 1 {
		value = strings.Trim(splittedLine[1], `"`)
		return variable, value
	}

	return variable, ""
}

func (p *dotenvFileParser) parseFromBytes(content []byte) (variables dotenvVariables) {
	bytesBuffer := bytes.NewBuffer(content)
	scanner := bufio.NewScanner(bytesBuffer)
	variables = make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()

		if !p.isLineValid(line) {
			errMetadata := yeterr.ErrorMetadata{"line": line}
			p.occurredErrors.AddFlaggedError(errInvalidDotenvLine, errMetadata, flagInvalidLine)
			continue
		}

		sanitizedLine := p.sanitizeLine(line)
		variable, value := p.parseSanitizedLine(sanitizedLine)
		variables[variable] = value
	}

	return variables
}
