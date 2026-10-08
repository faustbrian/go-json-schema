package jsonschema_test

import (
	"context"
	"encoding/json"
	"testing"

	jsonschema "github.com/faustbrian/go-json-schema/v2"
)

func TestURITemplateExpressionModifiers(t *testing.T) {
	for _, dialect := range []jsonschema.Dialect{
		jsonschema.Draft6, jsonschema.Draft7, jsonschema.Draft201909, jsonschema.Draft202012,
	} {
		compiler, err := jsonschema.NewCompiler(jsonschema.WithDialect(dialect), jsonschema.WithFormatAssertion())
		if err != nil {
			t.Fatal(err)
		}
		schema, err := compiler.Compile(context.Background(), []byte(`{"format":"uri-template"}`))
		if err != nil {
			t.Fatal(err)
		}
		for _, test := range []struct {
			value string
			valid bool
		}{
			{"{var:1*}", false}, {"{+var:1*}", false},
			{"{var:1}", true}, {"{var:9999}", true}, {"{var*}", true},
			{"{?x,y*}", true}, {"{a.b,%41:2}", true},
			{"a\U000F0000b/{var}", true}, {"a\U00100000b/{var}", true},
			{"{}", false}, {"{a,,b}", false}, {"{a..b}", false},
			{"{a.}", false}, {"{var:0}", false}, {"{var:01}", false},
			{"{var:10000}", false}, {"{var*:1}", false},
			{"{var:}", false}, {"{var-}", false}, {"{var,.a}", false},
			{"{a,}", false}, {"{+}", false}, {"a{var", false}, {"a}var", false},
			{"{_A0%41:12}", true}, {"{,var}", true}, {"{|var*}", true},
			{"{z}", true}, {"{Z}", true}, {"{9}", true},
			{"{var:10}", true}, {"a{b}", true},
		} {
			raw, err := json.Marshal(test.value)
			if err != nil {
				t.Fatal(err)
			}
			result, err := schema.Validate(context.Background(), raw)
			if err != nil {
				t.Fatal(err)
			}
			if result.Valid != test.valid {
				t.Errorf("dialect %v template %q: valid=%t err=%v, want %t", dialect, test.value, result.Valid, err, test.valid)
			}
		}
	}
}
