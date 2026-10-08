package jsonschema_test

import (
	"context"
	"testing"

	jsonschema "github.com/faustbrian/go-json-schema/v2"
)

func TestHostnameNormalizationPreservesCanonicalALabels(t *testing.T) {
	t.Parallel()

	// Literal RFC 3492 labels distinguish NFC composition from character
	// substitution; expected labels do not come from the IDNA implementation.
	tests := []struct {
		name   string
		format string
		value  string
		valid  bool
	}{
		{"supplementary Han with grave", "hostname", "xn--ksa5191x.example", true},
		{"noncanonical Hangul and decomposed Latin", "hostname", "xn--e-xbb1090k.example", false},
		{"canonical Hangul and composed Latin", "hostname", "xn--9ca1178f.example", true},
		{"mapped Jamo and decomposed Latin", "idn-hostname", "\u1100\u1161e\u0301.example", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			compiler, err := jsonschema.NewCompiler(
				jsonschema.WithDialect(jsonschema.Draft202012),
				jsonschema.WithFormatAssertion(),
			)
			if err != nil {
				t.Fatal(err)
			}
			schema, err := compiler.Compile(context.Background(), []byte(`{"format":"`+test.format+`"}`))
			if err != nil {
				t.Fatal(err)
			}
			result, err := schema.Validate(context.Background(), []byte(`"`+test.value+`"`))
			if err != nil {
				t.Fatal(err)
			}
			if result.Valid != test.valid {
				t.Fatalf("%s validity = %t, want %t", test.value, result.Valid, test.valid)
			}
		})
	}
}
