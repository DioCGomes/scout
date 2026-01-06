package rustparser

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/mlw157/scout/internal/models"
)

type RustParser struct{}

type CargoLock struct {
	Packages []Package `toml:"package"`
}

type Package struct {
	Name    string `toml:"name"`
	Version string `toml:"version"`
}

type FileData struct {
	Path string
	Data []byte
}

func NewRustParser() *RustParser {
	return &RustParser{}
}

func ReadFile(path string) (*FileData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return &FileData{
		Path: path,
		Data: data,
	}, nil
}

func (p *RustParser) ParseFile(path string) ([]models.Dependency, error) {
	fileData, err := ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cargoLock CargoLock
	_, err = toml.Decode(string(fileData.Data), &cargoLock)
	if err != nil {
		return nil, err
	}

	var dependencies []models.Dependency
	for _, pkg := range cargoLock.Packages {
		dependencies = append(dependencies, models.Dependency{
			Name:      pkg.Name,
			Version:   pkg.Version,
			Ecosystem: "crates.io",
		})
	}

	return dependencies, nil
}
