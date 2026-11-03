// Command pulumi-resource-maple is the Pulumi resource provider plugin.
package main

import (
	"context"
	_ "embed"

	pfbridge "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfbridge"

	maple "github.com/MapleTechLabs/pulumi-maple/provider"
)

//go:embed schema.json
var pulumiSchema []byte

func main() {
	pfbridge.Main(context.Background(), "maple", maple.Provider(), pfbridge.ProviderMetadata{PackageSchema: pulumiSchema})
}
