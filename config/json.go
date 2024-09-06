package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"

	"fiurthorn.de/goutils/internal"
)

func GetSchemaSibling(configPath string) (schemaPath string) {
	return internal.ResolveSibling(configPath, ".schema")
}

func WriteJsonSchema(schemaPath string, config any) (err error) {
	reflectType, err := GenerateSchema(reflect.TypeOf(config))
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

func WriteJsonConfigAndSchema(configPath, schemaPath string, config any) (err error) {
	configName := filepath.Base(configPath)

	if f, err := os.OpenFile(configName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
		return err
	} else {
		defer f.Close()

		schema, err := json.MarshalIndent(config, "", " ")
		if err != nil {
			return err
		}

		if _, err := f.Write(schema); err != nil {
			return err
		}
	}

	err = WriteJsonSchema(schemaPath, config)
	return
}
