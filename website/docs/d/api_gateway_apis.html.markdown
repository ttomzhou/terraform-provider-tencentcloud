---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_apis"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_apis"
description: |-
  Use this data source to query api gateway apis.
---

# tencentcloud_api_gateway_apis

Use this data source to query api gateway apis.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "ck"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_api" "api" {
  service_id               = tencentcloud_api_gateway_service.service.id
  api_name                 = "hello"
  api_desc                 = "my hello api"
  auth_type                = "NONE"
  protocol                 = "HTTP"
  enable_cors              = true
  request_config_path      = "/user/info"
  request_config_method    = "GET"
  service_config_type      = "HTTP"
  service_config_timeout   = 15
  service_config_url       = "http://www.qq.com"
  service_config_path      = "/user"
  service_config_method    = "GET"
  response_type            = "HTML"
  response_success_example = "success"
  response_fail_example    = "fail"
}

data "tencentcloud_api_gateway_apis" "id" {
  service_id = tencentcloud_api_gateway_service.service.id
  api_id     = tencentcloud_api_gateway_api.api.id
}

data "tencentcloud_api_gateway_apis" "name" {
  service_id = tencentcloud_api_gateway_service.service.id
  api_name   = tencentcloud_api_gateway_api.api.api_name
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required) Service id for query.
* `api_id` - (Optional) Created api id.
* `api_name` - (Optional) Custom api name.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of apis. Each element contains the following attributes:
  * `api_desc` - Custom api description.
  * `api_name` - Custom api name.
  * `auth_type` - API authentication type. Valid values: `SECRET` (key pair authentication),`NONE` (no authentication).
  * `create_time` - Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
  * `enable_cors` - Whether to enable CORS.
  * `modify_time` - Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
  * `protocol` - API frontend request type, such as `HTTP`,`WEBSOCKET`.
  * `request_config_method` - Request frontend method configuration. Like `GET`,`POST`,`PUT`,`DELETE`,`HEAD`,`ANY`.
  * `request_config_path` - Request frontend path configuration. Like `/user/getinfo`.
  * `request_parameters` - Frontend request parameters.
    * `default_value` - Parameter default value.
    * `desc` - Parameter description.
    * `name` - Parameter name.
    * `position` - Parameter location.
    * `required` - If this parameter required.
    * `type` - Parameter type.
  * `response_error_codes` - Custom error code configuration. Must keep at least one after set.
    * `code` - Custom response configuration error code.
    * `converted_code` - Custom error code conversion.
    * `desc` - Parameter description.
    * `msg` - Custom response configuration error message.
    * `need_convert` - Whether to enable error code conversion. Default value: `false`.
  * `response_fail_example` - Response failure sample of custom response configuration.
  * `response_success_example` - Successful response sample of custom response configuration.
  * `response_type` - Return type.
  * `service_config_method` - API backend service request method, such as `GET`. If `service_config_type` is `HTTP`, this parameter will be required. The frontend `request_config_method` and backend method `service_config_method` can be different.
  * `service_config_mock_return_message` - Returned information of API backend mocking. This parameter is required when `service_config_type`  is `MOCK`.
  * `service_config_path` - API backend service path, such as /path. If `service_config_type` is `HTTP`, this parameter will be required. The frontend `request_config_path` and backend path `service_config_path` can be different.
  * `service_config_product` - Backend type. This parameter takes effect when vpc is enabled. Currently, only `clb` is supported.
  * `service_config_scf_function_name` - SCF function name. This parameter takes effect when `service_config_type` is `SCF`.
  * `service_config_scf_function_namespace` - SCF function namespace. This parameter takes effect when  `service_config_type` is `SCF`.
  * `service_config_scf_function_qualifier` - SCF function version. This parameter takes effect when `service_config_type`  is `SCF`.
  * `service_config_timeout` - API backend service timeout period in seconds.
  * `service_config_type` - API backend service type.
  * `service_config_url` - API backend service url. This parameter is required when `service_config_type` is `HTTP`.
  * `service_config_vpc_id` - Unique vpc id.
  * `service_id` - Which service this api belongs.Refer to resource `tencentcloud_api_gateway_service`.


