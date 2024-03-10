package schema_registry

import (
	"embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/hamba/avro/v2"
)

type SchemaRegistry struct {
	V1 SchemaV1
	V2 SchemaV2
}

type SchemaV1 struct {
	Producer                 string
	Version                  string
	AccountCreatedSchema     avro.Schema
	AccountUpdatedSchema     avro.Schema
	TaskAssignedSchema       avro.Schema
	TaskCompletedSchema      avro.Schema
	PaymentMadeSchema        avro.Schema
	TransactionRevenueSchema avro.Schema
	TransactionCostSchema    avro.Schema
}

type SchemaV2 struct {
	Producer            string
	Version             string
	TaskAssignedSchema  avro.Schema
	TaskCompletedSchema avro.Schema
}

func NewSchemaRegistry(producer string) *SchemaRegistry {
	return &SchemaRegistry{
		V1: NewSchemaV1(producer),
		V2: NewSchemaV2(producer),
	}
}

func NewSchemaV1(producer string) SchemaV1 {
	return SchemaV1{
		Producer:                 producer,
		Version:                  "1",
		AccountCreatedSchema:     mustReadSchema("v1/account_created.json"),
		AccountUpdatedSchema:     mustReadSchema("v1/account_updated.json"),
		TaskAssignedSchema:       mustReadSchemaWithDeps("v1/task_assigned.json", "", map[string]string{"task": "v1/task.json"}),
		TaskCompletedSchema:      mustReadSchemaWithDeps("v1/task_completed.json", "", map[string]string{"task": "v1/task.json"}),
		PaymentMadeSchema:        mustReadSchema("v1/payment_made.json"),
		TransactionRevenueSchema: mustReadSchema("v1/transaction_revenue.json"),
		TransactionCostSchema:    mustReadSchema("v1/transaction_cost.json"),
	}
}

func NewSchemaV2(producer string) SchemaV2 {
	return SchemaV2{
		Producer:            producer,
		Version:             "2",
		TaskAssignedSchema:  mustReadSchemaWithDeps("v2/task_assigned.json", "", map[string]string{"task": "v2/task.json"}),
		TaskCompletedSchema: mustReadSchemaWithDeps("v2/task_completed.json", "", map[string]string{"task": "v2/task.json"}),
	}
}

const schemaFilesPrefix = "schemas/"

//go:embed schemas/v1/*.json schemas/v2/*.json
var f embed.FS

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
