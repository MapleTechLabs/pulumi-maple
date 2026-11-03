# Pulumi provider for Maple

Manage [Maple](https://maple.dev) dashboards, alert rules, alert destinations, uptime monitors,
scrape targets and API keys from Pulumi programs, through the public
[v2 API](https://api.maple.dev/v2/docs).

| Language   | Package                                                       |
| ---------- | ------------------------------------------------------------- |
| TypeScript | `@maple-dev/pulumi`                                           |
| Python     | `pulumi_maple`                                                |
| Go         | `github.com/MapleTechLabs/pulumi-maple/sdk/go/maple`          |

```typescript
import * as maple from "@maple-dev/pulumi"

new maple.Dashboard("ops", {
    name: "Operations",
    widgets: JSON.stringify([{ id: "w1", visualization: "chart", layout: { x: 0, y: 0, w: 6, h: 4 } }]),
})
```

Configure the provider with `pulumi config set maple:apiKey --secret maple_ak_…` or the
`MAPLE_API_KEY` environment variable (`maple:apiUrl` / `MAPLE_API_URL` for another region or a
local stack). More in [`examples/`](examples).

## How it works

This is the [Maple Terraform provider](https://github.com/MapleTechLabs/terraform-provider-maple)
wrapped with the Pulumi Terraform bridge. Resources, validation and API calls live there; this
repository maps them to Pulumi and generates the SDKs. A new resource in the Terraform provider
reaches Pulumi by bumping it and running `make generate`.

Differences from Terraform:

- **`ApiKey.rotate` replaces the key.** Terraform rolls the key in place, which changes its ID during
  an update. Pulumi cannot change an ID on update, so here a new key is created first and the old
  one revoked afterwards.
- **`getIngestKeys` returns plain values.** The bridge strips secret markers from function results
  ([pulumi/pulumi#12710](https://github.com/pulumi/pulumi/issues/12710)); wrap the keys with
  `pulumi.secret(...)` before exporting them.
- **JSON attributes are strings** (dashboard `widgets`, uptime `config`, alert `thresholds`): pass
  `JSON.stringify(...)` / `json.dumps(...)`. Defaults the API adds are not reported as drift.

## Develop

Needs Go, the Pulumi CLI, and Node and Python 3.10+ for the SDK builds.

```bash
make test       # mapping and secrets checks, then Pulumi YAML programs against an in-process fake API
make generate   # schema.json and the Node, Python and Go SDKs; commit the result
make provider   # bin/pulumi-resource-maple
```

The lifecycle tests run Pulumi YAML programs, so they need no SDK build. They use
`fakeapi` from the Terraform provider, the same fake v2 API its own tests run against.

To try a local build in a program:

```bash
make provider
pulumi plugin rm resource maple --yes
export PATH="$PWD/bin:$PATH"
```

## Releasing

Push a `vX.Y.Z` tag. The [release workflow](.github/workflows/release.yml) publishes the provider
plugin binaries to the GitHub release (Pulumi downloads them from there), tags `sdk/vX.Y.Z` for
the Go SDK, and publishes the npm package (`NPM_TOKEN`) and the PyPI package (trusted publishing,
`pypi` environment).
