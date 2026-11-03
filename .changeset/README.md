# Changesets

Releases are versioned with [`@maple-dev/changesets`](https://github.com/MapleTechLabs/changesets).
In a PR that changes what users get (a new resource from a Terraform provider bump, a fix, a
breaking change), run `npx changeset`, pick the bump, and describe the change.

`.github/workflows/release.yml` keeps a "Version Packages" PR open while changesets are pending.
Merging it tags `vX.Y.Z` and publishes the provider plugin, the Go SDK tag, npm and PyPI.
