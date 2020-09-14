package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

var testAPIGatewayUsagePlanAttachmentResourceName = "tencentcloud_api_gateway_usage_plan_attachment"
var testAPIGatewayUsagePlanAttachmentResourceKey = testAPIGatewayUsagePlanAttachmentResourceName + ".attach_service"

func TestAccTencentCloudAPIGateWayUsagePlanAttachmentResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAPIGatewayUsagePlanAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGatewayUsagePlanAttachment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGatewayUsagePlanAttachmentExists(testAPIGatewayUsagePlanAttachmentResourceKey),
					resource.TestCheckResourceAttrSet(testAPIGatewayUsagePlanAttachmentResourceKey, "service_id"),
					resource.TestCheckResourceAttr(testAPIGatewayUsagePlanAttachmentResourceKey, "usage_plan_ids.#", "1"),
					resource.TestCheckResourceAttr(testAPIGatewayUsagePlanAttachmentResourceKey, "environment", "test"),
					resource.TestCheckResourceAttr(testAPIGatewayUsagePlanAttachmentResourceKey, "bind_type", "SERVICE"),
				),
			},
			{
				ResourceName:      testAPIGatewayUsagePlanAttachmentResourceKey,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAPIGatewayUsagePlanAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAPIGatewayUsagePlanAttachmentResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		var idMap = make(map[string]interface{}, 5)
		if err := json.Unmarshal([]byte(rs.Primary.ID), &idMap); err != nil {
			return fmt.Errorf("id is broken,%s", err.Error())
		}
		var (
			usagePlanIds = idMap["usage_plan_ids"].([]interface{})
			serviceId    = idMap["service_id"].(string)
			environment  = idMap["environment"].(string)
			bindType     = idMap["bind_type"].(string)

			apiIds []interface{}
			outErr error
			has    bool
		)
		if v, ok := idMap["api_ids"]; ok {
			if v != nil {
				apiIds = v.([]interface{})
			}
		}
		if len(usagePlanIds) == 0 || serviceId == "" || environment == "" || bindType == "" {
			return fmt.Errorf("id is broken")
		}
		if bindType == API_GATEWAY_TYPE_API && len(apiIds) == 0 {
			return fmt.Errorf("id is broken")
		}

		service := APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		for _, usagePlanId := range usagePlanIds {
			_, has, outErr = service.DescribeUsagePlan(ctx, usagePlanId.(string))
			if outErr != nil {
				_, has, outErr = service.DescribeUsagePlan(ctx, usagePlanId.(string))
			}
			if outErr != nil {
				return outErr
			}
			if !has {
				return nil
			}
		}

		_, has, outErr = service.DescribeService(ctx, serviceId)
		if outErr != nil {
			_, has, outErr = service.DescribeService(ctx, serviceId)
		}
		if outErr != nil {
			return outErr
		}
		if !has {
			return nil
		}

		var plans []*apigateway.ApiUsagePlan

		if bindType == API_GATEWAY_TYPE_API {
			plans, outErr = service.DescribeApiUsagePlan(ctx, serviceId)
			if outErr != nil {
				plans, outErr = service.DescribeApiUsagePlan(ctx, serviceId)
			}
			if outErr != nil {
				return outErr
			}
		} else {
			plans, outErr = service.DescribeServiceUsagePlan(ctx, serviceId)
			if outErr != nil {
				plans, outErr = service.DescribeServiceUsagePlan(ctx, serviceId)
			}
			if outErr != nil {
				return outErr
			}
		}

		var (
			usagePlanMap = make(map[string]bool)
			apiIdMap     = make(map[string]bool)
		)
		for _, usagePlanId := range usagePlanIds {
			usagePlanMap[usagePlanId.(string)] = true
		}
		for _, apiId := range apiIds {
			apiIdMap[apiId.(string)] = true
		}

		for _, plan := range plans {
			if usagePlanMap[*plan.UsagePlanId] && *plan.Environment == environment {
				if bindType == API_GATEWAY_TYPE_API {
					if plan.ApiId != nil && apiIdMap[*plan.ApiId] {
						return fmt.Errorf("attachment  %s still exist on server", rs.Primary.ID)
					}
				} else {
					return fmt.Errorf("attachment  %s still exist on server", rs.Primary.ID)
				}
			}
		}

		return nil
	}
	return nil
}

func testAccCheckAPIGatewayUsagePlanAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		var idMap = make(map[string]interface{}, 5)
		if err := json.Unmarshal([]byte(rs.Primary.ID), &idMap); err != nil {
			return fmt.Errorf("id is broken,%s", err.Error())
		}
		var (
			usagePlanIds = idMap["usage_plan_ids"].([]interface{})
			serviceId    = idMap["service_id"].(string)
			environment  = idMap["environment"].(string)
			bindType     = idMap["bind_type"].(string)

			apiIds []interface{}
			outErr error
			has    bool
		)
		if v, ok := idMap["api_ids"]; ok {
			if v != nil {
				apiIds = v.([]interface{})
			}
		}
		if len(usagePlanIds) == 0 || serviceId == "" || environment == "" || bindType == "" {
			return fmt.Errorf("id is broken")
		}
		if bindType == API_GATEWAY_TYPE_API && len(apiIds) == 0 {
			return fmt.Errorf("id is broken")
		}

		service := APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		for _, usagePlanId := range usagePlanIds {
			_, has, outErr = service.DescribeUsagePlan(ctx, usagePlanId.(string))
			if outErr != nil {
				_, has, outErr = service.DescribeUsagePlan(ctx, usagePlanId.(string))
			}
			if outErr != nil {
				return outErr
			}
			if !has {
				return fmt.Errorf("usage plan %s not exsit on server", usagePlanId.(string))
			}
		}

		_, has, outErr = service.DescribeService(ctx, serviceId)
		if outErr != nil {
			_, has, outErr = service.DescribeService(ctx, serviceId)
		}
		if outErr != nil {
			return outErr
		}
		if !has {
			return fmt.Errorf("service %s not exsit on server", serviceId)
		}

		var plans []*apigateway.ApiUsagePlan

		if bindType == API_GATEWAY_TYPE_API {
			plans, outErr = service.DescribeApiUsagePlan(ctx, serviceId)
			if outErr != nil {
				plans, outErr = service.DescribeApiUsagePlan(ctx, serviceId)
			}
			if outErr != nil {
				return outErr
			}
		} else {
			plans, outErr = service.DescribeServiceUsagePlan(ctx, serviceId)
			if outErr != nil {
				plans, outErr = service.DescribeServiceUsagePlan(ctx, serviceId)
			}
			if outErr != nil {
				return outErr
			}
		}

		var (
			usagePlanMap = make(map[string]bool)
			apiIdMap     = make(map[string]bool)
		)
		for _, usagePlanId := range usagePlanIds {
			usagePlanMap[usagePlanId.(string)] = true
		}
		for _, apiId := range apiIds {
			apiIdMap[apiId.(string)] = true
		}

		for _, plan := range plans {
			if usagePlanMap[*plan.UsagePlanId] && *plan.Environment == environment {
				if bindType == API_GATEWAY_TYPE_API {
					if plan.ApiId != nil && apiIdMap[*plan.ApiId] {
						return nil
					}
				} else {
					return nil
				}
			}
		}
		return fmt.Errorf("attachment  %s not exist on server", rs.Primary.ID)
	}
}

const testAccAPIGatewayUsagePlanAttachment = `
	resource "tencentcloud_api_gateway_usage_plan" "plan" {
  		usage_plan_name         = "my_plan"
  		usage_plan_desc         = "nice plan"
  		max_request_num         = 100
  		max_request_num_pre_sec = 10
	}

	resource "tencentcloud_api_gateway_usage_plan_attachment" "attach_service" {
  		usage_plan_ids = [tencentcloud_api_gateway_usage_plan.plan.id]
  		service_id     = "service-ke4o2arm"
  		environment    = "test"
  		bind_type      = "SERVICE"
	}
`
