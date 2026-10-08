package comparison_test

import (
	"errors"
	"testing"

	owned "github.com/faustbrian/go-json-schema/v2"
)

func TestOwnedCompilerUsesReleasedUnicodeGrammar(t *testing.T) {
	compiler, err := owned.NewCompiler()
	if err != nil {
		t.Fatal(err)
	}
	schema, err := compiler.Compile(t.Context(), []byte(`{"type":"string","pattern":"^\\p{Script=Greek}+$"}`))
	if err != nil {
		t.Fatalf("released Unicode grammar compile: %v", err)
	}
	for _, test := range []struct {
		instance string
		valid    bool
	}{
		{`"αβ"`, true}, {`"AB"`, false},
	} {
		result, err := schema.Validate(t.Context(), []byte(test.instance))
		if err != nil || result.Valid != test.valid {
			t.Fatalf("released Unicode grammar result: valid=%t error=%v want=%t", result.Valid, err, test.valid)
		}
	}
	if _, err := compiler.Compile(t.Context(), []byte(`{"pattern":"\\p{Greek}"}`)); !errors.Is(err, owned.ErrInvalidSchema) {
		t.Fatalf("nonstandard Unicode grammar classification: %v", err)
	}
}

func TestOwnedCompilerUsesReleasedURITemplateGrammar(t *testing.T) {
	compiler, err := owned.NewCompiler(owned.WithFormatAssertion())
	if err != nil {
		t.Fatal(err)
	}
	schema, err := compiler.Compile(t.Context(), []byte(`{"format":"uri-template"}`))
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		instance string
		valid    bool
	}{
		{`"{var:1}"`, true}, {`"{var*}"`, true}, {`"{var:1*}"`, false},
	} {
		result, err := schema.Validate(t.Context(), []byte(test.instance))
		if err != nil || result.Valid != test.valid {
			t.Fatalf("released URI-template grammar result: valid=%t error=%v want=%t", result.Valid, err, test.valid)
		}
	}
}
