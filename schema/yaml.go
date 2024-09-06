package schema

import (
	"os"
	"path/filepath"

	"fiurthorn.de/goutils/internal"
	"gopkg.in/yaml.v3"
)

func WriteYamlConfigAndSchema(configPath string, config any) (err error) {
	configName := filepath.Base(configPath)
	schemaPath := internal.ResolveSibling(configPath, ".schema")
	schemaName := filepath.Base(schemaPath)

	if f, err := os.OpenFile(configName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
		return err
	} else {
		defer f.Close()

		if _, err := f.WriteString(`# yaml-language-server: $schema=./` + schemaName + "\n\n"); err != nil {
			return err
		}

		y := yaml.NewEncoder(f)
		y.SetIndent(2)
		if err = y.Encode(config); err != nil {
			return err
		}
	}

	err = WriteJsonSchema(schemaPath, config)
	return
}
