package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testAPIGatewayAPIKeysDataSourceName = "data.tencentcloud_api_gateway_api_keys"

func TestAccTencentAPIGatewayAPIKeysDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAPIGatewayAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTestAccTencentAPIGatewayAPIKeys(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAPIGatewayAPIKeyExists(testAPIGatewayAPIKeyResourceName+".test"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeysDataSourceName+".name", "list.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeysDataSourceName+".name", "list.0.api_key_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeysDataSourceName+".name", "list.0.status"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeysDataSourceName+".name", "list.0.access_key_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeysDataSourceName+".name", "list.0.access_key_secret"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeysDataSourceName+".name", "list.0.modify_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeysDataSourceName+".name", "list.0.create_time"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIKeysDataSourceName+".id", "list.#", "1"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeysDataSourceName+".id", "list.0.api_key_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeysDataSourceName+".id", "list.0.status"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeysDataSourceName+".id", "list.0.access_key_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeysDataSourceName+".id", "list.0.access_key_secret"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeysDataSourceName+".id", "list.0.modify_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeysDataSourceName+".id", "list.0.create_time"),
				),
			},
		},
	})
}

func testAccTestAccTencentAPIGatewayAPIKeys() string {
	return `
resource "tencentcloud_api_gateway_api_key" "test" {
  secret_name = "my_api_key"
  status      = "on"
}
data "tencentcloud_api_gateway_api_keys" "name" {
  secret_name = tencentcloud_api_gateway_api_key.test.secret_name
}

data "tencentcloud_api_gateway_api_keys" "id" {
  access_key_id = tencentcloud_api_gateway_api_key.test.access_key_id
}
`
}
