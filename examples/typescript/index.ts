import * as pulumi from "@pulumi/pulumi"
import * as maple from "@maple-dev/pulumi"

const oncall = new maple.AlertDestination("oncall", {
    name: "On-call webhook",
    type: "webhook",
    url: new pulumi.Config().requireSecret("webhookUrl"),
})

new maple.AlertRule("checkout-errors", {
    name: "Checkout error rate",
    signalType: "error_rate",
    comparator: "gt",
    threshold: 0.05, // error rates are 0 to 1 ratios
    windowMinutes: 5,
    serviceNames: ["checkout"],
    destinationIds: [oncall.id],
})

new maple.Dashboard("ops", {
    name: "Operations",
    tags: ["production"],
    widgets: JSON.stringify([{ id: "w1", visualization: "chart", layout: { x: 0, y: 0, w: 6, h: 4 } }]),
})

const ci = new maple.ApiKey("ci", { name: "ci", scopes: ["dashboards:read"] })

const keys = maple.getIngestKeysOutput({})

export const ciKey = ci.secret
export const ingestKey = pulumi.secret(keys.privateKey)
