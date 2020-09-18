/*
Use this data source to query api gateway throttling apis.

Example Usage

```hcl
resource "tencentcloud_api_gateway_throttling_api" "service" {
	service_id  	 = "service-4r4xrrz4"
	strategy 	     = "400"
	environment_name = "test"
	api_ids          = ["api-lukm33yk"]
}

data "tencentcloud_api_gateway_throttling_apis" "id" {
    service_id = tencentcloud_api_gateway_throttling_api.service.service_id
}

data "tencentcloud_api_gateway_throttling_apis" "foo" {
	service_id 		  = tencentcloud_api_gateway_throttling_api.service.service_id
	environment_names = ["release", "test"]
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

func dataSourceTencentCloudAPIGatewayThrottlingApis() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayThrottlingApisRead,

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique service ID of API.",
			},
			"environment_names": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Environment list.",
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
				Description: "A list of policies bound to API. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique service ID of API.",
						},
						"api_environment_strategies": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of throttling policies bound to API.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"api_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Unique API ID.",
									},
									"api_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom API name.",
									},
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API path.",
									},
									"method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API method.",
									},
									"strategy_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Environment throttling information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"environment_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Environment name.",
												},
												"quota": {
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
				},
			},
		},
	}
}

func dataSourceTencentCloudAPIGatewayThrottlingApisRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_api_gateway_throttling_apis.read")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}

		serviceID        string
		environmentNames []string
		err              error

		infos       []*apigateway.Service
		serviceIds  = make([]string, 0)
		resultLists = make([]map[string]interface{}, 0)
		ids         = make([]string, 0)
	)
	if v, ok := d.GetOk("service_id"); ok {
		serviceID = v.(string)
	}
	if v, ok := d.GetOk("environment_names"); ok {
		environmentNamesTmps := v.([]interface{})
		for _, v := range environmentNamesTmps {
			environmentNames = append(environmentNames, v.(string))
		}
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
		environmentList, err := apiGatewayService.DescribeApiEnvironmentStrategyList(ctx, serviceIdTmp, environmentNames)
		if err != nil {
			return err
		}

		environmentResults := make([]map[string]interface{}, 0, len(environmentList))
		for i := range environmentList {
			environmentSet := environmentList[i].EnvironmentStrategySet
			strategy_list := make([]map[string]interface{}, 0, len(environmentSet))
			for j := range environmentSet {
				if environmentSet[j] == nil {
					continue
				}
				strategy_list = append(strategy_list, map[string]interface{}{
					"environment_name": environmentSet[j].EnvironmentName,
					"quota":            environmentSet[j].Quota,
				})
			}

			item := map[string]interface{}{
				"api_id":        environmentList[i].ApiId,
				"api_name":      environmentList[i].ApiName,
				"path":          environmentList[i].Path,
				"method":        environmentList[i].Method,
				"strategy_list": strategy_list,
			}
			environmentResults = append(environmentResults, item)
		}
		resultLists = append(resultLists, map[string]interface{}{
			"service_id":                 serviceIdTmp,
			"api_environment_strategies": environmentResults,
		})
		ids = append(ids, serviceIdTmp)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("list", resultLists)
	if err != nil {
		log.Printf("[CRITICAL]%s provider set ThrottlingApiServices list fail, reason:%v ", logId, err)
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
