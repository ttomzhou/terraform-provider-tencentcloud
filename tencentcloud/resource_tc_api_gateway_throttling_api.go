/*
Use this resource to create api gateway throttling api.

Example Usage

```hcl
resource "tencentcloud_api_gateway_throttling_api" "service" {
	service_id       = "service-4r4xrrz4"
	strategy         = "400"
	environment_name = "test"
	api_ids          = ["api-lukm33yk"]
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudAPIGatewayThrottlingAPI() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayThrottlingAPICreate,
		Read:   resourceTencentCloudAPIGatewayThrottlingAPIRead,
		Update: resourceTencentCloudAPIGatewayThrottlingAPIUpdate,
		Delete: resourceTencentCloudAPIGatewayThrottlingAPIDelete,

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateNotEmpty,
				ForceNew:     true,
				Description:  "Service ID for query.",
			},
			"strategy": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Throttling value.",
			},
			"environment_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "List of Environment names.",
			},
			"api_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of API ID.",
			},
			//compute
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
	}
}

func resourceTencentCloudAPIGatewayThrottlingAPICreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_throttling_api.create")()
	var (
		logId = getLogId(contextNil)
		ctx   = context.WithValue(context.TODO(), logIdKey, logId)

		throttlingService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		serviceId         = d.Get("service_id").(string)
		strategy          = d.Get("strategy").(int)
		environmentName   = d.Get("environment_name").(string)
		apiIds            = d.Get("api_ids").([]interface{})
		apiIdResults      []string
		err               error
	)

	for _, v := range apiIds {
		apiIdResults = append(apiIdResults, v.(string))
	}

	_, err = throttlingService.ModifyApiEnvironmentStrategy(ctx, serviceId, int64(strategy), environmentName, apiIdResults)
	if err != nil {
		return err
	}
	d.SetId(strings.Join([]string{serviceId, environmentName}, FILED_SP))
	return resourceTencentCloudAPIGatewayThrottlingAPIRead(d, meta)
}

func resourceTencentCloudAPIGatewayThrottlingAPIRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_throttling_api.read")()
	var (
		logId = getLogId(contextNil)
		ctx   = context.WithValue(context.TODO(), logIdKey, logId)

		throttlingService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		id                = d.Id()
		err               error
	)

	results := strings.Split(id, FILED_SP)
	if len(results) != 2 {
		d.SetId("")
		return fmt.Errorf("ids param is error. setId:  %s", id)
	}
	serviceId := results[0]
	environmentName := results[1]
	environmentList, err := throttlingService.DescribeApiEnvironmentStrategyList(ctx, serviceId, []string{environmentName})
	if err != nil {
		d.SetId("")
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

	d.Set("service_id", serviceId)
	d.Set("api_environment_strategies", environmentResults)

	return nil
}

func resourceTencentCloudAPIGatewayThrottlingAPIUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_throttling_api.update")()
	var (
		logId = getLogId(contextNil)
		ctx   = context.WithValue(context.TODO(), logIdKey, logId)

		throttlingService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		id                = d.Id()
		err               error

		strategy        int64
		environmentName string
		apiIds          []string
	)
	results := strings.Split(id, FILED_SP)
	if len(results) != 2 {
		return fmt.Errorf("ids param is error. setId:  %s", id)
	}
	serviceId := results[0]

	oldInterfaceStrategy, newInterfaceStrategy := d.GetChange("strategy")
	if d.HasChange("strategy") {
		strategy = int64(newInterfaceStrategy.(int))
	} else {
		strategy = int64(oldInterfaceStrategy.(int))
	}

	oldInterfaceName, newInterfaceName := d.GetChange("environment_name")
	if d.HasChange("environment_name") {
		environmentName = newInterfaceName.(string)
	} else {
		environmentName = oldInterfaceName.(string)
	}

	oldInterfaceIds, newInterfaceIds := d.GetChange("api_ids")
	if d.HasChange("api_ids") {

		apiId := newInterfaceIds.([]interface{})
		for _, v := range apiId {
			apiIds = append(apiIds, v.(string))
		}
	} else {
		apiId := oldInterfaceIds.([]interface{})
		for _, v := range apiId {
			apiIds = append(apiIds, v.(string))
		}
	}

	_, err = throttlingService.ModifyApiEnvironmentStrategy(ctx, serviceId, strategy, environmentName, apiIds)
	if err != nil {
		return err
	}

	return resourceTencentCloudAPIGatewayThrottlingAPIRead(d, meta)
}

func resourceTencentCloudAPIGatewayThrottlingAPIDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_throttling_api.delete")()

	var (
		logId = getLogId(contextNil)
		ctx   = context.WithValue(context.TODO(), logIdKey, logId)
		err   error

		id                      = d.Id()
		throttlingService       = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		strategy          int64 = 5000
		apiList           []string
	)

	results := strings.Split(id, FILED_SP)
	if len(results) != 2 {
		return fmt.Errorf("ids param is error. setId:  %s", id)
	}
	serviceId := results[0]
	environmentName := results[1]

	environmentList, err := throttlingService.DescribeApiEnvironmentStrategyList(ctx, serviceId, []string{environmentName})
	if err != nil {
		return err
	}
	for i := range environmentList {
		if environmentList[i] == nil || environmentList[i].ApiId == nil {
			continue
		}
		apiList = append(apiList, *environmentList[i].ApiId)
	}

	if len(apiList) == 0 {
		return nil
	}

	_, err = throttlingService.ModifyApiEnvironmentStrategy(ctx, serviceId, strategy, environmentName, apiList)
	if err != nil {
		return err
	}

	return nil
}
