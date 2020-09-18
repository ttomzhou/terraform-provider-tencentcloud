package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testAPIGatewayDomainSourceName = "data.tencentcloud_api_gateway_customer_domains"

func TestAccTencentAPIGatewayCustomerDomain(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTestAccTencentAPIGatewayDomain(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAPIGatewayDomainSourceName+".id", "list.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewayDomainSourceName+".id", "list.0.domain_name"),
					resource.TestCheckResourceAttrSet(testAPIGatewayDomainSourceName+".id", "list.0.status"),
					resource.TestCheckResourceAttrSet(testAPIGatewayDomainSourceName+".id", "list.0.is_default_mapping"),
					resource.TestCheckResourceAttrSet(testAPIGatewayDomainSourceName+".id", "list.0.net_type"),
					resource.TestCheckResourceAttrSet(testAPIGatewayDomainSourceName+".id", "list.0.path_mappings.#"),
				),
			},
		},
	})
}

func testAccTestAccTencentAPIGatewayDomain() string {
	return `
		data "tencentcloud_api_gateway_customer_domains" "id" {
  			service_id = "service-ohxqslqe" 
		}
	`
}
