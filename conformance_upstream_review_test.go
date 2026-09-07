package jsonschema_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
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

type upstreamConformanceReview struct {
	Source              string                           `json:"source"`
	Compare             string                           `json:"compare"`
	BaseRevision        string                           `json:"base_revision"`
	HeadRevision        string                           `json:"head_revision"`
	ReviewedAt          string                           `json:"reviewed_at"`
	Classification      string                           `json:"classification"`
	Commits             []string                         `json:"commits"`
	SourceFiles         []upstreamConformanceReviewFile  `json:"source_files"`
	SupportedDialects   []string                         `json:"supported_dialects"`
	UnsupportedDialects []string                         `json:"unsupported_dialects"`
	Decisions           []string                         `json:"decisions"`
	Groups              []upstreamConformanceReviewGroup `json:"groups"`
}

type upstreamConformanceReviewFile struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type upstreamConformanceReviewGroup struct {
	Commit          string          `json:"commit"`
	Behavior        string          `json:"behavior"`
	Schema          json.RawMessage `json:"schema"`
	FormatAssertion bool            `json:"format_assertion"`
	Dialects        []string        `json:"dialects"`
	Tests           []fixtureCase   `json:"tests"`
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

func TestReviewedOfficialConformanceVectors(t *testing.T) {
	path := filepath.Join(
		"testdata",
		"regressions",
		"json-schema-test-suite-55e2372-to-c9510e3.json",
	)
	// #nosec G304 -- path is fixed to the reviewed upstream fixture range.
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var review upstreamConformanceReview
	if err := json.Unmarshal(raw, &review); err != nil {
		t.Fatal(err)
	}
	if review.BaseRevision != "55e23729473f4b629fd9266614280f355cd1b4fc" ||
		review.HeadRevision != "c9510e3bf8a896c3cba4e08509cf752b4f30dff8" {
		t.Fatalf(
			"reviewed range = %s...%s",
			review.BaseRevision,
			review.HeadRevision,
		)
	}
	if review.ReviewedAt != "2026-09-06" ||
		review.Classification != "conformance-significant behavior-neutral" {
		t.Fatalf(
			"review metadata = %q, %q",
			review.ReviewedAt,
			review.Classification,
		)
	}
	wantSource := "https://github.com/json-schema-org/JSON-Schema-Test-Suite"
	wantCompare := wantSource + "/compare/" + review.BaseRevision + "..." + review.HeadRevision
	if review.Source != wantSource || review.Compare != wantCompare {
		t.Fatalf("review source = %q, compare = %q", review.Source, review.Compare)
	}
	if !slices.Equal(review.Decisions, []string{"JSONSCHEMA-DEC-003", "JSONSCHEMA-DEC-015"}) {
		t.Fatalf("review decisions = %v", review.Decisions)
	}
	wantCommits := []string{
		"24012d0a2e3f974b53b9393b144ef8244fea9613",
		"5cef06b740b9c7b1ebb309afc24bf460dd397141",
		"c983956adb82dec1e3477f11322ff90647480970",
		"eebf50e89c02d7bd4976dc289fc6c629f5c6c6b0",
		"1872afdb1199cc9661157e3ef90b2d436a01e173",
		"bdfa72478b91f3e10270ed95187e448bb07620f7",
		"c9510e3bf8a896c3cba4e08509cf752b4f30dff8",
	}
	if !slices.Equal(review.Commits, wantCommits) {
		t.Fatalf("reviewed commits = %v, want %v", review.Commits, wantCommits)
	}
	wantSourceFiles := map[string]string{
		"tests/draft2019-09/enum.json":                 "d0a6df37fa174f37ea0540a9445f73c42de8a3c4138da4e8459b102f766e9cd9",
		"tests/draft2019-09/optional/format/iri.json":  "5289aa6b0931504491bdb5823fd5050050c2f671e7a986843f55e7e2ed193af4",
		"tests/draft2019-09/optional/format/uri.json":  "070492478e04ca141033c9b8fa5ebd4cab17d8c3caff67932969cda4327190d8",
		"tests/draft2019-09/optional/format/uuid.json": "7ab7d126a08d1f190297fcd62ca3bf5bd20fc2d526074bee8b003e4a579e5b79",
		"tests/draft2020-12/enum.json":                 "1abe7eb883bde4a78c62f08e484e16f32b72ca5c1be725aa54d85cc4546b5031",
		"tests/draft2020-12/optional/format/iri.json":  "f9522240705714cd53cc56217d04f5ecaabd12b3135a9e808ff02bb5089f4092",
		"tests/draft2020-12/optional/format/uri.json":  "47954ee6aef87c20a045ad2182fecd9c7eec53b6c8151ee88cba50381b2850a4",
		"tests/draft2020-12/optional/format/uuid.json": "d62b676e83a346a6fe4ad1812a4e8da1e7b99a15c5cfb9e9faa3168c37b9365e",
		"tests/draft3/enum.json":                       "96b218e653402e1e399ea7b0f80d6fddc41810a22f8898e3e8aa0062f53334d9",
		"tests/draft4/optional/format/hostname.json":   "271cd981abb58d28834e76a82f5f8435d2c00f0be17d66d2f4ec28a2200b661a",
		"tests/draft4/optional/format/uri.json":        "b35144c96aa5042d1464b14cabfdff5cf7d761bb34b84eb9b07ec54a68ffe1f9",
		"tests/draft6/enum.json":                       "069890108f4f83e62c9be9f42af63b5de911555e6644d83cd89c08ff644a43bb",
		"tests/draft6/optional/format/hostname.json":   "271cd981abb58d28834e76a82f5f8435d2c00f0be17d66d2f4ec28a2200b661a",
		"tests/draft6/optional/format/uri.json":        "b35144c96aa5042d1464b14cabfdff5cf7d761bb34b84eb9b07ec54a68ffe1f9",
		"tests/draft7/enum.json":                       "069890108f4f83e62c9be9f42af63b5de911555e6644d83cd89c08ff644a43bb",
		"tests/draft7/optional/format/iri.json":        "a814f74881cab8f91a9d38517869475b1ca292268c9853f4f4d8c04a85c01293",
		"tests/draft7/optional/format/uri.json":        "b35144c96aa5042d1464b14cabfdff5cf7d761bb34b84eb9b07ec54a68ffe1f9",
		"tests/v1/enum.json":                           "87d930b8863ddecdf71840886b367b8a141f0b2f88db8d719be2da280b660fa8",
		"tests/v1/format/iri.json":                     "8804ccf1a7100adf90c4fa2eb58dc4bc48a0c49fb2810310aaba2ae484ccbaa8",
		"tests/v1/format/uri.json":                     "02b400d8cc578ae6e96910a37d08257056203b2e97b7f5b5c3964144670bca2c",
		"tests/v1/format/uuid.json":                    "dfd532e255bc749a0a84a438e493f51135c4a190e20afd266d1e52fdc741a46e",
	}
	if len(review.SourceFiles) != len(wantSourceFiles) {
		t.Fatalf("reviewed source-file count = %d, want %d", len(review.SourceFiles), len(wantSourceFiles))
	}
	seenFiles := make(map[string]struct{}, len(wantSourceFiles))
	for _, source := range review.SourceFiles {
		if _, exists := seenFiles[source.Path]; exists {
			t.Fatalf("duplicate reviewed source path %q", source.Path)
		}
		seenFiles[source.Path] = struct{}{}
		if want, exists := wantSourceFiles[source.Path]; !exists || source.SHA256 != want {
			t.Fatalf("reviewed source %q has SHA-256 %q, want %q", source.Path, source.SHA256, want)
		}
	}
	wantDialects := []string{
		"draft3",
		"draft4",
		"draft6",
		"draft7",
		"draft2019-09",
		"draft2020-12",
	}
	if !slices.Equal(review.SupportedDialects, wantDialects) {
		t.Fatalf("supported dialects = %v, want %v", review.SupportedDialects, wantDialects)
	}
	if !slices.Equal(review.UnsupportedDialects, []string{"v1"}) {
		t.Fatalf("unsupported dialects = %v, want [v1]", review.UnsupportedDialects)
	}

	dialects := map[string]jsonschema.Dialect{
		"draft3":       jsonschema.Draft3,
		"draft4":       jsonschema.Draft4,
		"draft6":       jsonschema.Draft6,
		"draft7":       jsonschema.Draft7,
		"draft2019-09": jsonschema.Draft201909,
		"draft2020-12": jsonschema.Draft202012,
	}
	type groupContract struct {
		commit          string
		formatAssertion bool
		dialects        []string
		vectors         int
	}
	wantGroups := map[string]groupContract{
		"hostname overall length": {
			commit:          "24012d0a2e3f974b53b9393b144ef8244fea9613",
			formatAssertion: true,
			dialects:        []string{"draft4", "draft6"},
			vectors:         1,
		},
		"URI trailing newline": {
			commit:          "5cef06b740b9c7b1ebb309afc24bf460dd397141",
			formatAssertion: true,
			dialects:        []string{"draft4", "draft6", "draft7", "draft2019-09", "draft2020-12"},
			vectors:         1,
		},
		"enum codepoint equality": {
			commit:          "1872afdb1199cc9661157e3ef90b2d436a01e173",
			formatAssertion: false,
			dialects:        []string{"draft3", "draft6", "draft7", "draft2019-09", "draft2020-12"},
			vectors:         2,
		},
		"UUID braces": {
			commit:          "eebf50e89c02d7bd4976dc289fc6c629f5c6c6b0",
			formatAssertion: true,
			dialects:        []string{"draft2019-09", "draft2020-12"},
			vectors:         1,
		},
		"IRI percent encoding": {
			commit:          "c9510e3bf8a896c3cba4e08509cf752b4f30dff8",
			formatAssertion: true,
			dialects:        []string{"draft7", "draft2019-09", "draft2020-12"},
			vectors:         3,
		},
	}
	seenGroups := make(map[string]struct{}, len(wantGroups))
	vectorCount := 0
	for _, group := range review.Groups {
		want, exists := wantGroups[group.Behavior]
		if !exists || group.Commit != want.commit || group.FormatAssertion != want.formatAssertion ||
			!slices.Equal(group.Dialects, want.dialects) || len(group.Tests) != want.vectors {
			t.Fatalf("reviewed group %q has unexpected provenance or applicability", group.Behavior)
		}
		if _, duplicate := seenGroups[group.Behavior]; duplicate {
			t.Fatalf("duplicate reviewed group %q", group.Behavior)
		}
		seenGroups[group.Behavior] = struct{}{}
		vectorCount += len(group.Tests)
		for _, dialectName := range group.Dialects {
			dialect, exists := dialects[dialectName]
			if !exists {
				t.Fatalf("unreviewed or unsupported dialect %q", dialectName)
			}
			t.Run(fmt.Sprintf("%s/%s/%s", group.Commit, dialectName, group.Behavior), func(t *testing.T) {
				compiler, err := jsonschema.NewCompiler(jsonschema.WithDialect(dialect))
				if group.FormatAssertion {
					compiler, err = jsonschema.NewCompiler(
						jsonschema.WithDialect(dialect),
						jsonschema.WithFormatAssertion(),
					)
				}
				if err != nil {
					t.Fatal(err)
				}
				schema, err := compiler.Compile(context.Background(), group.Schema)
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
	if vectorCount != 8 {
		t.Fatalf("reviewed vector count = %d, want 8", vectorCount)
	}
	if len(seenGroups) != len(wantGroups) {
		t.Fatalf("reviewed group count = %d, want %d", len(seenGroups), len(wantGroups))
	}
}

func TestReviewedOfficialURITemplatePercentEncodingVectors(t *testing.T) {
	path := filepath.Join(
		"testdata",
		"regressions",
		"json-schema-test-suite-c9510e3-to-f6fd52a.json",
	)
	// #nosec G304 -- path is fixed to the reviewed upstream fixture range.
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var review upstreamConformanceReview
	if err := json.Unmarshal(raw, &review); err != nil {
		t.Fatal(err)
	}
	if review.BaseRevision != "c9510e3bf8a896c3cba4e08509cf752b4f30dff8" ||
		review.HeadRevision != "f6fd52a0a95472e079cbfc6ef7f089702b80e045" {
		t.Fatalf("reviewed range = %s...%s", review.BaseRevision, review.HeadRevision)
	}
	if review.ReviewedAt != "2026-09-07" ||
		review.Classification != "conformance-significant behavior-changing" {
		t.Fatalf("review metadata = %q, %q", review.ReviewedAt, review.Classification)
	}
	wantSource := "https://github.com/json-schema-org/JSON-Schema-Test-Suite"
	wantCompare := wantSource + "/compare/" + review.BaseRevision + "..." + review.HeadRevision
	if review.Source != wantSource || review.Compare != wantCompare {
		t.Fatalf("review source = %q, compare = %q", review.Source, review.Compare)
	}
	if !slices.Equal(review.Commits, []string{"f6fd52a0a95472e079cbfc6ef7f089702b80e045"}) {
		t.Fatalf("reviewed commits = %v", review.Commits)
	}
	if !slices.Equal(review.Decisions, []string{"JSONSCHEMA-DEC-003", "JSONSCHEMA-DEC-015"}) {
		t.Fatalf("review decisions = %v", review.Decisions)
	}
	wantFiles := map[string]string{
		"tests/draft2019-09/optional/format/uri-template.json": "12ef3e09f6ce0cc54ae43d98407eedd35d68299ea38a44aa2f572243064d9e77",
		"tests/draft2020-12/optional/format/uri-template.json": "8af3285cc0ff20fc79d9884c2e5e26e20c053ec6fc4081ca02995f8f70568b38",
		"tests/draft6/optional/format/uri-template.json":       "57f3209f2cf6ef0e2edad057a708850f74552c6599f3437f3a7e9578e6580f62",
		"tests/draft7/optional/format/uri-template.json":       "57f3209f2cf6ef0e2edad057a708850f74552c6599f3437f3a7e9578e6580f62",
		"tests/v1/format/uri-template.json":                    "d5c21031a6df543696fe6abf09f55a53ccf6ec0fd7d3b5c6f64eb16c31bfc7ba",
	}
	if len(review.SourceFiles) != len(wantFiles) {
		t.Fatalf("reviewed source-file count = %d, want %d", len(review.SourceFiles), len(wantFiles))
	}
	seenFiles := make(map[string]struct{}, len(wantFiles))
	for _, source := range review.SourceFiles {
		if _, exists := seenFiles[source.Path]; exists {
			t.Fatalf("duplicate reviewed source path %q", source.Path)
		}
		seenFiles[source.Path] = struct{}{}
		if want, exists := wantFiles[source.Path]; !exists || source.SHA256 != want {
			t.Fatalf("reviewed source %q has SHA-256 %q, want %q", source.Path, source.SHA256, want)
		}
	}
	wantDialects := []string{"draft6", "draft7", "draft2019-09", "draft2020-12"}
	if !slices.Equal(review.SupportedDialects, wantDialects) ||
		!slices.Equal(review.UnsupportedDialects, []string{"v1"}) {
		t.Fatalf("reviewed dialects = %v, unsupported = %v", review.SupportedDialects, review.UnsupportedDialects)
	}
	if len(review.Groups) != 1 {
		t.Fatalf("reviewed group count = %d, want 1", len(review.Groups))
	}
	group := review.Groups[0]
	if group.Commit != review.HeadRevision || group.Behavior != "URI-template percent encoding" ||
		!group.FormatAssertion || !slices.Equal(group.Dialects, wantDialects) || len(group.Tests) != 6 {
		t.Fatalf("reviewed group has unexpected provenance or applicability")
	}
	var reviewedSchema map[string]any
	if err := json.Unmarshal(group.Schema, &reviewedSchema); err != nil {
		t.Fatal(err)
	}
	if len(reviewedSchema) != 1 || reviewedSchema["format"] != "uri-template" {
		t.Fatalf("reviewed schema = %s", group.Schema)
	}
	wantVectors := map[string]string{
		"a literal with a lone percent sign is invalid":                         `"a%"`,
		"a literal with an incomplete percent-encoded triplet is invalid":       `"a%4"`,
		"a literal with non-hex percent encoding is invalid":                    `"a%GG"`,
		"a variable name with a lone percent sign is invalid":                   `"{%}"`,
		"a variable name with an incomplete percent-encoded triplet is invalid": `"{%4}"`,
		"a variable name with non-hex percent encoding is invalid":              `"{%GG}"`,
	}
	seenVectors := make(map[string]struct{}, len(wantVectors))
	for _, test := range group.Tests {
		want, exists := wantVectors[test.Description]
		if _, duplicate := seenVectors[test.Description]; duplicate {
			t.Fatalf("duplicate reviewed vector %q", test.Description)
		}
		seenVectors[test.Description] = struct{}{}
		if !exists || string(test.Data) != want || test.Valid {
			t.Fatalf("reviewed vector %q has data %s and valid=%t", test.Description, test.Data, test.Valid)
		}
	}

	dialects := map[string]jsonschema.Dialect{
		"draft6":       jsonschema.Draft6,
		"draft7":       jsonschema.Draft7,
		"draft2019-09": jsonschema.Draft201909,
		"draft2020-12": jsonschema.Draft202012,
	}
	for _, dialectName := range group.Dialects {
		dialectName := dialectName
		t.Run(dialectName, func(t *testing.T) {
			compiler, err := jsonschema.NewCompiler(
				jsonschema.WithDialect(dialects[dialectName]),
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
				result, err := schema.Validate(context.Background(), test.Data)
				if err != nil {
					t.Fatalf("%s: validate: %v", test.Description, err)
				}
				if result.Valid != test.Valid {
					t.Errorf("%s: got valid=%t, want %t", test.Description, result.Valid, test.Valid)
				}
			}
		})
	}
}
