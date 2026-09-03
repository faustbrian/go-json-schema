package jsonschema_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	jsonschema "github.com/faustbrian/go-json-schema"
)

type upstreamFormatReview struct {
	BaseRevision string                      `json:"base_revision"`
	HeadRevision string                      `json:"head_revision"`
	Groups       []upstreamFormatReviewGroup `json:"groups"`
}

type upstreamFormatReviewGroup struct {
	Commit   string        `json:"commit"`
	Format   string        `json:"format"`
	Dialects []string      `json:"dialects"`
	Tests    []fixtureCase `json:"tests"`
}

func TestReviewedOfficialFormatVectors(t *testing.T) {
	path := filepath.Join(
		"testdata",
		"regressions",
		"json-schema-test-suite-3c25e5f-to-55e2372.json",
	)
	// #nosec G304 -- path is fixed to the reviewed upstream fixture range.
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var review upstreamFormatReview
	if err := json.Unmarshal(raw, &review); err != nil {
		t.Fatal(err)
	}
	if review.BaseRevision != "3c25e5f709192aadf67cf7f2eb19771a57131fec" ||
		review.HeadRevision != "55e23729473f4b629fd9266614280f355cd1b4fc" {
		t.Fatalf(
			"reviewed range = %s...%s",
			review.BaseRevision,
			review.HeadRevision,
		)
	}
	vectorCount := 0
	for _, group := range review.Groups {
		vectorCount += len(group.Tests)
	}
	if vectorCount != 8 {
		t.Fatalf("reviewed vector count = %d, want 8", vectorCount)
	}

	dialects := map[string]jsonschema.Dialect{
		"draft4":       jsonschema.Draft4,
		"draft6":       jsonschema.Draft6,
		"draft7":       jsonschema.Draft7,
		"draft2019-09": jsonschema.Draft201909,
	}
	for _, group := range review.Groups {
		group := group
		for _, dialectName := range group.Dialects {
			dialectName := dialectName
			t.Run(fmt.Sprintf("%s/%s/%s", group.Commit, dialectName, group.Format), func(t *testing.T) {
				dialect, exists := dialects[dialectName]
				if !exists {
					t.Fatalf("unreviewed dialect %q", dialectName)
				}
				compiler, err := jsonschema.NewCompiler(
					jsonschema.WithDialect(dialect),
					jsonschema.WithFormatAssertion(),
				)
				if err != nil {
					t.Fatal(err)
				}
				schema, err := compiler.Compile(
					context.Background(),
					[]byte(`{"format":"`+group.Format+`"}`),
				)
				if err != nil {
					t.Fatal(err)
				}
				for _, test := range group.Tests {
					result, err := schema.Validate(context.Background(), test.Data)
					if err != nil {
						t.Fatalf("%s: validate: %v", test.Description, err)
					}
					if result.Valid != test.Valid {
						t.Errorf(
							"%s: got valid=%t, want %t",
							test.Description,
							result.Valid,
							test.Valid,
						)
					}
				}
			})
		}
	}
}
