// Command pulumi-resource-maple is the Pulumi resource provider plugin.
package main

import (
	"context"
	_ "embed"

	pfbridge "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfbridge"

	maple "github.com/MapleTechLabs/pulumi-maple/provider"
	"github.com/MapleTechLabs/pulumi-maple/provider/pkg/version"
)

//go:embed schema.json
var pulumiSchema []byte

func main() {
	schema, err := maple.WithSchemaVersion(pulumiSchema, version.Version)
	if err != nil {
		panic(err)
	}
	pfbridge.Main(context.Background(), "maple", maple.Provider(), pfbridge.ProviderMetadata{PackageSchema: schema})
}
