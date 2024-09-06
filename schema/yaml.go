package schema

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"

	"fiurthorn.de/goutils/internal"
	"gopkg.in/yaml.v3"
)

func WriteYamlConfig(configPath string, config any) (schemaPath string, err error) {
	configName := filepath.Base(configPath)
	schemaPath = internal.ResolveSibling(configPath, ".schema")
	schemaName := filepath.Base(schemaPath)

	if f, err := os.OpenFile(configName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
		return "", err
	} else {
		defer f.Close()

		if _, err := f.WriteString(`# yaml-language-server: $schema=./` + schemaName + "\n\n"); err != nil {
			return "", err
		}

		y := yaml.NewEncoder(f)
		y.SetIndent(2)
		if err = y.Encode(config); err != nil {
			return "", err
		}
	}

	reflectType, err := Generate(reflect.TypeOf(config))
	if err != nil {
		return
	}
	schema, err := json.MarshalIndent(reflectType, "", " ")
	if err != nil {
		return
	}

	err = os.WriteFile(schemaPath, schema, 0644)
	return
}
