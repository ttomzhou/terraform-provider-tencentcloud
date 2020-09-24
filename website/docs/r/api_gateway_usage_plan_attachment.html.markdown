---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_usage_plan_attachment"
sidebar_current: "docs-tencentcloud-resource-api_gateway_usage_plan_attachment"
description: |-
  Use this resource to attach api gateway usage plan to service.
---

# tencentcloud_api_gateway_usage_plan_attachment

Use this resource to attach api gateway usage plan to service.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_usage_plan" "plan" {
  usage_plan_name         = "my_plan"
  usage_plan_desc         = "nice plan"
  max_request_num         = 100
  max_request_num_pre_sec = 10
}

resource "tencentcloud_api_gateway_service" "service" {
  service_name = "niceservice"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_usage_plan_attachment" "attach_service" {
  usage_plan_id = tencentcloud_api_gateway_usage_plan.plan.id
  service_id    = tencentcloud_api_gateway_service.service.id
  environment   = "test"
  bind_type     = "SERVICE"
}
```

## Argument Reference

The following arguments are supported:

* `environment` - (Required, ForceNew) Environment to be bound `test`,`prepub` or `release`.
* `service_id` - (Required, ForceNew) ID of the service.
* `usage_plan_id` - (Required, ForceNew) ID of the usage plan.
* `api_id` - (Optional, ForceNew) ID of the api. This parameter will be required when `bind_type` is `API`.
* `bind_type` - (Optional, ForceNew) Binding type. Valid values: `API`, `SERVICE` (default value).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

api gateway usage plan attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_usage_plan_attachment.attach_service '{"api_id":"","bind_type":"SERVICE","environment":"test","service_id":"service-pkegyqmc","usage_plan_id":"usagePlan-26t0l0w3"}'
```

