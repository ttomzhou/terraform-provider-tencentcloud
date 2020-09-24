---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_throttling_api"
sidebar_current: "docs-tencentcloud-resource-api_gateway_throttling_api"
description: |-
  Use this resource to create api gateway throttling api.
---

# tencentcloud_api_gateway_throttling_api

Use this resource to create api gateway throttling api.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_throttling_api" "service" {
  service_id       = "service-4r4xrrz4"
  strategy         = "400"
  environment_name = "test"
  api_ids          = ["api-lukm33yk"]
}
```

## Argument Reference

The following arguments are supported:

* `api_ids` - (Required) List of API ID.
* `environment_name` - (Required) List of Environment names.
* `service_id` - (Required, ForceNew) Service ID for query.
* `strategy` - (Required) Throttling value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `api_environment_strategies` - List of throttling policies bound to API.
  * `api_id` - Unique API ID.
  * `api_name` - Custom API name.
  * `method` - API method.
  * `path` - API path.
  * `strategy_list` - Environment throttling information.
    * `environment_name` - Environment name.
    * `quota` - Throttling value.


