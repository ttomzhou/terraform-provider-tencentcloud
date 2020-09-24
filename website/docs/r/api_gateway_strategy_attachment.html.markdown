---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_strategy_attachment"
sidebar_current: "docs-tencentcloud-resource-api_gateway_strategy_attachment"
description: |-
  Use this resource to create api gateway ip strategy attachment.
---

# tencentcloud_api_gateway_strategy_attachment

Use this resource to create api gateway ip strategy attachment.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_strategy_attachment" "test" {
  service_id       = "service-ohxqslqe"
  strategy_id      = tencentcloud_api_gateway_ip_strategy.test.strategy_id
  environment_name = "release"
  bind_api_id      = "api-jbtpu758"
}
```

## Argument Reference

The following arguments are supported:

* `bind_api_id` - (Required, ForceNew) The API that needs to be bound.
* `environment_name` - (Required, ForceNew) The environment of the strategy association.
* `service_id` - (Required, ForceNew) The id of the API gateway service.
* `strategy_id` - (Required, ForceNew) The id of the API gateway strategy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

api gateway ip strategy attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_strategy_attachment.test service-ohxqslqe#IPStrategy-nbxqk56k#api-jbtpu758#release
```

