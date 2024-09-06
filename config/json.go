package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"

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

func WriteJsonConfig(configPath, schemaName string, config any) (err error) {
	if len(schemaName) > 0 {
		val := reflect.ValueOf(config)
		typ := val.Type()
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			if tag := f.Tag.Get("json"); tag != "" {
				jsonTags := strings.Split(tag, ",")
				if jsonTags[0] == "$schema" {
					val.Field(i).SetString(schemaName)
					break
				}

			}
		}
	}

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

	return
}
