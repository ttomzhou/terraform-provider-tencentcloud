package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testAPIGatewaythrottlingServiceDataSourceName = "data.tencentcloud_api_gateway_throttling_services"

func TestAccTencentAPIGatewayThrottlingServicesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckThrottlingServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTestAccTencentAPIGatewayThrottlingServices(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckThrottlingServiceExists("tencentcloud_api_gateway_throttling_service.service"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingServiceDataSourceName+".id", "list.#"),
					resource.TestCheckResourceAttr(testAPIGatewaythrottlingServiceDataSourceName+".id", "list.0.service_id", "service-4r4xrrz4"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingServiceDataSourceName+".id", "list.0.environments.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingServiceDataSourceName+".id", "list.0.environments.0.environment_name"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingServiceDataSourceName+".id", "list.0.environments.0.url"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingServiceDataSourceName+".id", "list.0.environments.0.status"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingServiceDataSourceName+".id", "list.0.environments.0.version_name"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingServiceDataSourceName+".id", "list.0.environments.0.strategy"),
				),
			},
		},
	})
}

func testAccTestAccTencentAPIGatewayThrottlingServices() string {
	return `

resource "tencentcloud_api_gateway_throttling_service" "service" {
	service_id = "service-4r4xrrz4"
	strategy = "400"
	environment_names = ["release"]
}

data "tencentcloud_api_gateway_throttling_services" "id" {
    service_id = "service-4r4xrrz4"
}

`
}
