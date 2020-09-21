---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_service"
sidebar_current: "docs-tencentcloud-resource-api_gateway_service"
description: |-
  Use this resource to create api gateway service.
---

# tencentcloud_api_gateway_service

Use this resource to create api gateway service.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "niceservice"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}
```

## Argument Reference

The following arguments are supported:

* `net_type` - (Required) Network type list, which is used to specify the supported network types. `INNER` indicates access over private network, and `OUTER` indicates access over public network.
* `protocol` - (Required) Service frontend request type, such as `http`, `https`, and `http&https`.
* `service_name` - (Required) Custom service name.
* `appid_type` - (Optional, ForceNew) User type, which is reserved and can be used by serverless users.
* `exclusive_set_name` - (Optional, ForceNew) Self-deployed cluster name, which is used to specify the self-deployed cluster where the service is to be created.
* `ip_version` - (Optional, ForceNew) IP version number. Valid values: `IPv4` (default value), `IPv6`.
* `service_desc` - (Optional) Custom service description.
* `set_server_name` - (Optional, ForceNew) Cluster name, which is reserved and used by the tsf serverless type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `api_list` - A list of apis. Each element contains the following attributes:
  * `api_desc` - Description of the api.
  * `api_id` - ID of the api.
  * `api_name` - Name of the api.
  * `method` - Method of the api.
  * `path` - Path of the api.
* `create_time` - Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
* `inner_http_port` - Port number for http access over private network.
* `inner_https_port` - Port number for https access over private network.
* `internal_sub_domain` - Private network access subdomain name.
* `modify_time` - Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
* `outer_sub_domain` - Public network access subdomain name.
* `service_id` - Service ID for query.
* `usage_plan_list` - A list of attach usage plans. Each element contains the following attributes:
  * `api_id` - ID of the api.
  * `bind_type` - Binding type.
  * `usage_plan_id` - ID of the usage plan.
  * `usage_plan_name` - Name of the usage plan.


## Import

api gateway service can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_service.service service-pg6ud8pa
```

