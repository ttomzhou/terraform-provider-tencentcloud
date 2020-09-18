---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_usage_plan_environments"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_usage_plan_environments"
description: |-
  Used to query the environment list bound by the plan
---

# tencentcloud_api_gateway_usage_plan_environments

Used to query the environment list bound by the plan

## Example Usage

```hcl
resource "tencentcloud_api_gateway_usage_plan" "plan" {
  usage_plan_name         = "my_plan"
  usage_plan_desc         = "nice plan"
  max_request_num         = 100
  max_request_num_pre_sec = 10
}

resource "tencentcloud_api_gateway_usage_plan_attachment" "attach_service" {
  usage_plan_id = tencentcloud_api_gateway_usage_plan.plan.id
  service_id    = "service-ke4o2arm"
  environment   = "test"
  bind_type     = "SERVICE"
}

data "tencentcloud_api_gateway_usage_plan_environments" "environment_test" {
  usage_plan_id = tencentcloud_api_gateway_usage_plan_attachment.attach_service.usage_plan_id
  bind_type     = "SERVICE"
}
```

## Argument Reference

The following arguments are supported:

* `usage_plan_id` - (Required) ID of the usage plan to be queried.
* `bind_type` - (Optional) Binding type. Valid values: `API`, `SERVICE` (default value).
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of usage plan binding details. Each element contains the following attributes:
  * `api_id` - The api id, This value is empty if attach service.
  * `api_name` - The api name, This value is empty if attach service.
  * `create_time` - Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
  * `environment` - The environment name.
  * `method` - The api method, This value is empty if attach service.
  * `modify_time` - Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
  * `path` - The api path, This value is empty if attach service.
  * `service_id` - The service id.
  * `service_name` - The service name.


