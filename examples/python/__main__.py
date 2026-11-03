import json

import pulumi
import pulumi_maple as maple

maple.Dashboard(
    "ops",
    name="Operations",
    tags=["production"],
    widgets=json.dumps([{"id": "w1", "visualization": "chart", "layout": {"x": 0, "y": 0, "w": 6, "h": 4}}]),
)

ci = maple.ApiKey("ci", name="ci", scopes=["dashboards:read"])
pulumi.export("ciKey", ci.secret)
