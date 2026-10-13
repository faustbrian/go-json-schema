package jsonschema_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	jsonschema "github.com/faustbrian/go-json-schema/v2"
)

func TestReferenceFailuresRedactEveryURIComponent(t *testing.T) {
	for _, suffix := range []string{"#missing-anchor", "#/missing", "#/array/01", "#/array/9", "#/scalar/child", "#fragment-marker%"} {
		t.Run(suffix, func(t *testing.T) {
			compiler, err := jsonschema.NewCompiler(jsonschema.WithResourceLoader(jsonschema.ResourceLoaderFunc(func(context.Context, string) ([]byte, error) {
				return []byte(`{"array":[true],"scalar":1}`), nil
			})))
			if err != nil {
				t.Fatal(err)
			}
			reference := "https://user-marker:password-marker@example.test/schema?query-marker=hidden" + suffix
			raw, _ := json.Marshal(map[string]string{"$ref": reference})
			_, err = compiler.Compile(context.Background(), raw)
			if !errors.Is(err, jsonschema.ErrInvalidSchema) {
				t.Fatal("expected invalid schema classification")
			}
			for _, marker := range []string{"user-marker", "password-marker", "query-marker", strings.TrimPrefix(suffix, "#")} {
				if strings.Contains(err.Error(), marker) {
					t.Fatal("reference diagnostic disclosed a sensitive component")
				}
			}
		})
	}
	compiler, err := jsonschema.NewCompiler(jsonschema.WithResourceLoader(jsonschema.ResourceLoaderFunc(func(context.Context, string) ([]byte, error) { return []byte(`{`), nil })))
	if err != nil {
		t.Fatal(err)
	}
	_, err = compiler.Compile(context.Background(), []byte(`{"$ref":"https://user-marker:password-marker@example.test/schema?query-marker=hidden"}`))
	if !errors.Is(err, jsonschema.ErrInvalidJSON) {
		t.Fatal("expected invalid JSON classification")
	}
	for _, marker := range []string{"user-marker", "password-marker", "query-marker"} {
		if strings.Contains(err.Error(), marker) {
			t.Fatal("resource diagnostic disclosed a sensitive component")
		}
	}
}

func TestExplicitNullAnnotationsRemainPresentOnWire(t *testing.T) {
	compiler, err := jsonschema.NewCompiler()
	if err != nil {
		t.Fatal(err)
	}
	for _, raw := range []string{`{"default":null}`, `{"$defs":{"x":{"default":null}},"$ref":"#/$defs/x"}`, `{"allOf":[{"default":null}]}`} {
		schema, err := compiler.Compile(context.Background(), []byte(raw))
		if err != nil {
			t.Fatal(err)
		}
		annotations, err := schema.CollectAnnotations(context.Background(), []byte(`null`))
		if err != nil || len(annotations) != 1 {
			t.Fatalf("collection count=%d err=%v", len(annotations), err)
		}
		assertNullAnnotationWire(t, annotations[0])
		for _, format := range []jsonschema.OutputFormat{jsonschema.OutputBasic, jsonschema.OutputDetailed, jsonschema.OutputVerbose} {
			for _, valuePath := range []bool{false, true} {
				var output jsonschema.OutputUnit
				if valuePath {
					output, err = schema.ValidateValueOutput(context.Background(), nil, format)
				} else {
					output, err = schema.ValidateOutput(context.Background(), []byte(`null`), format)
				}
				if err != nil {
					t.Fatal(err)
				}
				found := 0
				var visit func(jsonschema.OutputUnit)
				visit = func(unit jsonschema.OutputUnit) {
					if strings.HasSuffix(unit.KeywordLocation, "/default") {
						found++
						assertNullAnnotationWire(t, unit)
					}
					for _, child := range unit.Annotations {
						visit(child)
					}
					for _, child := range unit.Errors {
						visit(child)
					}
				}
				visit(output)
				if found != 1 {
					t.Fatalf("format=%s valuePath=%t: found %d null annotations", format, valuePath, found)
				}
			}
		}
	}
	absent, err := json.Marshal(jsonschema.OutputUnit{Valid: true})
	if err != nil || strings.Contains(string(absent), `"annotation"`) {
		t.Fatal("absent annotation was emitted")
	}
}

func assertNullAnnotationWire(t *testing.T, unit jsonschema.OutputUnit) {
	t.Helper()
	if unit.Annotation != nil {
		t.Fatal("explicit null must retain its nil Go value")
	}
	raw, err := json.Marshal(unit)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		t.Fatal(err)
	}
	if string(fields["annotation"]) != "null" {
		t.Fatal("explicit null annotation missing on wire")
	}
}

