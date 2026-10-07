package jsonschema_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	jsonschema "github.com/faustbrian/go-json-schema"
)

func TestReviewedOfficialRecentFormatVectors(t *testing.T) {
	raw, err := os.ReadFile("testdata/regressions/json-schema-test-suite-f6fd52a-to-5b0ee16.json")
	if err != nil {
		t.Fatal(err)
	}
	var review upstreamConformanceReview
	if err := json.Unmarshal(raw, &review); err != nil {
		t.Fatal(err)
	}
	if review.BaseRevision != "f6fd52a0a95472e079cbfc6ef7f089702b80e045" ||
		review.HeadRevision != "5b0ee1613e45fcc2bddac00e07c19cd49b00d8a8" {
		t.Fatalf("unexpected reviewed range %s...%s", review.BaseRevision, review.HeadRevision)
	}
	if len(review.Groups) == 0 {
		t.Fatal("review contains no executable format cases")
	}
	dialects := map[string]jsonschema.Dialect{
		"draft6":       jsonschema.Draft6,
		"draft7":       jsonschema.Draft7,
		"draft2019-09": jsonschema.Draft201909,
		"draft2020-12": jsonschema.Draft202012,
	}
	for _, group := range review.Groups {
		if !group.FormatAssertion || len(group.Dialects) == 0 || len(group.Tests) == 0 {
			t.Fatalf("%s: missing asserted format cases or dialects", group.Behavior)
		}
		for _, dialectName := range group.Dialects {
			t.Run(fmt.Sprintf("%s/%s", group.Behavior, dialectName), func(t *testing.T) {
				dialect, exists := dialects[dialectName]
				if !exists {
					t.Fatalf("unsupported reviewed dialect %q", dialectName)
				}
				compiler, err := jsonschema.NewCompiler(
					jsonschema.WithDialect(dialect),
					jsonschema.WithFormatAssertion(),
				)
				if err != nil {
					t.Fatal(err)
				}
				schema, err := compiler.Compile(context.Background(), group.Schema)
				if err != nil {
					t.Fatal(err)
				}
				for _, test := range group.Tests {
					t.Run(test.Description, func(t *testing.T) {
						result, err := schema.Validate(context.Background(), test.Data)
						if err != nil {
							t.Fatal(err)
						}
						if result.Valid != test.Valid {
							t.Fatalf("got valid=%t, want %t", result.Valid, test.Valid)
						}
					})
				}
			})
		}
	}
}
