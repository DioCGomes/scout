package rubyparser

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/mlw157/scout/internal/models"
)

type RubyParser struct {
}

func NewRubyParser() *RubyParser {
	return &RubyParser{}
}

type FileData struct {
	Path string
	Data []byte
}

func ReadFile(path string) (*FileData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return &FileData{Path: path, Data: data}, nil
}

func (r *RubyParser) ParseFile(path string) ([]models.Dependency, error) {
	fileData, err := ReadFile(path)
	if err != nil {
		return nil, err
	}

	var dependencies []models.Dependency

	// Needs to have exactly 4 spaces in the beginning of the line
	gemLineRegex := regexp.MustCompile(`^    ([a-zA-Z0-9_.-]+)\s+\(([^)]+)\)`)

	scanner := bufio.NewScanner(strings.NewReader(string(fileData.Data)))
	inSpecs := false

	for scanner.Scan() {
		line := scanner.Text()

		// Start at specs
		if strings.TrimSpace(line) == "specs:" {
			inSpecs = true
			continue
		}

		// Stop when we hit a new section
		if inSpecs && len(line) > 0 && line[0] != ' ' {
			break
		}

		if inSpecs {
			if matches := gemLineRegex.FindStringSubmatch(line); matches != nil {
				dependencies = append(dependencies, models.Dependency{
					Name:      matches[1],
					Version:   matches[2],
					Ecosystem: "gem",
				})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return dependencies, nil
}
