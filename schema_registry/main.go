package schema_registry

import (
	"embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/hamba/avro/v2"
)

type SchemaRegistry struct {
	Producer             string
	AccountCreatedSchema avro.Schema
	AccountUpdatedSchema avro.Schema
	TaskAssignedSchema   avro.Schema
	TaskCompletedSchema  avro.Schema
}

func NewSchemaRegistry(producer string) *SchemaRegistry {
	return &SchemaRegistry{
		Producer:             producer,
		AccountCreatedSchema: mustReadSchema("account_created.json"),
		AccountUpdatedSchema: mustReadSchema("account_updated.json"),
		TaskAssignedSchema:   mustReadSchemaWithDeps("task_assigned.json", "", map[string]string{"task": "task.json"}),
		TaskCompletedSchema:  mustReadSchemaWithDeps("task_completed.json", "", map[string]string{"task": "task.json"}),
	}
}

const schemaFilesPrefix = "schemas/"

//go:embed schemas/*.json
var f embed.FS

var EventRawSchema = mustReadSchema("event_raw.json")

func MarshalAndValidate(schema avro.Schema, v interface{}) ([]byte, error) {
	bytes, err := avro.Marshal(schema, v)
	if err != nil {
		slog.Error("error while marshalling to avro bytes", "error", err)
		return nil, err
	}
	return bytes, nil
}

func UnmarshalAndValidate(schema avro.Schema, bytes []byte, v interface{}) error {
	err := avro.Unmarshal(schema, bytes, &v)
	if err != nil {
		slog.Error("error while unmarshalling from avro bytes", "error", err)
		return err
	}
	return nil
}

func mustReadSchema(file string) avro.Schema {
	path := fmt.Sprintf("%s%s", schemaFilesPrefix, file)
	bytes, err := f.ReadFile(path)
	if err != nil {
		slog.Error("error while reading from file", "file", file, "error", err)
		os.Exit(1)
	}
	schema, err := avro.ParseBytes(bytes)
	if err != nil {
		slog.Error("error while parsing avro schema", "file", file, "error", err)
		os.Exit(1)
	}
	return schema
}
func mustReadSchemaWithDeps(file, namespace string, deps map[string]string) avro.Schema {
	cache := &avro.SchemaCache{}
	for name, fileDep := range deps {
		path := fmt.Sprintf("%s%s", schemaFilesPrefix, fileDep)
		bytes, err := f.ReadFile(path)
		if err != nil {
			slog.Error("error while reading from file", "file", file, "error", err)
			os.Exit(1)
		}
		schema, err := avro.ParseBytesWithCache(bytes, namespace, cache)
		if err != nil {
			slog.Error("error while parsing avro schema", "file", file, "error", err)
			os.Exit(1)
		}
		cache.Add(name, schema)
	}

	path := fmt.Sprintf("%s%s", schemaFilesPrefix, file)
	bytes, err := f.ReadFile(path)
	if err != nil {
		slog.Error("error while reading from file", "file", file, "error", err)
		os.Exit(1)
	}
	schema, err := avro.ParseBytesWithCache(bytes, namespace, cache)
	if err != nil {
		slog.Error("error while parsing avro schema", "file", file, "error", err)
		os.Exit(1)
	}
	return schema
}
