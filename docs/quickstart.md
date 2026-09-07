# Quick start

Install the module using the version selected by your application's dependency
policy:

```sh
go get github.com/faustbrian/go-json-schema@v1
```

Compile once and validate many times. This exact program is compiled and run
by the documentation gate. The package-level [`Example`](../example_test.go)
exercises the same workflow through Go's executable-example contract:

```go
package main

import (
	"context"
	"fmt"

	jsonschema "github.com/faustbrian/go-json-schema"
)

func main() {
	compiler, err := jsonschema.NewCompiler(
		jsonschema.WithDialect(jsonschema.Draft202012),
	)
	if err != nil {
		panic(err)
	}

	schema, err := compiler.Compile(context.Background(), []byte(`{
		"type": "object",
		"required": ["name"],
		"properties": {"name": {"type": "string"}}
	}`))
	if err != nil {
		panic(err)
	}

	result, err := schema.Validate(
		context.Background(),
		[]byte(`{"name":"Ada"}`),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Valid)
}

// Output: true
```

Schema mismatch is reported as `Valid == false`. A non-nil error means input,
resource, cancellation, extension, or limit processing prevented a decision.

Next, choose a [dialect](dialects.md), configure [resource loading](resolvers.md)
when `$ref` crosses documents, and select [output](output.md) when callers need
machine-readable diagnostics.
