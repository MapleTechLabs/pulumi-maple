package maple

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"
	shim "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfshim"
)

type packageSchema struct {
	Resources map[string]struct {
		Properties map[string]struct {
			Secret bool `json:"secret"`
		} `json:"properties"`
	} `json:"resources"`
	Functions map[string]struct {
		Outputs struct {
			Properties map[string]struct {
				Secret bool `json:"secret"`
			} `json:"properties"`
		} `json:"outputs"`
	} `json:"functions"`
}

func readSchema(t *testing.T) packageSchema {
	t.Helper()
	raw, err := os.ReadFile("cmd/pulumi-resource-maple/schema.json")
	if err != nil {
		t.Fatal(err)
	}
	var s packageSchema
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatal(err)
	}
	return s
}

// Every Terraform resource and data source reaches the committed schema, and every
// attribute Terraform marks sensitive is a Pulumi secret there.
func TestUpstreamSchemaIsFullyMappedWithSecrets(t *testing.T) {
	prov := Provider()
	s := readSchema(t)

	prov.P.ResourcesMap().Range(func(name string, res shim.Resource) bool {
		info := prov.Resources[name]
		if info == nil {
			t.Errorf("resource %s has no Pulumi mapping", name)
			return true
		}
		spec, ok := s.Resources[string(info.Tok)]
		if !ok {
			t.Errorf("resource %s (%s) is missing from schema.json; run make schema", name, info.Tok)
			return true
		}
		res.Schema().Range(func(attr string, sch shim.Schema) bool {
			prop := tfbridge.TerraformToPulumiNameV2(attr, res.Schema(), info.Fields)
			if sch.Sensitive() && !spec.Properties[prop].Secret {
				t.Errorf("%s.%s is sensitive in Terraform but not a secret in Pulumi", info.Tok, prop)
			}
			return true
		})
		return true
	})

	prov.P.DataSourcesMap().Range(func(name string, ds shim.Resource) bool {
		info := prov.DataSources[name]
		if info == nil {
			t.Errorf("data source %s has no Pulumi mapping", name)
			return true
		}
		spec, ok := s.Functions[string(info.Tok)]
		if !ok {
			t.Errorf("data source %s (%s) is missing from schema.json; run make schema", name, info.Tok)
			return true
		}
		ds.Schema().Range(func(attr string, sch shim.Schema) bool {
			prop := tfbridge.TerraformToPulumiNameV2(attr, ds.Schema(), info.Fields)
			if sch.Sensitive() && !spec.Outputs.Properties[prop].Secret {
				t.Errorf("%s.%s is sensitive in Terraform but not a secret in Pulumi", info.Tok, prop)
			}
			return true
		})
		return true
	})
}
