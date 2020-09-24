---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_api_keys"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_api_keys"
description: |-
  Use this data source to query api gateway access keys.
---

# tencentcloud_api_gateway_api_keys

Use this data source to query api gateway access keys.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_api_key" "test_cluster" {
  secret_name = "my_api_key"
  status      = "on"
}
resource "tencentcloud_api_gateway_api_key" "test_cluster2" {
  secret_name = "my_api_key"
  status      = "on"
}
data "tencentcloud_api_gateway_api_keys" "name" {
  secret_name = tencentcloud_api_gateway_api_key.test_cluster.secret_name
}

data "tencentcloud_api_gateway_api_keys" "id" {
  access_key_id = tencentcloud_api_gateway_api_key.test_cluster.access_key_id
}
```

## Argument Reference

The following arguments are supported:

* `access_key_id` - (Optional) Created API key ID, This field is exactly the same as ID.
* `result_output_file` - (Optional) Used to save results.
* `secret_name` - (Optional) Custom key name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of API keys. Each element contains the following attributes:
  * `access_key_id` - Created API key ID, This field is exactly the same as `api_key_id`.
  * `access_key_secret` - Created API key.
  * `api_key_id` - API key ID.
  * `create_time` - Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
  * `modify_time` - Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.


