package jsonschema_test

import (
	"context"
	"encoding/json"
	"fmt"

	jsonschema "github.com/faustbrian/go-json-schema"
)

func must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}

func Example() {
	compiler := must(jsonschema.NewCompiler(
		jsonschema.WithDialect(jsonschema.Draft202012),
	))
	schema := must(compiler.Compile(
		context.Background(),
		[]byte(`{
			"type": "object",
			"required": ["name"],
			"properties": {"name": {"type": "string"}}
		}`),
	))

	result := must(schema.Validate(context.Background(), []byte(`{"name":"Ada"}`)))
	fmt.Println(result.Valid)

	// Output: true
}

func ExampleSchema_ValidateValue() {
	compiler := must(jsonschema.NewCompiler())
	schema := must(compiler.Compile(
		context.Background(),
		[]byte(`{"type":"number","multipleOf":0.1}`),
	))

	result := must(schema.ValidateValue(context.Background(), json.Number("0.3")))
	fmt.Println(result.Valid)

	// Output: true
}

func ExampleMapLoader() {
	loader := must(jsonschema.NewMapLoader(map[string][]byte{
		"https://schemas.example.test/name": []byte(`{
			"$id":"https://schemas.example.test/name",
			"type":"string",
			"minLength":1
		}`),
	}))
	compiler := must(jsonschema.NewCompiler(jsonschema.WithResourceLoader(loader)))
	schema := must(compiler.Compile(
		context.Background(),
		[]byte(`{"$ref":"https://schemas.example.test/name"}`),
	))

	result := must(schema.Validate(context.Background(), []byte(`"Ada"`)))
	fmt.Println(result.Valid)

	// Output: true
}

func ExampleSchema_ValidateOutput() {
	compiler := must(jsonschema.NewCompiler())
	schema := must(compiler.Compile(context.Background(), []byte(`{"type":"string"}`)))

	output := must(schema.ValidateOutput(
		context.Background(),
		[]byte(`42`),
		jsonschema.OutputFlag,
	))
	encoded := must(json.Marshal(output))
	fmt.Println(string(encoded))

	// Output: {"valid":false}
}
