package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudApiStrategyAttachment_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testApiStrategyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testApiStrategyAttachment_basic,
				Check: resource.ComposeTestCheckFunc(
					testApiStrategyAttachmentExists("tencentcloud_api_gateway_strategy_attachment.test"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_strategy_attachment.test", "service_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_strategy_attachment.test", "strategy_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_strategy_attachment.test", "environment_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_strategy_attachment.test", "bind_api_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_api_gateway_strategy_attachment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testApiStrategyAttachmentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_api_gateway_strategy_attachment" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		serviceId := idSplit[0]
		strategyId := idSplit[1]
		bindApiId := idSplit[2]

		Has, err := service.DescribeStrategyAttachment(ctx, serviceId, strategyId, bindApiId)
		if err != nil {
			return err
		}

		if Has {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][api ip strategy][Destroy] check: api ip strategy still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testApiStrategyAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][api ip strategy][Exists] check:  %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][api ip strategy][Exists] check: id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		serviceId := idSplit[0]
		strategyId := idSplit[1]
		bindApiId := idSplit[2]
		Has, err := service.DescribeStrategyAttachment(ctx, serviceId, strategyId, bindApiId)
		if err != nil {
			return err
		}

		if Has {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][api ip strategy][Exists] check: not exists: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testApiStrategyAttachment_basic = `

resource "tencentcloud_api_gateway_ip_strategy" "test"{
   service_id = "service-ohxqslqe"
   strategy_name = "tf_test"
   strategy_type = "BLACK"
   strategy_data = "9.9.9.9"
}

resource "tencentcloud_api_gateway_strategy_attachment" "att_test"{
   service_id = "service-ohxqslqe"
   strategy_id = tencentcloud_api_gateway_ip_strategy.test.strategy_id
   environment_name = "release"
   bind_api_id = "api-jbtpu758"
}
}
`
