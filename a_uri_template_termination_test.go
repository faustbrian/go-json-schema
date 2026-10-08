package jsonschema_test

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"testing"
	"time"

	jsonschema "github.com/faustbrian/go-json-schema"
)

func TestURITemplateValidationTerminatesForFiniteExpressions(t *testing.T) {
	const helper = "JSONSCHEMA_URI_TERMINATION_HELPER"
	if os.Getenv(helper) != "1" {
		executable, err := os.Executable()
		if err != nil {
			t.Fatal(err)
		}
		// Contain a failed termination assertion in a killable, waited process;
		// this is a test budget, not a public validation latency guarantee.
		ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
		defer cancel()
		command := exec.CommandContext(ctx, executable,
			"-test.run=^TestURITemplateValidationTerminatesForFiniteExpressions$",
			"-test.count=1")
		command.Env = append(os.Environ(), helper+"=1")
		command.WaitDelay = time.Second
		output, err := command.CombinedOutput()
		if ctx.Err() != nil {
			t.Fatal("finite URI-template validation did not terminate within the test budget")
		}
		if err != nil {
			t.Fatalf("finite URI-template validation failed: %v\n%s", err, output)
		}
		return
	}

	for _, dialect := range []jsonschema.Dialect{
		jsonschema.Draft6, jsonschema.Draft7, jsonschema.Draft201909, jsonschema.Draft202012,
	} {
		compiler, err := jsonschema.NewCompiler(jsonschema.WithDialect(dialect), jsonschema.WithFormatAssertion())
		if err != nil {
			t.Fatal(err)
		}
		schema, err := compiler.Compile(t.Context(), []byte(`{"format":"uri-template"}`))
		if err != nil {
			t.Fatal(err)
		}
		for _, value := range []string{"a{b}", "{abc%41}", "{a.b}"} {
			raw, err := json.Marshal(value)
			if err != nil {
				t.Fatal(err)
			}
			result, err := schema.Validate(t.Context(), raw)
			if err != nil || !result.Valid {
				t.Fatalf("dialect %v template %q: valid=%t error=%v", dialect, value, result.Valid, err)
			}
		}
	}
}
