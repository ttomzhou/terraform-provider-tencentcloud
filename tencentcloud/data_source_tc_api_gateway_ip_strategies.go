/*
Use this data source to query api gateway ip strategy.

Example Usage

```hcl
resource "tencentcloud_api_gateway_ip_strategy" "test"{
	service_id 	  = "service-ohxqslqe"
	strategy_name = "tf_test"
	strategy_type = "BLACK"
	strategy_data = "9.9.9.9"
}

data "tencentcloud_api_gateway_ip_strategies" "id" {
	service_id = tencentcloud_api_gateway_ip_strategy.test.service_id
}
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

func dataSourceTencentCloudAPIGatewayIpStrategy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayIpStrategyRead,

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the service to be queried.",
			},
			"strategy_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of ip policy.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values.
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of strategy. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"strategy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the strategy.",
						},
						"strategy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the strategy.",
						},
						"strategy_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the strategy. Valid values: `WHITE` (white list) and `BLACK` (black list).",
						},
						"ip_list": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The list of ip.",
						},
						"service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The service id.",
						},
						"bind_api_total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of api bound to the strategy.",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.",
						},
						"attach_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of bound api details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The service id.",
									},
									"api_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The api id.",
									},
									"api_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The api interface description.",
									},
									"api_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the api interface.",
									},
									"vpc_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The vpc id.",
									},
									"uniq_vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The vpc unique id.",
									},
									"api_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The api type. The values are `NORMAL` (common api) and `TSF` (microservice api).",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The api protocol.",
									},
									"auth_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The api authentication type. The value is `SECRET` (key pair authentication), `NONE` (no authentication), and `OAUTH`.",
									},
									"api_business_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of oauth api. This field is valid when the authType is `OAUTH`, and the values are `NORMAL` (business api) and `OAUTH` (authorization api).",
									},
									"auth_relation_api_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique id of the associated authorization api, which takes effect when the authType is `OAUTH` and `ApiBusinessType` is normal. Identifies the unique id of the oauth2.0 authorization api bound to the business api.",
									},
									"tags": {
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Computed:    true,
										Description: "The label information associated with the api.",
									},
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The api path.",
									},
									"method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The api request method.",
									},
									"relation_business_api_ids": {
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Computed:    true,
										Description: "List of business api associated with authorized api.",
									},
									"oauth_config": {
										Type:        schema.TypeMap,
										Computed:    true,
										Description: "OAUTH configuration information. It takes effect when authType is `OAUTH`.",
									},
									"modify_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.",
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

func dataSourceTencentCloudAPIGatewayIpStrategyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_api_gateway_ip_strategy.read")

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}

		serviceId    = d.Get("service_id").(string)
		strategyName = d.Get("strategy_name").(string)
		infos        []*apigateway.IPStrategy
		list         []map[string]interface{}

		err error
	)

	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		infos, err = apiGatewayService.DescribeIPStrategysStatus(ctx, serviceId, strategyName)
		if err != nil {
			return retryError(err)
		}
		return nil
	}); err != nil {
		return err
	}

	for _, info := range infos {
		var attachListInfo []map[string]interface{}

		for _, env := range API_GATEWAY_SERVICE_ENVS {
			var strategy *apigateway.IPStrategy
			if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				strategy, err = apiGatewayService.DescribeIPStrategies(ctx, serviceId, *info.StrategyId, env)
				if err != nil {
					return retryError(err)
				}
				return nil
			}); err != nil {
				return err
			}

			for _, api := range strategy.BindApis {
				attachListInfo = append(attachListInfo, map[string]interface{}{
					"service_id":                api.ServiceId,
					"api_id":                    api.ApiId,
					"api_desc":                  api.ApiDesc,
					"api_name":                  api.ApiName,
					"vpc_id":                    api.VpcId,
					"uniq_vpc_id":               api.UniqVpcId,
					"api_type":                  api.ApiType,
					"protocol":                  api.Protocol,
					"auth_type":                 api.AuthType,
					"api_business_type":         api.ApiBusinessType,
					"auth_relation_api_id":      api.AuthRelationApiId,
					"tags":                      api.Tags,
					"path":                      api.Path,
					"method":                    api.Method,
					"relation_business_api_ids": api.RelationBuniessApiIds,
					"oauth_config":              flattenOauthConfigMappings(api.OauthConfig),
					"modify_time":               api.ModifiedTime,
					"create_time":               api.CreatedTime,
				})
			}
		}

		infoMap := map[string]interface{}{
			"strategy_id":          info.StrategyId,
			"strategy_name":        info.StrategyName,
			"strategy_type":        info.StrategyType,
			"ip_list":              info.StrategyData,
			"service_id":           info.ServiceId,
			"bind_api_total_count": info.BindApiTotalCount,
			"modify_time":          info.ModifiedTime,
			"create_time":          info.CreatedTime,
			"attach_list":          attachListInfo,
		}

		list = append(list, infoMap)
	}

	if err = d.Set("list", list); err != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s", logId, err.Error())
	}

	d.SetId(strings.Join([]string{serviceId, strategyName}, FILED_SP))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		return writeToFile(output.(string), list)
	}
	return nil
}
