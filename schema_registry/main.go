package schema_registry

import (
	"embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/hamba/avro/v2"
)

type SchemaRegistry struct {
	V1 Schemas
}

type Schemas struct {
	Producer             string
	AccountCreatedSchema []byte
	AccountUpdatedSchema []byte
	TaskAssignedSchema   []byte
	TaskCompletedSchema  []byte
}

func NewSchemaRegistry(producer string) *SchemaRegistry {
	return &SchemaRegistry{
		V1: Schemas{
			Producer:             producer,
			AccountCreatedSchema: mustReadSchema("v1/", "account.created.json"),
			AccountUpdatedSchema: mustReadSchema("v1/", "account.updated.json"),
			TaskAssignedSchema:   mustReadSchema("v1/", "task.assigned.json"),
			TaskCompletedSchema:  mustReadSchema("v1/", "task.completed.json"),
		},
	}
}

//go:embed schemas/v1/* schemas/event.json
var f embed.FS

var EventSchema = mustReadSchema("", "event.json")

func MarshalAndValidate(schemaBytes []byte, v interface{}) ([]byte, error) {
	schema, err := avro.ParseBytes(schemaBytes)
	if err != nil {
		slog.Error("unable to parse schema", "error", err)
		return nil, err
	}
	bytes, err := avro.Marshal(schema, v)
	if err != nil {
		slog.Error("error while marshalling to avro bytes", "error", err)
		return nil, err
	}
	return bytes, nil
}

func UnmarshalAndValidate(schemaBytes []byte, bytes []byte, v interface{}) error {
	schema, err := avro.ParseBytes(schemaBytes)
	if err != nil {
		slog.Error("unable to parse schema", "error", err)
		return err
	}
	err = avro.Unmarshal(schema, bytes, &v)
	if err != nil {
		slog.Error("error while unmarshalling from avro bytes", "error", err)
		return err
	}
	return nil
}

func mustReadSchema(version, file string) []byte {
	path := fmt.Sprintf("schemas/%s%s", version, file)
	bytes, err := f.ReadFile(path)
	if err != nil {
		slog.Error("error while reading from file", "file", file, "error", err)
		os.Exit(1)
	}
	return bytes
}
