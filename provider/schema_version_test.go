package maple

import (
	"encoding/json"
	"os"
	"testing"
)

// The served schema carries the plugin version, so `pulumi package add maple` can
// generate an SDK from it.
func TestServedSchemaHasVersion(t *testing.T) {
	raw, err := os.ReadFile("cmd/pulumi-resource-maple/schema.json")
	if err != nil {
		t.Fatal(err)
	}
	out, err := WithSchemaVersion(raw, "1.2.3")
	if err != nil {
		t.Fatal(err)
	}
	var s struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	if err := json.Unmarshal(out, &s); err != nil {
		t.Fatal(err)
	}
	if s.Name != "maple" || s.Version != "1.2.3" {
		t.Fatalf("name %q version %q, want maple 1.2.3", s.Name, s.Version)
	}
}
