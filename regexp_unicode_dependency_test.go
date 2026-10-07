package jsonschema_test

import (
	"context"
	"errors"
	"testing"

	jsonschema "github.com/faustbrian/go-json-schema"
)

func TestECMAScriptUnicodePropertyPatterns(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name   string
		schema string
		cases  []struct {
			instance string
			valid    bool
		}
	}{
		{
			name:   "string pattern",
			schema: `{"type":"string","pattern":"^\\p{Script=Greek}+$"}`,
			cases: []struct {
				instance string
				valid    bool
			}{
				{instance: `"α"`, valid: true},
				{instance: `"A"`, valid: false},
			},
		},
		{
			name:   "evaluated property routing",
			schema: `{"type":"object","patternProperties":{"^\\p{Script=Greek}+$":{"type":"integer"}},"unevaluatedProperties":false}`,
			cases: []struct {
				instance string
				valid    bool
			}{
				{instance: `{"α":1}`, valid: true},
				{instance: `{"α":"invalid"}`, valid: false},
				{instance: `{"A":1}`, valid: false},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			compiler, err := jsonschema.NewCompiler()
			if err != nil {
				t.Fatal(err)
			}
			schema, err := compiler.Compile(context.Background(), []byte(test.schema))
			if err != nil {
				t.Fatalf("compile standard Unicode property: %v", err)
			}
			for _, instance := range test.cases {
				result, err := schema.Validate(context.Background(), []byte(instance.instance))
				if err != nil || result.Valid != instance.valid {
					t.Fatalf("instance %s: valid=%t, err=%v, want valid=%t", instance.instance, result.Valid, err, instance.valid)
				}
			}
		})
	}
}

func TestECMAScriptPatternRejectsNonstandardGrammar(t *testing.T) {
	t.Parallel()

	for _, pattern := range []string{`^*`, `(?=a)*`, `\p{Greek}`, `\p{GCB=RI}`} {
		t.Run(pattern, func(t *testing.T) {
			compiler, err := jsonschema.NewCompiler()
			if err != nil {
				t.Fatal(err)
			}
			_, err = compiler.Compile(
				context.Background(),
				[]byte(`{"pattern":`+quoteJSON(pattern)+`}`),
			)
			if !errors.Is(err, jsonschema.ErrInvalidSchema) {
				t.Fatalf("nonstandard Unicode-mode expression %q: got %v, want ErrInvalidSchema", pattern, err)
			}
		})
	}
}

func TestAssertedRegexFormatUsesECMAScriptPropertyGrammar(t *testing.T) {
	t.Parallel()

	compiler, err := jsonschema.NewCompiler(jsonschema.WithFormatAssertion())
	if err != nil {
		t.Fatal(err)
	}
	schema, err := compiler.Compile(context.Background(), []byte(`{"format":"regex"}`))
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		pattern string
		valid   bool
	}{
		{pattern: `\p{Script=Greek}`, valid: true},
		{pattern: `\p{Greek}`, valid: false},
		{pattern: `\p{GCB=RI}`, valid: false},
		{pattern: `^*`, valid: false},
	} {
		t.Run(test.pattern, func(t *testing.T) {
			result, err := schema.Validate(context.Background(), []byte(quoteJSON(test.pattern)))
			if err != nil || result.Valid != test.valid {
				t.Fatalf("regex %q: valid=%t, err=%v, want valid=%t", test.pattern, result.Valid, err, test.valid)
			}
		})
	}
}

func TestECMAScriptBackreferenceSchemaRepeatedConcurrentUse(t *testing.T) {
	t.Parallel()

	compiler, err := jsonschema.NewCompiler()
	if err != nil {
		t.Fatal(err)
	}
	schema, err := compiler.Compile(context.Background(), []byte(`{"type":"string","pattern":"^(a)\\1$"}`))
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		instance string
		valid    bool
	}{
		{instance: `"aa"`, valid: true},
		{instance: `"ab"`, valid: false},
	} {
		t.Run(test.instance, func(t *testing.T) {
			t.Parallel()
			for range 8 {
				result, err := schema.Validate(context.Background(), []byte(test.instance))
				if err != nil || result.Valid != test.valid {
					t.Fatalf("instance %s: valid=%t, err=%v, want valid=%t", test.instance, result.Valid, err, test.valid)
				}
			}
		})
	}
}
