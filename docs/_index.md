---
title: Maple
meta_desc: Use the Maple provider for Pulumi to manage Maple dashboards, alert rules, alert destinations, uptime monitors, scrape targets and API keys.
layout: package
---

The Maple provider for Pulumi lets you manage [Maple](https://maple.dev) observability configuration (dashboards, alert rules and their notification destinations, uptime and heartbeat monitors, scrape targets and API keys) as part of your Pulumi programs. It talks to Maple's public [v2 API](https://maple.dev/docs/reference/api) and is built from the [Maple Terraform provider](https://registry.terraform.io/providers/MapleTechLabs/maple/latest), so both have the same resources and fields.

## Installation

{{< chooser language "typescript,python,go,csharp,java,yaml,hcl" >}}
{{% choosable language typescript %}}

```bash
npm install @maple-dev/pulumi
```

{{% /choosable %}}
{{% choosable language python %}}

The provider publishes no Python package. Generate the SDK in your project:

```bash
pulumi package add maple
```

{{% /choosable %}}
{{% choosable language go %}}

```bash
go get github.com/MapleTechLabs/pulumi-maple/sdk/go/maple
```

{{% /choosable %}}
{{% choosable language csharp %}}

The provider publishes no .NET package. Generate the SDK in your project:

```bash
pulumi package add maple
```

{{% /choosable %}}
{{% choosable language java %}}

The provider publishes no Java package. Generate the SDK in your project:

```bash
pulumi package add maple
```

{{% /choosable %}}
{{% choosable language yaml %}}

```bash
pulumi package add maple
```

{{% /choosable %}}
{{% choosable language hcl %}}

Declare the provider in your program, then run `pulumi install`:

```hcl
terraform {
  required_providers {
    maple = {
      source  = "mapletechlabs/maple"
      version = "~> 0.1"
    }
  }
}
```

```bash
pulumi install
```

{{% /choosable %}}
{{< /chooser >}}

Pulumi downloads the provider plugin from the GitHub release on first use.

## Example Usage

An alert rule that pages a PagerDuty service when the checkout service's error rate stays above 5% for five minutes:

{{< chooser language "typescript,go" >}}
{{% choosable language typescript %}}

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as maple from "@maple-dev/pulumi";

const config = new pulumi.Config();

const oncall = new maple.AlertDestination("oncall", {
    type: "pagerduty",
    name: "Primary on-call",
    integrationKey: config.requireSecret("pagerdutyIntegrationKey"),
});

const checkoutErrors = new maple.AlertRule("checkout-errors", {
    name: "Checkout error rate",
    signalType: "error_rate",
    comparator: "gt",
    severity: "critical",
    threshold: 0.05,
    windowMinutes: 5,
    serviceNames: ["checkout"],
    destinationIds: [oncall.id],
});

export const alertRuleId = checkoutErrors.id;
```

{{% /choosable %}}
{{% choosable language go %}}

```go
package main

import (
	"github.com/MapleTechLabs/pulumi-maple/sdk/go/maple"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		oncall, err := maple.NewAlertDestination(ctx, "oncall", &maple.AlertDestinationArgs{
			Type:           pulumi.String("pagerduty"),
			Name:           pulumi.String("Primary on-call"),
			IntegrationKey: cfg.RequireSecret("pagerdutyIntegrationKey"),
		})
		if err != nil {
			return err
		}
		checkoutErrors, err := maple.NewAlertRule(ctx, "checkout-errors", &maple.AlertRuleArgs{
			Name:           pulumi.String("Checkout error rate"),
			SignalType:     pulumi.String("error_rate"),
			Comparator:     pulumi.String("gt"),
			Severity:       pulumi.String("critical"),
			Threshold:      pulumi.Float64(0.05),
			WindowMinutes:  pulumi.Int(5),
			ServiceNames:   pulumi.StringArray{pulumi.String("checkout")},
			DestinationIds: pulumi.StringArray{oncall.ID()},
		})
		if err != nil {
			return err
		}
		ctx.Export("alertRuleId", checkoutErrors.ID())
		return nil
	})
}
```

{{% /choosable %}}
{{< /chooser >}}

Configure it with a Maple API key from **Settings → Keys** and the PagerDuty integration key:

```bash
pulumi config set --secret maple:apiKey maple_ak_...
pulumi config set --secret pagerdutyIntegrationKey <integration-key>
pulumi up
```

## Configuration

- `apiKey` (Required, Secret): a Maple API key (`maple_ak_…`). May also be set with the `MAPLE_API_KEY` environment variable. Alert rules, alert destinations, API keys and ingest keys need a key created by an org admin; a key with scopes needs the matching `:write` scope for each resource type it manages, e.g. `dashboards:write`.
- `apiUrl` (Optional): base URL of the Maple API. May also be set with the `MAPLE_API_URL` environment variable. Defaults to `https://api.maple.dev`; organizations in the EU region use `https://api.eu.maple.dev`.

## Resources

- `maple.Dashboard`
- `maple.AlertDestination`: PagerDuty, webhook, Discord, Telegram, Google Chat, email and chat workspace channels
- `maple.AlertRule`
- `maple.UptimeMonitor`: HTTP, TCP, DNS and heartbeat monitors
- `maple.ScrapeTarget`
- `maple.ApiKey`
- `maple.getIngestKeys`: reads the organization's ingest keys

JSON fields (a dashboard's `widgets`, an uptime monitor's `config`, an alert rule's `thresholds`) take a JSON string in the API's wire shape, e.g. `JSON.stringify(...)`. Credentials and API key secrets are Pulumi secrets. `getIngestKeys` returns plain values, so wrap them with `pulumi.secret(...)` before exporting.

Every resource can be imported by its Maple ID, for example `pulumi import maple:index/dashboard:Dashboard ops dash_2xH8d...`.

The [Maple Pulumi guide](https://maple.dev/docs/iac/pulumi) covers import, JSON fields, secrets and API key rotation in detail.
