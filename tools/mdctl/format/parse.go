package format

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
)

const (
	CREATE       = "create"
	FROM         = "from"
	NAME         = "name"
	FAMILY       = "family"
	ARCHITECTURE = "architecture"
	LICENSE      = "license"
	DESCRIPTION  = "description"
	PARAMSIZE    = "param_size"
	WEIGHTS      = "weights"
	TOKENIZER    = "tokenizer"
	PRECISION    = "precision"
	FORMAT       = "format"
	QUANTIZATION = "quantization"
	CONFIG       = "config"
)

type Command struct {
	Name string
	Args string
}

func (c *Command) Reset() {
	c.Name = ""
	c.Args = ""
}

func Parse(reader io.Reader) ([]Command, error) {
	var commands []Command
	var command Command
	var modelCommand Command
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, bufio.MaxScanTokenSize), bufio.MaxScanTokenSize)
	scanner.Split(scanModelfile)
	for scanner.Scan() {
		line := scanner.Bytes()

		fields := bytes.SplitN(line, []byte(" "), 2)
		if len(fields) == 0 || len(fields[0]) == 0 {
			continue
		}

		switch string(bytes.ToUpper(fields[0])) {
		case strings.ToUpper(CREATE), strings.ToUpper(FROM):
			command.Name = string(bytes.ToLower(fields[0]))
			command.Args = string(bytes.TrimSpace(fields[1]))
			modelCommand = command
		case strings.ToUpper(NAME),
			strings.ToUpper(FAMILY),
			strings.ToUpper(ARCHITECTURE),
			strings.ToUpper(LICENSE),
			strings.ToUpper(DESCRIPTION),
			strings.ToUpper(FORMAT),
			strings.ToUpper(PRECISION),
			strings.ToUpper(QUANTIZATION),
			strings.ToUpper(PARAMSIZE),
			strings.ToUpper(WEIGHTS),
			strings.ToUpper(CONFIG),
			strings.ToUpper(TOKENIZER):
			command.Name = string(bytes.ToLower(fields[0]))
			command.Args = string(bytes.TrimSpace(fields[1]))

		default:
			if !bytes.HasPrefix(fields[0], []byte("#")) {
				log.Printf("WARNING: unknown command: %s", fields[0])
			}
			continue
		}

		commands = append(commands, command)
		command.Reset()
	}

	if modelCommand.Args == "" {
		return nil, fmt.Errorf("no FROM or CREATE line was specified")
	}

	return commands, scanner.Err()
}

func scanModelfile(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = scan([]byte(`"""`), []byte(`"""`), data, atEOF)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to scan modelfile: %w", err)
	}

	if advance > 0 && token != nil {
		return advance, token, nil
	}

	advance, token, err = scan([]byte(`"`), []byte(`"`), data, atEOF)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to scan modelfile: %w", err)
	}

	if advance > 0 && token != nil {
		return advance, token, nil
	}

	return bufio.ScanLines(data, atEOF)
}

func scan(openBytes, closeBytes, data []byte, atEOF bool) (advance int, token []byte, err error) {
	newline := bytes.IndexByte(data, '\n')

	if start := bytes.Index(data, openBytes); start >= 0 && start < newline {
		end := bytes.Index(data[start+len(openBytes):], closeBytes)
		if end < 0 {
			if atEOF {
				return 0, nil, fmt.Errorf("unterminated %s: expecting %s", openBytes, closeBytes)
			}
			return 0, nil, nil
		}

		n := start + len(openBytes) + end + len(closeBytes)

		newData := data[:start]
		newData = append(newData, data[start+len(openBytes):n-len(closeBytes)]...)
		return n, newData, nil
	}

	return 0, nil, nil
}
