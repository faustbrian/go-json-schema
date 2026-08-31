package comparison_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	owned "github.com/faustbrian/go-json-schema"
	kaptin "github.com/kaptinlin/jsonschema"
	tekuri "github.com/santhosh-tekuri/jsonschema/v6"
)

const differentialFixture = "../../specification/differential/maintained-peers.json"

type differentialReport struct {
	SchemaVersion   int                       `json:"schema_version"`
	ObservedAt      string                    `json:"observed_at"`
	Harness         string                    `json:"harness"`
	Implementations map[string]implementation `json:"implementations"`
	Cases           []differentialCase        `json:"cases"`
}

type implementation struct {
	Module  string `json:"module"`
	Version string `json:"version"`
}

type differentialCase struct {
	ID             string             `json:"id"`
	Decision       string             `json:"decision"`
	Schema         string             `json:"schema"`
	Instance       string             `json:"instance"`
	Outcomes       map[string]outcome `json:"outcomes"`
	Classification string             `json:"classification"`
	Rationale      string             `json:"rationale"`
}

type outcome struct {
	Schema     string `json:"schema"`
	Instance   string `json:"instance"`
	Validation string `json:"validation"`
}

func TestMaintainedPeerDecisionOutcomes(t *testing.T) {
	report := loadDifferentialReport(t)
	assertPeerVersions(t, report.Implementations)

	runners := map[string]func([]byte, []byte) outcome{
		"faustbrian":      runOwned,
		"kaptinlin":       runKaptin,
		"santhosh-tekuri": runTekuri,
	}
	for _, test := range report.Cases {
		t.Run(test.ID, func(t *testing.T) {
			for name, run := range runners {
				got := run([]byte(test.Schema), []byte(test.Instance))
				if want := test.Outcomes[name]; !reflect.DeepEqual(got, want) {
					t.Errorf("%s outcome = %+v, want %+v", name, got, want)
				}
			}
		})
	}
}

func loadDifferentialReport(t *testing.T) differentialReport {
	t.Helper()
	path := filepath.Clean(differentialFixture)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var report differentialReport
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&report); err != nil {
		t.Fatal(err)
	}
	if report.SchemaVersion != 1 || report.ObservedAt == "" || report.Harness != "benchmarks/comparison:TestMaintainedPeerDecisionOutcomes" || len(report.Cases) == 0 {
		t.Fatalf("invalid differential report: schema=%d observed=%q harness=%q cases=%d", report.SchemaVersion, report.ObservedAt, report.Harness, len(report.Cases))
	}
	seen := make(map[string]struct{}, len(report.Cases))
	for _, test := range report.Cases {
		if _, duplicate := seen[test.ID]; test.ID == "" || duplicate || test.Decision == "" || test.Schema == "" || test.Instance == "" || test.Classification != "deliberate policy difference" || test.Rationale == "" || len(test.Outcomes) != 3 {
			t.Fatalf("incomplete or duplicate differential case %q", test.ID)
		}
		seen[test.ID] = struct{}{}
	}
	return report
}

func assertPeerVersions(t *testing.T, implementations map[string]implementation) {
	t.Helper()
	command := exec.CommandContext(t.Context(), "go", "list", "-m", "-json", "all")
	output, err := command.Output()
	if err != nil {
		t.Fatalf("list selected modules: %v", err)
	}
	versions := make(map[string]string)
	decoder := json.NewDecoder(bytes.NewReader(output))
	for {
		var module struct {
			Path    string
			Version string
		}
		if err := decoder.Decode(&module); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			t.Fatalf("decode selected modules: %v", err)
		}
		versions[module.Path] = module.Version
	}
	for _, name := range []string{"kaptinlin", "santhosh-tekuri"} {
		identity := implementations[name]
		if got := versions[identity.Module]; got != identity.Version {
			t.Errorf("%s version = %q, want %q", identity.Module, got, identity.Version)
		}
	}
}

func runOwned(schemaJSON, instanceJSON []byte) outcome {
	compiler, err := owned.NewCompiler()
	if err != nil {
		return rejectedSchema()
	}
	schema, err := compiler.Compile(context.Background(), schemaJSON)
	if err != nil {
		return rejectedSchema()
	}
	result, err := schema.Validate(context.Background(), instanceJSON)
	if err != nil {
		return outcome{Schema: "accepted", Instance: "rejected", Validation: "not-run"}
	}
	return validationOutcome(result.Valid)
}

func runKaptin(schemaJSON, instanceJSON []byte) outcome {
	schema, err := kaptin.NewCompiler().Compile(schemaJSON)
	if err != nil {
		return rejectedSchema()
	}
	result := schema.ValidateJSON(instanceJSON)
	if invalidJSON := result.Errors["format"]; invalidJSON != nil && invalidJSON.Code == "invalid_json" {
		return outcome{Schema: "accepted", Instance: "rejected", Validation: "not-run"}
	}
	return validationOutcome(result.IsValid())
}

func runTekuri(schemaJSON, instanceJSON []byte) outcome {
	document, err := tekuri.UnmarshalJSON(bytes.NewReader(schemaJSON))
	if err != nil {
		return rejectedSchema()
	}
	compiler := tekuri.NewCompiler()
	if err := compiler.AddResource("schema.json", document); err != nil {
		return rejectedSchema()
	}
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return rejectedSchema()
	}
	instance, err := tekuri.UnmarshalJSON(bytes.NewReader(instanceJSON))
	if err != nil {
		return outcome{Schema: "accepted", Instance: "rejected", Validation: "not-run"}
	}
	return validationOutcome(schema.Validate(instance) == nil)
}

func rejectedSchema() outcome {
	return outcome{Schema: "rejected", Instance: "not-run", Validation: "not-run"}
}

func validationOutcome(valid bool) outcome {
	if valid {
		return outcome{Schema: "accepted", Instance: "accepted", Validation: "valid"}
	}
	return outcome{Schema: "accepted", Instance: "accepted", Validation: "invalid"}
}

func (value outcome) String() string {
	return fmt.Sprintf("schema=%s instance=%s validation=%s", value.Schema, value.Instance, value.Validation)
}
