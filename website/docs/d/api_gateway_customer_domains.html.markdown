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
resource "tencentcloud_api_gateway_custom_domain" "foo" {
  service_id         = "service-ohxqslqe"
  sub_domain         = "tic-test.dnsv1.com"
  protocol           = "http"
  net_type           = "OUTER"
  is_default_mapping = "false"
  default_domain     = "service-ohxqslqe-1259649581.gz.apigw.tencentcs.com"
  path_mappings      = ["/good#test", "/root#release"]
}

data "tencentcloud_api_gateway_customer_domains" "id" {
  service_id = tencentcloud_api_gateway_custom_domain.foo.service_id
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


