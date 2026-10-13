package jsonschema_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"

	jsonschema "github.com/faustbrian/go-json-schema/v2"
)

// FuzzCompileAndValidate exercises hostile schema/instance bytes across every
// released dialect, including pinned type seeds, large exponents and cycles.
// Compile/validation refusal is allowed; the target guards panic/termination
// boundaries, not agreement with an external specification oracle.
func FuzzCompileAndValidate(f *testing.F) {
	for _, seed := range []struct {
		dialect  uint8
		schema   string
		instance string
	}{
		{0, `true`, `null`},
		{2, `{"type":"integer"}`, `1e100000`},
		{5, `{"pattern":"(?<=a)b"}`, `"ab"`},
		{5, `{"oneOf":[{"type":"null"},{"type":"array"}]}`, `[]`},
		{5, `{"$defs":{"node":{"$ref":"#/$defs/node"}},"$ref":"#/$defs/node"}`, `0`},
		{5, `{"allOf":[` + strings.Repeat(`{},`, 31) + `{}` + `]}`, `{}`},
		{5, strings.Repeat(`{"properties":{"x":`, 32) + `{}` + strings.Repeat(`}}`, 32), strings.Repeat(`{"x":`, 32) + `0` + strings.Repeat(`}`, 32)},
	} {
		f.Add(seed.dialect, seed.schema, seed.instance)
	}
	addOfficialFuzzSeeds(f)
	f.Fuzz(func(t *testing.T, dialectIndex uint8, schemaJSON string, instanceJSON string) {
		dialects := [...]jsonschema.Dialect{
			jsonschema.Draft3,
			jsonschema.Draft4,
			jsonschema.Draft6,
			jsonschema.Draft7,
			jsonschema.Draft201909,
			jsonschema.Draft202012,
		}
		compiler, err := jsonschema.NewCompiler(
			jsonschema.WithDialect(dialects[int(dialectIndex)%len(dialects)]),
		)
		if err != nil {
			t.Fatal(err)
		}
		schema, err := compiler.Compile(context.Background(), []byte(schemaJSON))
		if err != nil {
			return
		}
		_, _ = schema.Validate(context.Background(), []byte(instanceJSON))
	})
}

func addOfficialFuzzSeeds(f *testing.F) {
	f.Helper()
	for dialectIndex, directory := range []string{
		"draft3", "draft4", "draft6", "draft7", "draft2019-09", "draft2020-12",
	} {
		// #nosec G304 -- directory is selected from the fixed dialect list.
		raw, err := os.ReadFile(filepath.Join(
			"testdata", "official", "JSON-Schema-Test-Suite", "tests", directory, "type.json",
		))
		if err != nil {
			f.Fatal(err)
		}
		var groups []fixtureGroup
		if err := json.Unmarshal(raw, &groups); err != nil {
			f.Fatal(err)
		}
		for _, group := range groups {
			for _, test := range group.Tests {
				f.Add(uint8(dialectIndex), string(group.Schema), string(test.Data))
			}
		}
	}
}

// FuzzValidationOutput traverses all four output forms for compiled hostile
// schemas. It tolerates parse/budget refusal and exercises encoding on success;
// it does not assert complete diagnostic equivalence or conformance.
func FuzzValidationOutput(f *testing.F) {
	f.Add(`{"properties":{"name":{"type":"string"}},"required":["name"]}`, `{"name":1}`)
	f.Add(`{"prefixItems":[{"type":"integer"}],"items":false}`, `[1,2]`)
	f.Fuzz(func(t *testing.T, schemaJSON string, instanceJSON string) {
		compiler, err := jsonschema.NewCompiler()
		if err != nil {
			t.Fatal(err)
		}
		schema, err := compiler.Compile(context.Background(), []byte(schemaJSON))
		if err != nil {
			return
		}
		for _, format := range []jsonschema.OutputFormat{
			jsonschema.OutputFlag,
			jsonschema.OutputBasic,
			jsonschema.OutputDetailed,
			jsonschema.OutputVerbose,
		} {
			output, err := schema.ValidateOutput(context.Background(), []byte(instanceJSON), format)
			if err == nil {
				_, _ = json.Marshal(output)
			}
		}
	})
}

