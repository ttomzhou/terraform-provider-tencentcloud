/*
Use this data source to query api gateway throttling services.

Example Usage

```hcl
resource "tencentcloud_api_gateway_throttling_service" "service" {
	service_id = "service-4r4xrrz4"
	strategy = "400"
	environment_names = ["release"]
}

data "tencentcloud_api_gateway_throttling_services" "id" {
    service_id = tencentcloud_api_gateway_throttling_service.service.service_id
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudAPIGatewayThrottlingServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayThrottlingServicesRead,

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Service ID for query.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			//compute
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of Throttling policy. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service ID for query.",
						},
						"environments": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of Throttling policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"environment_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Environment name.",
									},
									"url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Access service environment URL.",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Release status.",
									},
									"version_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Published version number.",
									},
									"strategy": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Throttling value.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAPIGatewayThrottlingServicesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_api_gateway_throttling_services.read")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}

		serviceID string
		err       error

		infos       []*apigateway.Service
		serviceIds  = make([]string, 0)
		resultLists = make([]map[string]interface{}, 0)
		ids         = make([]string, 0)
	)
	if v, ok := d.GetOk("service_id"); ok {
		serviceID = v.(string)
	}

	if serviceID == "" {
		infos, err = apiGatewayService.DescribeServicesStatus(ctx, "", "")
		if err != nil {
			return err
		}

		for i := range infos {
			result := infos[i]
			if result.ServiceId == nil {
				continue
			}
			serviceIds = append(serviceIds, *result.ServiceId)
		}
	} else {
		serviceIds = append(serviceIds, serviceID)
	}

	for i := range serviceIds {
		serviceIdTmp := serviceIds[i]
		environmentList, err := apiGatewayService.DescribeServiceEnvironmentStrategyList(ctx, serviceIdTmp)
		if err != nil {
			return err
		}

		environmentResults := make([]map[string]interface{}, 0, len(environmentList))
		for i := range environmentList {
			value := environmentList[i]
			if value == nil {
				continue
			}
			item := map[string]interface{}{
				"environment_name": value.EnvironmentName,
				"url":              value.Url,
				"status":           value.Status,
				"version_name":     value.VersionName,
				"strategy":         value.Strategy,
			}
			environmentResults = append(environmentResults, item)
		}

		resultLists = append(resultLists, map[string]interface{}{
			"service_id":   serviceIdTmp,
			"environments": environmentResults,
		})
		ids = append(ids, serviceIdTmp)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("list", resultLists)
	if err != nil {
		log.Printf("[CRITICAL]%s provider set ThrottlingServices list fail, reason:%v ", logId, err)
		return err
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), resultLists); err != nil {
			return err
		}
	}
	return nil
}
