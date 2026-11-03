// Command pulumi-tfgen-maple generates the Pulumi schema and SDKs.
package main

import (
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfgen"

	maple "github.com/MapleTechLabs/pulumi-maple/provider"
)

func main() {
	tfgen.Main("maple", maple.Provider())
}
