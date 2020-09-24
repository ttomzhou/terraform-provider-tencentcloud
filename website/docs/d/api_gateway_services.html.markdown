---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_services"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_services"
description: |-
  Use this data source to query api gateway services.
---

# tencentcloud_api_gateway_services

Use this data source to query api gateway services.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "niceservice"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

data "tencentcloud_api_gateway_services" "name" {
  service_name = tencentcloud_api_gateway_service.service.service_name
}

data "tencentcloud_api_gateway_services" "ids" {
  service_id = tencentcloud_api_gateway_service.service.id
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional) Used to save results.
* `service_id` - (Optional) Service id for query.
* `service_name` - (Optional) Service name for query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of services. Each element contains the following attributes:
  * `api_list` - A list of apis. Each element contains the following attributes:
    * `api_desc` - Description of the api.
    * `api_id` - ID of the api.
    * `api_name` - Name of the api.
    * `method` - Method of the api.
    * `path` - Path of the api.
  * `create_time` - Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
  * `exclusive_set_name` - Self-deployed cluster name, which is used to specify the self-deployed cluster where the service is to be created.
  * `inner_http_port` - Port number for http access over private network.
  * `inner_https_port` - Port number for https access over private network.
  * `internal_sub_domain` - Private network access subdomain name.
  * `ip_version` - IP version number. Valid values: `IPv4` (default value), `IPv6`.
  * `modify_time` - Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
  * `net_type` - Network type list, which is used to specify the supported network types. `INNER` indicates access over private network, and `OUTER` indicates access over public network.
  * `outer_sub_domain` - Public network access subdomain name.
  * `protocol` - Service frontend request type, such as `http`, `https`, and `http&https`.
  * `service_desc` - Custom service description.
  * `service_id` - Custom service id.
  * `service_name` - Custom service name.
  * `usage_plan_list` - A list of attach usage plans. Each element contains the following attributes:
    * `api_id` - ID of the api.
    * `bind_type` - Binding type.
    * `usage_plan_id` - ID of the usage plan.
    * `usage_plan_name` - Name of the usage plan.


