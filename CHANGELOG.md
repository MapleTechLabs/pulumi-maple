# pulumi-maple

## 0.2.0

### Minor Changes

- f806890: `AlertRule.groupBies` is now `groupBy`, matching the Terraform `group_by` attribute. Rename it in
  programs that set it. The Go SDK now asks Pulumi for the plugin of its own version: `sdk/v0.1.0`
  asked for `0.0.1-dev`, which has no download, so Go programs need this release. The provider logo
  in the Pulumi schema now resolves.

### Patch Changes

- 2f1131a: API calls identify as `pulumi-maple/<version>` instead of the Terraform provider. The schema logo
  is the Maple wordmark SVG, and `docs/_index.md` is the overview page for the Pulumi Registry.

## 0.1.0

### Minor Changes

- a843a6e: First release: Maple dashboards, alert rules, alert destinations, uptime monitors, scrape targets, API keys and the ingest keys lookup, bridged from terraform-provider-maple, with TypeScript and Go SDKs.
