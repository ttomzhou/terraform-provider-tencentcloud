---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_ip_strategy"
sidebar_current: "docs-tencentcloud-resource-api_gateway_ip_strategy"
description: |-
  Use this resource to create api gateway ip strategy.
---

# tencentcloud_api_gateway_ip_strategy

Use this resource to create api gateway ip strategy.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_ip_strategy" "test" {
  service_id    = "service-ohxqslqe"
  strategy_name = "tf_test"
  strategy_type = "BLACK"
  strategy_data = "9.9.9.9"
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required, ForceNew) The ID of the API gateway service.
* `strategy_data` - (Required) IP address data.
* `strategy_name` - (Required, ForceNew) User defined strategy name.
* `strategy_type` - (Required, ForceNew) Blacklist or whitelist.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
* `strategy_id` - IP policy ID.


## Import

api gateway ip strategy can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_ip_strategy.test service-ohxqslqe#IPStrategy-q1lk8ud2
```

