package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudAPIGateWayThrottlingService(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckThrottlingServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccThrottlingService,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckThrottlingServiceExists("tencentcloud_api_gateway_throttling_service.foo"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_throttling_service.foo", "service_id", "service-4r4xrrz4"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_throttling_service.foo", "strategy", "400"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_throttling_service.foo", "environment_names.#", "1"),
				),
			},
		},
	})
}

func testAccCheckThrottlingServiceDestroy(s *terraform.State) error {
	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		throttlingService = APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_api_gateway_throttling_service" {
			continue
		}

		serviceId := rs.Primary.Attributes["service_id"]
		environmentNames := rs.Primary.Attributes["environment_names"]
		environmentList, err := throttlingService.DescribeServiceEnvironmentStrategyList(ctx, serviceId)
		if err != nil {
			log.Printf("test DescribeApiEnvironmentStrategyList: %v", err)
			return err
		}

		for i := range environmentList {
			if environmentList[i] == nil {
				continue
			}
			if !strings.Contains(environmentNames, *environmentList[i].EnvironmentName) {
				continue
			}
			if *environmentList[i].Strategy == -1 || *environmentList[i].Strategy == 5000 {
				continue
			}

			return fmt.Errorf("throttling service still not restore: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckThrottlingServiceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var (
			logId             = getLogId(contextNil)
			ctx               = context.WithValue(context.TODO(), logIdKey, logId)
			throttlingService = APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("api Getway throttling service %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("api Getway throttling service id is not set")
		}

		serviceId := rs.Primary.Attributes["service_id"]
		environmentNames := rs.Primary.Attributes["environment_names"]
		environmentList, err := throttlingService.DescribeServiceEnvironmentStrategyList(ctx, serviceId)
		if err != nil {
			log.Printf("test DescribeApiEnvironmentStrategyList: %v", err)
			return err
		}

		for i := range environmentList {
			if environmentList[i] == nil {
				continue
			}
			if !strings.Contains(environmentNames, *environmentList[i].EnvironmentName) {
				continue
			}
			if *environmentList[i].Strategy == -1 {
				return fmt.Errorf("throttling service still not set value: %s", rs.Primary.ID)
			}
		}
		return nil
	}
}

const testAccThrottlingService = `
resource "tencentcloud_api_gateway_throttling_service" "foo" {
	service_id        = "service-4r4xrrz4"
	strategy          = "400"
	environment_names = ["release"]
}
`
