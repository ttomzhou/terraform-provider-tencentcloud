---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_ip_strategies"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_ip_strategies"
description: |-
  Use this data source to query api gateway IP strategy.
---

# tencentcloud_api_gateway_ip_strategies

Use this data source to query api gateway IP strategy.

## Example Usage

```hcl
data "tencentcloud_api_gateway_ip_strategies" "id" {
  service_id = "service-ohxqslqe"
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required) ID of the service to be queried.
* `result_output_file` - (Optional) Used to save results.
* `strategy_name` - (Optional) The name of ip policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of strategy. Each element contains the following attributes:
  * `attach_list` - List of bound API details.
    * `api_business_type` - The type of OAUTH API. This field is valid when the AuthType is OAUTH, and the values are NORMAL (business API) and OAUTH (authorization API).
    * `api_desc` - The api interface description.
    * `api_id` - The api id.
    * `api_name` - The name of the API interface.
    * `api_type` - The api type. The values are NORMAL (common API) and TSF (microservice API).
    * `auth_relation_api_id` - The unique ID of the associated authorization API, which takes effect when the AuthType is OAUTH and ApiBusinessType is NORMAL. Identifies the unique ID of the oauth2.0 authorization API bound to the business API.
    * `auth_type` - The api authentication type. The value is SECRET (key pair authentication), NONE (no authentication), and OAUTH.
    * `create_time` - Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
    * `method` - The api request method.
    * `modify_time` - Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
    * `oauth_config` - OAUTH configuration information. It takes effect when AuthType is OAUTH.
    * `path` - The api path.
    * `protocol` - The api protocol.
    * `relation_business_api_ids` - List of business APIs associated with authorized APIs.
    * `service_id` - The service id.
    * `tags` - The label information associated with the api.
    * `uniq_vpc_id` - The vpc unique id.
    * `vpc_id` - The vpc id.
  * `bind_api_total_count` - The number of APIs bound to the strategy.
  * `create_time` - Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
  * `ip_list` - The list of ip.
  * `modify_time` - Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
  * `service_id` - The service id.
  * `strategy_id` - ID of the strategy.
  * `strategy_name` - Name of the strategy.
  * `strategy_type` - Type of the strategy. Valid values: WHITE (white list) and BLACK (black list).


