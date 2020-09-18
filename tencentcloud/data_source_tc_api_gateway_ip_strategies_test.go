package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testAPIGatewayIpStrategySourceName = "data.tencentcloud_api_gateway_ip_strategies"

// @TODO 添加检查IP策略是否存在，删除服务IP策略
func TestAccTencentAPIGatewayIpStrategyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccCheckAPIGate,
		Steps: []resource.TestStep{
			{
				Config: testAccTestAccTencentAPIGatewayIpStrategy(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".id", "list.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".id", "list.0.strategy_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".id", "list.0.strategy_name"),
					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".id", "list.0.bind_api_total_count"),
					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".id", "list.0.create_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".id", "list.0.attach_list.#"),

					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".name", "list.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".name", "list.0.strategy_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".name", "list.0.strategy_name"),
					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".name", "list.0.bind_api_total_count"),
					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".name", "list.0.create_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".name", "list.0.attach_list.#"),
				),
			},
		},
	})
}

func testAccTestAccTencentAPIGatewayIpStrategy() string {
	return `
		data "tencentcloud_api_gateway_ip_strategies" "id" {
  			service_id = "service-ohxqslqe" 
		}
		
		data "tencentcloud_api_gateway_ip_strategies" "name" {
			service_id 	  = "service-ohxqslqe"
  			strategy_name = "test2"
		}
	`
}
