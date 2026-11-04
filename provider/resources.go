// Package maple bridges the Maple Terraform provider into a Pulumi provider.
package maple

import (
	_ "embed"
	"path"

	tfmaple "github.com/MapleTechLabs/terraform-provider-maple/provider"
	pfbridge "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfbridge"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge/tokens"

	"github.com/MapleTechLabs/pulumi-maple/provider/pkg/version"
)

const (
	mainPkg = "maple"
	mainMod = "index"
)

//go:embed cmd/pulumi-resource-maple/bridge-metadata.json
var metadata []byte

// Provider maps every resource and data source of the Terraform provider into the
// maple:index module. Sensitive Terraform attributes become Pulumi secrets.
// WithRotateReplaces: Pulumi cannot take the new ID an in-place key roll returns.
func Provider() tfbridge.ProviderInfo {
	prov := tfbridge.ProviderInfo{
		P:                 pfbridge.ShimProvider(tfmaple.New(version.Version, tfmaple.WithRotateReplaces())()),
		Name:              "maple",
		DisplayName:       "Maple",
		Version:           version.Version,
		Publisher:         "Maple",
		LogoURL:           "https://maple.dev/favicon.png",
		PluginDownloadURL: "github://api.github.com/MapleTechLabs/pulumi-maple",
		Description:       "Manage Maple dashboards, alert rules, alert destinations, uptime monitors, scrape targets and API keys.",
		Keywords:          []string{"maple", "observability", "opentelemetry", "category/monitoring"},
		License:           "MIT",
		Homepage:          "https://maple.dev/docs/iac/pulumi",
		Repository:        "https://github.com/MapleTechLabs/pulumi-maple",
		GitHubOrg:         "MapleTechLabs",
		MetadataInfo:      tfbridge.NewProviderMetadata(metadata),
		Config: map[string]*tfbridge.SchemaInfo{
			"api_key": {Secret: tfbridge.True()},
		},
		JavaScript: &tfbridge.JavaScriptInfo{
			PackageName:          "@maple-dev/pulumi",
			RespectSchemaVersion: true,
		},
		Golang: &tfbridge.GolangInfo{
			ImportBasePath: path.Join(
				"github.com/MapleTechLabs/pulumi-maple/sdk/",
				tfbridge.GetModuleMajorVersion(version.Version),
				"go",
				mainPkg,
			),
			GenerateResourceContainerTypes: true,
			GenerateExtraInputTypes:        true,
			RespectSchemaVersion:           true,
		},
	}

	prov.MustComputeTokens(tokens.SingleModule("maple_", mainMod, tokens.MakeStandard(mainPkg)))
	prov.MustApplyAutoAliases()
	return prov
}