// FuzzRawAndValueValidationAgree compares raw validation with a UseNumber
// decoded value under an empty schema. UTF-8, trailing-input and duplicate-name filters avoid
// intentional admission differences; this is not an independent JSON oracle.
func FuzzRawAndValueValidationAgree(f *testing.F) {
	for _, seed := range []string{
		`null`, `true`, `1e1000`, `"text"`, `[1,{"x":2}]`, `{"a":0.3}`,
		`{"x":1,"x":2}`, `{"outer":{"x":1,"\u0078":2}}`,
	} {
		f.Add(seed)
	}
	compiler, err := jsonschema.NewCompiler()
	if err != nil {
		f.Fatal(err)
	}
	schema, err := compiler.Compile(context.Background(), []byte(`{}`))
	if err != nil {
		f.Fatal(err)
	}
	f.Fuzz(func(t *testing.T, raw string) {
		// Validate rejects invalid UTF-8 while encoding/json deliberately
		// replaces it when constructing a caller-provided Go string.
		if !utf8.ValidString(raw) {
			return
		}
		decoder := json.NewDecoder(bytes.NewBufferString(raw))
		decoder.UseNumber()
		var value any
		if err := decoder.Decode(&value); err != nil {
			return
		}
		var trailing any
		if err := decoder.Decode(&trailing); err != io.EOF {
			return
		}
		// encoding/json collapses duplicate names, including escaped aliases.
		// Exclude only that intentional admission difference, not all parser errors.
		if hasDuplicateJSONMembers(raw) {
			return
		}
		rawResult, rawErr := schema.Validate(context.Background(), []byte(raw))
		valueResult, valueErr := schema.ValidateValue(context.Background(), value)
		if (rawErr == nil) != (valueErr == nil) {
			t.Fatalf("raw error %v; value error %v", rawErr, valueErr)
		}
		if rawErr == nil && rawResult.Valid != valueResult.Valid {
			t.Fatalf("raw valid %t; value valid %t", rawResult.Valid, valueResult.Valid)
		}
	})
}

func hasDuplicateJSONMembers(raw string) bool {
	type objectFrame struct {
		names    map[string]struct{}
		nameNext bool
	}
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.UseNumber()
	var frames []objectFrame
	for {
		token, err := decoder.Token()
		if err != nil {
			return false // The caller already admitted exactly one JSON value.
		}
		if delimiter, ok := token.(json.Delim); ok && (delimiter == '}' || delimiter == ']') {
			frames = frames[:len(frames)-1]
			continue
		}
		if len(frames) > 0 && frames[len(frames)-1].names != nil {
			frame := &frames[len(frames)-1]
			if frame.nameNext {
				name, ok := token.(string)
				if !ok {
					return false
				}
				if _, exists := frame.names[name]; exists {
					return true
				}
				frame.names[name] = struct{}{}
				frame.nameNext = false
				continue
			}
			frame.nameNext = true
		}
		if delimiter, ok := token.(json.Delim); ok {
			frame := objectFrame{}
			if delimiter == '{' {
				frame.names = make(map[string]struct{})
				frame.nameNext = true
			}
			frames = append(frames, frame)
		}
	}
}

func TestFuzzAgreementDuplicateAdmission(t *testing.T) {
	for _, test := range []struct {
		raw       string
		duplicate bool
	}{
		{`{"x":1,"x":2}`, true},
		{`{"outer":{"x":1,"\u0078":2}}`, true},
		{`[{"x":1,"x":2}]`, true},
		{`{"x":{"x":1},"y":{"x":2}}`, false},
		{`[{"x":1},{"x":2}]`, false},
		{`{"x":[1,2],"y":{},"z":"x"}`, false},
		{`[null,true,1e1000,"x"]`, false},
	} {
		if got := hasDuplicateJSONMembers(test.raw); got != test.duplicate {
			t.Errorf("duplicate=%t, want %t for %s", got, test.duplicate, test.raw)
		}
	}
}

// FuzzReferenceResolution exercises generated local definition names and
// escaped JSON Pointer references; no network/custom loader participates.
func FuzzReferenceResolution(f *testing.F) {
	f.Add("value", `"ok"`)
	f.Add("a~b/c", `0`)
	f.Fuzz(func(t *testing.T, property string, instance string) {
		schemaDocument, err := json.Marshal(map[string]any{
			"$defs": map[string]any{property: map[string]any{"const": json.Number("0")}},
			"$ref":  "#/$defs/" + escapeJSONPointer(property),
		})
		if err != nil {
			return
		}
		compiler, err := jsonschema.NewCompiler()
		if err != nil {
			t.Fatal(err)
		}
		schema, err := compiler.Compile(context.Background(), schemaDocument)
		if err != nil {
			return
		}
		_, _ = schema.Validate(context.Background(), []byte(instance))
	})
}

func escapeJSONPointer(value string) string {
	var result bytes.Buffer
	for _, character := range value {
		switch character {
		case '~':
			result.WriteString("~0")
		case '/':
			result.WriteString("~1")
		default:
			result.WriteRune(character)
		}
	}
	return result.String()
}
