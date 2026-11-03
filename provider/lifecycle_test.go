package maple

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MapleTechLabs/terraform-provider-maple/fakeapi"
	"github.com/pulumi/providertest/providers"
	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/assertpreview"
	"github.com/pulumi/providertest/pulumitest/assertup"
	"github.com/pulumi/providertest/pulumitest/opttest"
	pfbridge "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfbridge"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

// newProgram runs a Pulumi YAML program against the in-process provider and fake API.
// YAML needs no generated SDK, so these tests exercise the bridge alone.
func newProgram(t *testing.T, api *fakeapi.API, program string) *pulumitest.PulumiTest {
	t.Helper()
	schema, err := os.ReadFile("cmd/pulumi-resource-maple/schema.json")
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	yaml := "name: maple-test\nruntime: yaml\n" + program
	if err := os.WriteFile(filepath.Join(dir, "Pulumi.yaml"), []byte(yaml), 0o600); err != nil {
		t.Fatal(err)
	}
	server := func(providers.PulumiTest) (pulumirpc.ResourceProviderServer, error) {
		return pfbridge.NewProviderServer(context.Background(), nil, Provider(), pfbridge.ProviderMetadata{PackageSchema: schema})
	}
	pt := pulumitest.NewPulumiTest(t, dir,
		opttest.TestInPlace(),
		opttest.SkipInstall(),
		opttest.AttachProviderServer("maple", server),
	)
	pt.SetConfig(t, "maple:apiUrl", api.URL())
	pt.SetConfig(t, "maple:apiKey", fakeapi.Key)
	return pt
}

const dashboardProgram = `
config:
  dashName:
    type: string
resources:
  ops:
    type: maple:Dashboard
    properties:
      name: ${dashName}
      tags: [production]
      widgets:
        fn::toJSON:
          - id: w1
            visualization: chart
            layout: { x: 0, y: 0, w: 6, h: 4 }
outputs:
  id: ${ops.id}
  timeRange: ${ops.timeRange}
`

func TestDashboardLifecycle(t *testing.T) {
	api := fakeapi.New(t)
	pt := newProgram(t, api, dashboardProgram)
	pt.SetConfig(t, "dashName", "Operations")

	up := pt.Up(t)
	id, _ := up.Outputs["id"].Value.(string)
	if !strings.HasPrefix(id, "dash_") {
		t.Fatalf("id = %q, want dash_ prefix", id)
	}
	// The API's default for an unset attribute is read back.
	if got := up.Outputs["timeRange"].Value; got != `{"type":"relative","value":"12h"}` {
		t.Errorf("timeRange = %v", got)
	}
	// The API returns widgets in its own form; the program's JSON must not drift.
	assertpreview.HasNoChanges(t, pt.Preview(t))

	pt.SetConfig(t, "dashName", "Operations v2")
	up = pt.Up(t)
	assertup.HasNoReplacements(t, up)
	if got := api.LastBody("PATCH", "/v2/dashboards/")["name"]; got != "Operations v2" {
		t.Errorf("PATCH name = %v", got)
	}

	// A change made in the Maple UI shows up on refresh and is reverted by up.
	api.Mutate("/v2/dashboards/"+id, func(obj map[string]any) { obj["name"] = "Renamed in UI" })
	pt.Refresh(t)
	if preview := pt.Preview(t); preview.ChangeSummary["update"] != 1 {
		t.Errorf("preview after outside edit = %v, want 1 update", preview.ChangeSummary)
	}
	pt.Up(t)

	pt.Destroy(t)
	if n := api.Count("DELETE", "/v2/dashboards/"); n != 1 {
		t.Errorf("DELETE count = %d, want 1", n)
	}
}

func TestApiKeySecretAndRotate(t *testing.T) {
	api := fakeapi.New(t)
	pt := newProgram(t, api, `
config:
  rotate:
    type: string
resources:
  ci:
    type: maple:ApiKey
    properties:
      name: ci
      scopes: ["dashboards:read"]
      rotate: ${rotate}
outputs:
  id: ${ci.id}
  secret: ${ci.secret}
`)
	pt.SetConfig(t, "rotate", "2026-10")
	up := pt.Up(t)
	firstID, _ := up.Outputs["id"].Value.(string)
	if !up.Outputs["secret"].Secret {
		t.Error("api key secret output is not a Pulumi secret")
	}

	// Pulumi cannot change an ID during an update, so a new keeper replaces the
	// key: the new one is created, then the old one revoked.
	pt.SetConfig(t, "rotate", "2026-11")
	up = pt.Up(t)
	if (*up.Summary.ResourceChanges)["replace"] != 1 {
		t.Errorf("changes = %v, want 1 replace", up.Summary.ResourceChanges)
	}
	if replaced, _ := up.Outputs["id"].Value.(string); replaced == firstID {
		t.Errorf("replacement kept ID %s", firstID)
	}
	if n := api.Count("DELETE", "/v2/api_keys/"+firstID); n != 1 {
		t.Errorf("old key revoked %d times, want 1", n)
	}
	// The new key is read back, not mistaken for a revoked one.
	pt.Refresh(t)
	assertpreview.HasNoChanges(t, pt.Preview(t))
}