func TestRequiredAssertionVocabularyRejectsUnsupportedFormats(t *testing.T) {
	for _, required := range []bool{false, true} {
		meta, _ := json.Marshal(map[string]any{"$schema": string(jsonschema.Draft202012), "$vocabulary": map[string]bool{"https://json-schema.org/draft/2020-12/vocab/core": true, "https://json-schema.org/draft/2020-12/vocab/applicator": true, "https://json-schema.org/draft/2020-12/vocab/format-assertion": required}})
		for _, format := range []string{"unknown-format", "email", "registered-format"} {
			compiler, err := jsonschema.NewCompiler(jsonschema.WithResourceLoader(jsonschema.ResourceLoaderFunc(func(context.Context, string) ([]byte, error) { return meta, nil })), jsonschema.WithFormat("registered-format", jsonschema.FormatFunc(func(context.Context, string) (bool, error) { return false, nil })))
			if err != nil {
				t.Fatal(err)
			}
			raw, _ := json.Marshal(map[string]any{"$schema": "https://example.test/meta", "properties": map[string]any{"x": map[string]string{"format": format}}})
			schema, err := compiler.Compile(context.Background(), raw)
			if required && format == "unknown-format" {
				if !errors.Is(err, jsonschema.ErrInvalidSchema) {
					t.Fatal("required assertion vocabulary admitted unsupported format")
				}
			} else if err != nil {
				t.Fatal(err)
			} else if format == "registered-format" {
				result, err := schema.Validate(context.Background(), []byte(`{"x":"anything"}`))
				if err != nil || result.Valid {
					t.Fatal("registered format must remain active in recognized assertion vocabulary")
				}
			}
		}
	}
	compiler, err := jsonschema.NewCompiler(jsonschema.WithFormatAssertion())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := compiler.Compile(context.Background(), []byte(`{"format":"unknown-format"}`)); err != nil {
		t.Fatal("application opt-in must retain unknown-format annotations")
	}
}

func TestIdentifierFailuresRedactURIComponents(t *testing.T) {
	identifier := "https://user-marker:password-marker@example.test/schema?query-marker=hidden#fragment-marker"
	_, err := jsonschema.NewCompiler(jsonschema.WithDialect(jsonschema.Dialect(identifier)))
	if !errors.Is(err, jsonschema.ErrUnsupportedDialect) {
		t.Fatal("expected unsupported dialect")
	}
	for _, marker := range []string{"user-marker", "password-marker", "query-marker", "fragment-marker"} {
		if strings.Contains(err.Error(), marker) {
			t.Fatal("dialect error exposed sensitive URI component")
		}
	}
	meta, err := json.Marshal(map[string]any{"$schema": string(jsonschema.Draft202012), "$vocabulary": map[string]bool{identifier: true}})
	if err != nil {
		t.Fatal(err)
	}
	compiler, err := jsonschema.NewCompiler(jsonschema.WithResourceLoader(jsonschema.ResourceLoaderFunc(func(context.Context, string) ([]byte, error) { return meta, nil })))
	if err != nil {
		t.Fatal(err)
	}
	_, err = compiler.Compile(context.Background(), []byte(`{"$schema":"https://example.test/meta"}`))
	if !errors.Is(err, jsonschema.ErrUnsupportedVocabulary) {
		t.Fatal("expected unsupported vocabulary")
	}
	for _, marker := range []string{"user-marker", "password-marker", "query-marker", "fragment-marker"} {
		if strings.Contains(err.Error(), marker) {
			t.Fatal("vocabulary error exposed sensitive URI component")
		}
	}
}

func TestCustomNullAnnotationDiffersFromAbsentAnnotation(t *testing.T) {
	for _, annotation := range []json.RawMessage{nil, json.RawMessage(`null`)} {
		compiler, err := jsonschema.NewCompiler(jsonschema.WithDialect(jsonschema.Draft7), jsonschema.WithVocabulary("https://example.test/vocab", map[string]jsonschema.KeywordCompiler{
			"customNull": jsonschema.KeywordCompilerFunc(func(context.Context, jsonschema.Dialect, jsonschema.Value) (jsonschema.KeywordEvaluator, error) {
				return jsonschema.KeywordEvaluatorFunc(func(context.Context, jsonschema.Value) (jsonschema.KeywordResult, error) {
					return jsonschema.KeywordResult{Valid: true, Annotation: annotation}, nil
				}), nil
			}),
		}))
		if err != nil {
			t.Fatal(err)
		}
		schema, err := compiler.Compile(context.Background(), []byte(`{"customNull":true}`))
		if err != nil {
			t.Fatal(err)
		}
		units, err := schema.CollectAnnotations(context.Background(), []byte(`0`))
		if err != nil {
			t.Fatal(err)
		}
		if annotation == nil {
			if len(units) != 0 {
				t.Fatal("nil callback annotation must be absent")
			}
		} else {
			if len(units) != 1 {
				t.Fatal("null callback annotation must be present")
			}
			assertNullAnnotationWire(t, units[0])
		}
	}
}

func TestVocabularyRegistrationRedactsURIComponents(t *testing.T) {
	identifier := "https://user-marker:password-marker@example.test/vocab?query-marker=hidden"
	for _, options := range [][]jsonschema.Option{
		{jsonschema.WithVocabulary(identifier+"#fragment-marker", nil)},
		{jsonschema.WithVocabulary(identifier, nil), jsonschema.WithVocabulary(identifier, nil)},
	} {
		_, err := jsonschema.NewCompiler(options...)
		if !errors.Is(err, jsonschema.ErrInvalidSchema) {
			t.Fatal("expected invalid vocabulary registration")
		}
		for _, marker := range []string{"user-marker", "password-marker", "query-marker", "fragment-marker"} {
			if strings.Contains(err.Error(), marker) {
				t.Fatal("vocabulary registration exposed sensitive URI component")
			}
		}
	}
}
