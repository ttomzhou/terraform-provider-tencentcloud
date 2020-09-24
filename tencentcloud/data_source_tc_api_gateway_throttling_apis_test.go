package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testAPIGatewaythrottlingApiDataSourceName = "data.tencentcloud_api_gateway_throttling_apis"

func TestAccTencentAPIGatewayThrottlingApisDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckThrottlingAPIDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTestAccTencentAPIGatewayThrottlingApis(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckThrottlingAPIExists("tencentcloud_api_gateway_throttling_api.service"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".id", "list.#"),
					resource.TestCheckResourceAttr(testAPIGatewaythrottlingApiDataSourceName+".id", "list.0.service_id", "service-4r4xrrz4"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".id", "list.0.api_environment_strategies.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".id", "list.0.api_environment_strategies.0.api_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".id", "list.0.api_environment_strategies.0.api_name"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".id", "list.0.api_environment_strategies.0.path"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".id", "list.0.api_environment_strategies.0.method"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".id", "list.0.api_environment_strategies.0.strategy_list.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".id", "list.0.api_environment_strategies.0.strategy_list.0.environment_name"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".id", "list.0.api_environment_strategies.0.strategy_list.0.quota"),

					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".foo", "list.#"),
					resource.TestCheckResourceAttr(testAPIGatewaythrottlingApiDataSourceName+".foo", "list.0.service_id", "service-4r4xrrz4"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".foo", "list.0.api_environment_strategies.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".foo", "list.0.api_environment_strategies.0.api_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".foo", "list.0.api_environment_strategies.0.api_name"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".foo", "list.0.api_environment_strategies.0.path"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".foo", "list.0.api_environment_strategies.0.method"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".foo", "list.0.api_environment_strategies.0.strategy_list.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".foo", "list.0.api_environment_strategies.0.strategy_list.0.environment_name"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".foo", "list.0.api_environment_strategies.0.strategy_list.0.quota"),
				),
			},
		},
	})
}

func testAccTestAccTencentAPIGatewayThrottlingApis() string {
	return `
resource "tencentcloud_api_gateway_throttling_api" "service" {
	service_id       = "service-4r4xrrz4"
	strategy         = "400"
	environment_name = "test"
	api_ids          = ["api-lukm33yk"]
}

data "tencentcloud_api_gateway_throttling_apis" "id" {
    service_id = "service-4r4xrrz4"
}

data "tencentcloud_api_gateway_throttling_apis" "foo" {
	service_id        = "service-4r4xrrz4"
	environment_names = ["release", "test"]
}
`
}
