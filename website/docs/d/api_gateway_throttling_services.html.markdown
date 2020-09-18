---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_throttling_services"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_throttling_services"
description: |-
  Use this data source to query api gateway throttling services.
---

# tencentcloud_api_gateway_throttling_services

Use this data source to query api gateway throttling services.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_throttling_service" "service" {
  service_id        = "service-4r4xrrz4"
  strategy          = "400"
  environment_names = ["release"]
}

data "tencentcloud_api_gateway_throttling_services" "id" {
  service_id = tencentcloud_api_gateway_throttling_service.service.service_id
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional) Used to save results.
* `service_id` - (Optional) Service ID for query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of Throttling policy. Each element contains the following attributes:
  * `environments` - A list of Throttling policy.
    * `environment_name` - Environment name.
    * `status` - Release status.
    * `strategy` - Throttling value.
    * `url` - Access service environment URL.
    * `version_name` - Published version number.
  * `service_id` - Service ID for query.


