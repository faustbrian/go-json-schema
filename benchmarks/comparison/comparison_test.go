// Package comparison_test compares a pinned published owned release with
// maintained peers. It does not execute this worktree's root implementation;
// historical results/fixtures remain distinct from current source qualification.
package comparison_test

import (
	"bytes"
	"context"
	"testing"

	owned "github.com/faustbrian/go-json-schema/v2"
	kaptin "github.com/kaptinlin/jsonschema"
	tekuri "github.com/santhosh-tekuri/jsonschema/v6"
)

var schemaJSON = []byte(`{
	"$schema":"https://json-schema.org/draft/2020-12/schema",
	"type":"object",
	"required":["id","tags"],
	"properties":{
		"id":{"type":"integer","minimum":1},
		"tags":{
			"type":"array",
			"items":{"type":"string","minLength":2},
			"uniqueItems":true
		}
	},
	"additionalProperties":false
}`)

var instanceValue = map[string]any{
	"id": 42,
	"tags": []any{
		"alpha", "beta", "gamma", "delta", "epsilon",
	},
}

// BenchmarkCompile includes constructor/registration/compilation for each
// pinned implementation. These setup boundaries differ from the root reused-
// compiler benchmark and must be preserved when comparing results.
func BenchmarkCompile(b *testing.B) {
	b.Run("faustbrian", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			compiler, err := owned.NewCompiler()
			if err != nil {
				b.Fatal(err)
			}
			if _, err := compiler.Compile(
				context.Background(), schemaJSON,
			); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("kaptinlin", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			if _, err := kaptin.NewCompiler().Compile(schemaJSON); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("santhosh-tekuri", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			document, err := tekuri.UnmarshalJSON(bytes.NewReader(schemaJSON))
			if err != nil {
				b.Fatal(err)
			}
			compiler := tekuri.NewCompiler()
			if err := compiler.AddResource("schema.json", document); err != nil {
				b.Fatal(err)
			}
			if _, err := compiler.Compile("schema.json"); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkValidateDecoded uses precompiled schemas and one decoded object.
// The owned ValidateValue path includes encoding/parsing while peers consume
// decoded values; results compare those APIs, not identical internal work.
func BenchmarkValidateDecoded(b *testing.B) {
	ownedSchema := compileOwned(b)
	kaptinSchema := compileKaptin(b)
	tekuriSchema := compileTekuri(b)

	b.Run("faustbrian", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			result, err := ownedSchema.ValidateValue(
				context.Background(), instanceValue,
			)
			if err != nil || !result.Valid {
				b.Fatalf("result %#v: %v", result, err)
			}
		}
	})
	b.Run("kaptinlin", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			if result := kaptinSchema.ValidateMap(instanceValue); !result.IsValid() {
				b.Fatalf("invalid result: %#v", result)
			}
		}
	})
	b.Run("santhosh-tekuri", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			if err := tekuriSchema.Validate(instanceValue); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func compileOwned(tb testing.TB) *owned.Schema {
	tb.Helper()
	compiler, err := owned.NewCompiler()
	if err != nil {
		tb.Fatal(err)
	}
	schema, err := compiler.Compile(context.Background(), schemaJSON)
	if err != nil {
		tb.Fatal(err)
	}
	return schema
}

func compileKaptin(tb testing.TB) *kaptin.Schema {
	tb.Helper()
	schema, err := kaptin.NewCompiler().Compile(schemaJSON)
	if err != nil {
		tb.Fatal(err)
	}
	return schema
}

func compileTekuri(tb testing.TB) *tekuri.Schema {
	tb.Helper()
	document, err := tekuri.UnmarshalJSON(bytes.NewReader(schemaJSON))
	if err != nil {
		tb.Fatal(err)
	}
	compiler := tekuri.NewCompiler()
	if err := compiler.AddResource("schema.json", document); err != nil {
		tb.Fatal(err)
	}
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		tb.Fatal(err)
	}
	return schema
}
