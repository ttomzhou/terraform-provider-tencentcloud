---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_customer_domains"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_customer_domains"
description: |-
  Use this data source to query api gateway domain list.
---

# tencentcloud_api_gateway_customer_domains

Use this data source to query api gateway domain list.

## Example Usage

```hcl
data "tencentcloud_api_gateway_customer_domains" "id" {
  service_id = "service-ohxqslqe"
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required) The id of service.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - Service custom domain name list.
  * `certificate_id` - The id of certificate.
  * `domain_name` - Domain name.
  * `is_default_mapping` - Whether to use default path mapping, true means to use default path mapping; if false, means to use custom path mapping.
  * `net_type` - Network type, valid value: INNER or OUTER.
  * `path_mappings` - Domain name mapping path and environment list.
    * `environment` - Release environment, optional values are [test, prepub, release].
    * `path` - The domain mapping path.
  * `protocol` - Custom domain name agreement type.
  * `status` - Domain name resolution status. true means normal parsing, false means parsing failed.