func TestUptimeMonitorTypeChangeReplaces(t *testing.T) {
	api := fakeapi.New(t)
	program := func(config string) string {
		return `
resources:
  site:
    type: maple:UptimeMonitor
    properties:
      name: Website
      config:
        fn::toJSON:
` + config + `
outputs:
  kind: ${site.kind}
  heartbeatUrl: ${site.heartbeatUrl}
`
	}
	pt := newProgram(t, api, program(`          type: http
          url: https://example.com/health
          assertions: [{ type: status, op: is, value: 2xx }]`))
	if kind := pt.Up(t).Outputs["kind"].Value; kind != "http" {
		t.Errorf("kind = %v, want http", kind)
	}
	// The API adds method = "GET" to the config; that is not drift.
	assertpreview.HasNoChanges(t, pt.Preview(t))

	pt.WritePulumiYaml(t, "name: maple-test\nruntime: yaml\n"+program(`          type: heartbeat
          schedule: { type: interval, seconds: 300 }
          grace_seconds: 60`))
	up := pt.Up(t)
	if (*up.Summary.ResourceChanges)["replace"] != 1 {
		t.Errorf("changes = %v, want 1 replace", up.Summary.ResourceChanges)
	}
	if url, _ := up.Outputs["heartbeatUrl"].Value.(string); !strings.Contains(url, "/hb/upm_") {
		t.Errorf("heartbeatUrl = %q", url)
	}
}

func TestScrapeTargetCredentialsAreWriteOnly(t *testing.T) {
	api := fakeapi.New(t)
	pt := newProgram(t, api, `
resources:
  node:
    type: maple:ScrapeTarget
    properties:
      name: node-exporter
      url: https://metrics.example.com/metrics
      authType: bearer
      authCredentials:
        fn::secret: token-123
outputs:
  interval: ${node.scrapeIntervalSeconds}
`)
	up := pt.Up(t)
	if got := up.Outputs["interval"].Value; got != float64(15) {
		t.Errorf("scrapeIntervalSeconds = %v, want the API default 15", got)
	}
	// The API never returns the credentials, so state must keep the configured value.
	pt.Refresh(t)
	assertpreview.HasNoChanges(t, pt.Preview(t))
}

// The bridge returns data source results without secret markers (pulumi/pulumi#12710),
// so this checks the value only; the README tells users to wrap it in a secret.
func TestIngestKeysInvoke(t *testing.T) {
	api := fakeapi.New(t)
	pt := newProgram(t, api, `
variables:
  keys:
    fn::invoke:
      function: maple:getIngestKeys
      arguments: {}
outputs:
  privateKey: ${keys.privateKey}
`)
	if got := pt.Up(t).Outputs["privateKey"].Value; got != "maple_sk_test" {
		t.Errorf("privateKey = %v", got)
	}
}

// An existing dashboard imported with the same program has no diff: the API's
// canonical widgets JSON matches fn::toJSON output.
func TestDashboardImport(t *testing.T) {
	api := fakeapi.New(t)
	source := newProgram(t, api, dashboardProgram)
	source.SetConfig(t, "dashName", "Operations")
	id, _ := source.Up(t).Outputs["id"].Value.(string)

	// YAML takes only a literal import ID, so it is written into the program.
	pt := newProgram(t, api, strings.Replace(dashboardProgram, "    type: maple:Dashboard\n",
		"    type: maple:Dashboard\n    options:\n      import: "+id+"\n", 1))
	pt.SetConfig(t, "dashName", "Operations")
	up := pt.Up(t)
	if (*up.Summary.ResourceChanges)["import"] != 1 {
		t.Errorf("changes = %v, want 1 import", up.Summary.ResourceChanges)
	}
	assertpreview.HasNoChanges(t, pt.Preview(t))
}
