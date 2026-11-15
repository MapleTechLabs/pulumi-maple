---
"@maple-dev/pulumi": minor
---

`AlertRule.groupBies` is now `groupBy`, matching the Terraform `group_by` attribute. Rename it in
programs that set it. The Go SDK now asks Pulumi for the plugin of its own version: `sdk/v0.1.0`
asked for `0.0.1-dev`, which has no download, so Go programs need this release. The provider logo
in the Pulumi schema now resolves.
