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

func TestAccTencentCloudAPIGateWayThrottlingAPI(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckThrottlingAPIDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccThrottlingAPI,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckThrottlingAPIExists("tencentcloud_api_gateway_throttling_api.foo"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_throttling_api.foo", "service_id", "service-4r4xrrz4"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_throttling_api.foo", "strategy", "400"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_throttling_api.foo", "environment_name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_throttling_api.foo", "api_ids.#", "1"),
				),
			},
		},
	})
}

func testAccCheckThrottlingAPIDestroy(s *terraform.State) error {
	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		throttlingService = APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_api_gateway_throttling_api" {
			continue
		}

		serviceId := rs.Primary.Attributes["service_id"]
		environmentName := rs.Primary.Attributes["environment_name"]
		apiIds := rs.Primary.Attributes["api_ids"]
		environmentList, err := throttlingService.DescribeApiEnvironmentStrategyList(ctx, serviceId, []string{environmentName})
		if err != nil {
			log.Printf("test DescribeApiEnvironmentStrategyList: %v", err)
			return err
		}

		for i := range environmentList {
			if environmentList[i] == nil {
				continue
			}
			if !strings.Contains(apiIds, *environmentList[i].ApiId) {
				continue
			}
			environmentSet := environmentList[i].EnvironmentStrategySet
			for j := range environmentSet {
				if environmentSet[j] == nil {
					continue
				}
				if *environmentSet[j].EnvironmentName != environmentName {
					continue
				}

				if *environmentSet[j].Quota == -1 || *environmentSet[j].Quota == 5000 {
					continue
				}
				return fmt.Errorf("throttling api still not restore: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckThrottlingAPIExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var (
			logId             = getLogId(contextNil)
			ctx               = context.WithValue(context.TODO(), logIdKey, logId)
			throttlingService = APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("api Getway throttling api %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("api Getway throttling api id is not set")
		}
		serviceId := rs.Primary.Attributes["service_id"]
		environmentName := rs.Primary.Attributes["environment_name"]
		apiIds := rs.Primary.Attributes["api_ids"]
		environmentList, err := throttlingService.DescribeApiEnvironmentStrategyList(ctx, serviceId, []string{environmentName})
		if err != nil {
			log.Printf("test DescribeApiEnvironmentStrategyList: %v", err)
			return err
		}

		for i := range environmentList {
			if environmentList[i] == nil {
				continue
			}
			if !strings.Contains(apiIds, *environmentList[i].ApiId) {
				continue
			}
			environmentSet := environmentList[i].EnvironmentStrategySet
			for j := range environmentSet {
				if environmentSet[j] == nil {
					continue
				}
				if *environmentSet[j].EnvironmentName != environmentName {
					continue
				}

				if *environmentSet[j].Quota == -1 {
					return fmt.Errorf("throttling api still not set value: %s", rs.Primary.ID)
				}

			}
		}
		return nil
	}
}

const testAccThrottlingAPI = `
resource "tencentcloud_api_gateway_throttling_api" "foo" {
	service_id = "service-4r4xrrz4"
	strategy = "400"
	environment_name = "test"
	api_ids = ["api-lukm33yk"]
}
`
