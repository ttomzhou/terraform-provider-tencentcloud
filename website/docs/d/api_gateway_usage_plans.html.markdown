---
subcategory: "API Gateway"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_usage_plans"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_usage_plans"
description: |-
  Use this data source to query api gateway usage plans.
---

# tencentcloud_api_gateway_usage_plans

Use this data source to query api gateway usage plans.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_usage_plan" "plan" {
  usage_plan_name         = "my_plan"
  usage_plan_desc         = "nice plan"
  max_request_num         = 100
  max_request_num_pre_sec = 10
}

data "tencentcloud_api_gateway_usage_plans" "name" {
  usage_plan_name = tencentcloud_api_gateway_usage_plan.plan.usage_plan_name
}

data "tencentcloud_api_gateway_usage_plans" "id" {
  usage_plan_id = tencentcloud_api_gateway_usage_plan.plan.id
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional) Used to save results.
* `usage_plan_id` - (Optional) ID of the usage plan to be queried.
* `usage_plan_name` - (Optional) Name of the usage plan to be queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of usage plans. Each element contains the following attributes:
  * `attach_list` - Attach service and api list.
    * `api_id` - The api id, This value is empty if attach service.
    * `api_name` - The api name, This value is empty if attach service.
    * `create_time` - Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
    * `environment` - The environment name.
    * `method` - The api method, This value is empty if attach service.
    * `modify_time` - Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
    * `path` - The api path, This value is empty if attach service.
    * `service_id` - The service id.
    * `service_name` - The service name.
  * `create_time` - Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
  * `max_request_num_pre_sec` - Limit of requests per second. Valid values: -1, [1,2000]. The default value is -1, which indicates no limit.
  * `max_request_num` - Total number of requests allowed. Valid values: -1, [1,99999999]. The default value is -1, which indicates no limit.
  * `modify_time` - Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
  * `usage_plan_desc` - Custom usage plan description.
  * `usage_plan_id` - ID of the usage plan.
  * `usage_plan_name` - Name of the usage plan.


