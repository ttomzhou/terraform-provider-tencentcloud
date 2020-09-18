---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_throttling_service"
sidebar_current: "docs-tencentcloud-resource-api_gateway_throttling_service"
description: |-
  Use this resource to create api gateway throttling service.
---

# tencentcloud_api_gateway_throttling_service

Use this resource to create api gateway throttling service.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_throttling_service" "service" {
  service_id        = "service-4r4xrrz4"
  strategy          = "400"
  environment_names = ["release"]
}
```

## Argument Reference

The following arguments are supported:

* `environment_names` - (Required) List of Environment names.
* `service_id` - (Required, ForceNew) Service ID for query.
* `strategy` - (Required) Throttling value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `environments` - A list of Throttling policy.
  * `environment_name` - Environment name.
  * `status` - Release status.
  * `strategy` - Throttling value.
  * `url` - Access service environment URL.
  * `version_name` - Published version number.


