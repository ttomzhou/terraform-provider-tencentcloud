---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_throttling_apis"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_throttling_apis"
description: |-
  Use this data source to query api gateway throttling apis.
---

# tencentcloud_api_gateway_throttling_apis

Use this data source to query api gateway throttling apis.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_throttling_api" "service" {
  service_id       = "service-4r4xrrz4"
  strategy         = "400"
  environment_name = "test"
  api_ids          = ["api-lukm33yk"]
}

data "tencentcloud_api_gateway_throttling_apis" "id" {
  service_id = tencentcloud_api_gateway_throttling_api.service.service_id
}

data "tencentcloud_api_gateway_throttling_apis" "foo" {
  service_id        = tencentcloud_api_gateway_throttling_api.service.service_id
  environment_names = ["release", "test"]
}
```

## Argument Reference

The following arguments are supported:

* `environment_names` - (Optional) Environment list.
* `result_output_file` - (Optional) Used to save results.
* `service_id` - (Optional) Unique service ID of API.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of policies bound to API. Each element contains the following attributes:
  * `api_environment_strategies` - List of throttling policies bound to API.
    * `api_id` - Unique API ID.
    * `api_name` - Custom API name.
    * `method` - API method.
    * `path` - API path.
    * `strategy_list` - Environment throttling information.
      * `environment_name` - Environment name.
      * `quota` - Throttling value.
  * `service_id` - Unique service ID of API.


