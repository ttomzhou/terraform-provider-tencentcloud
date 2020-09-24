---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_ip_strategies"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_ip_strategies"
description: |-
  Use this data source to query api gateway ip strategy.
---

# tencentcloud_api_gateway_ip_strategies

Use this data source to query api gateway ip strategy.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_ip_strategy" "test" {
  service_id    = "service-ohxqslqe"
  strategy_name = "tf_test"
  strategy_type = "BLACK"
  strategy_data = "9.9.9.9"
}

data "tencentcloud_api_gateway_ip_strategies" "id" {
  service_id = tencentcloud_api_gateway_ip_strategy.test.service_id
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required) ID of the service to be queried.
* `result_output_file` - (Optional) Used to save results.
* `strategy_name` - (Optional) Name of ip policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of strategy. Each element contains the following attributes:
  * `attach_list` - List of bound api details.
    * `api_business_type` - The type of oauth api. This field is valid when the authType is `OAUTH`, and the values are `NORMAL` (business api) and `OAUTH` (authorization api).
    * `api_desc` - The api interface description.
    * `api_id` - The api id.
    * `api_name` - The name of the api interface.
    * `api_type` - The api type. The values are `NORMAL` (common api) and `TSF` (microservice api).
    * `auth_relation_api_id` - The unique id of the associated authorization api, which takes effect when the authType is `OAUTH` and `ApiBusinessType` is normal. Identifies the unique id of the oauth2.0 authorization api bound to the business api.
    * `auth_type` - The api authentication type. The value is `SECRET` (key pair authentication), `NONE` (no authentication), and `OAUTH`.
    * `create_time` - Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
    * `method` - The api request method.
    * `modify_time` - Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
    * `oauth_config` - OAUTH configuration information. It takes effect when authType is `OAUTH`.
    * `path` - The api path.
    * `protocol` - The api protocol.
    * `relation_business_api_ids` - List of business api associated with authorized api.
    * `service_id` - The service id.
    * `tags` - The label information associated with the api.
    * `uniq_vpc_id` - The vpc unique id.
    * `vpc_id` - The vpc id.
  * `bind_api_total_count` - The number of api bound to the strategy.
  * `create_time` - Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
  * `ip_list` - The list of ip.
  * `modify_time` - Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
  * `service_id` - The service id.
  * `strategy_id` - ID of the strategy.
  * `strategy_name` - Name of the strategy.
  * `strategy_type` - Type of the strategy. Valid values: `WHITE` (white list) and `BLACK` (black list).


