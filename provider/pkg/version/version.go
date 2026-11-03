// Package version holds the provider version, set at build time with -ldflags.
package version

// Version is the provider's semver. Release builds override it with
// -X github.com/MapleTechLabs/pulumi-maple/provider/pkg/version.Version=<v>.
var Version = "0.0.1-dev"
