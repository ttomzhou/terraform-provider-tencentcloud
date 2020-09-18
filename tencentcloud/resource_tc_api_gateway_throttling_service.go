/*
Use this resource to create api gateway throttling service.

Example Usage

```hcl
resource "tencentcloud_api_gateway_throttling_service" "service" {
	service_id        = "service-4r4xrrz4"
	strategy          = "400"
	environment_names = ["release"]
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudAPIGatewayThrottlingService() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayThrottlingServiceCreate,
		Read:   resourceTencentCloudAPIGatewayThrottlingServiceRead,
		Update: resourceTencentCloudAPIGatewayThrottlingServiceUpdate,
		Delete: resourceTencentCloudAPIGatewayThrottlingServiceDelete,

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
			"environment_names": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of Environment names.",
			},

			//compute
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
	}
}

func resourceTencentCloudAPIGatewayThrottlingServiceCreate(data *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_throttling_service.create")()
	var (
		logId = getLogId(contextNil)
		ctx   = context.WithValue(context.TODO(), logIdKey, logId)

		throttlingService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		serviceId         = data.Get("service_id").(string)
		strategy          = data.Get("strategy").(int)
		environmentName   = data.Get("environment_names").([]interface{})
		nameResults       []string
		err               error
	)

	for _, v := range environmentName {
		nameResults = append(nameResults, v.(string))
	}

	_, err = throttlingService.ModifyServiceEnvironmentStrategy(ctx, serviceId, int64(strategy), nameResults)
	if err != nil {
		return err
	}
	data.SetId(serviceId)
	return resourceTencentCloudAPIGatewayThrottlingServiceRead(data, meta)
}

func resourceTencentCloudAPIGatewayThrottlingServiceRead(data *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_throttling_service.read")()
	var (
		logId = getLogId(contextNil)
		ctx   = context.WithValue(context.TODO(), logIdKey, logId)

		throttlingService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		serviceId         = data.Id()
		err               error
	)

	environmentList, err := throttlingService.DescribeServiceEnvironmentStrategyList(ctx, serviceId)
	if err != nil {
		data.SetId("")
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
	data.Set("service_id", serviceId)
	data.Set("environments", environmentResults)
	return nil
}

func resourceTencentCloudAPIGatewayThrottlingServiceUpdate(data *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_throttling_service.update")()
	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		throttlingService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		serviceId         = data.Id()
		err               error

		strategy         int64
		environmentNames []string
	)

	oldInterfaceStrategy, newInterfaceStrategy := data.GetChange("strategy")
	if data.HasChange("strategy") {
		strategy = int64(newInterfaceStrategy.(int))
	} else {
		strategy = int64(oldInterfaceStrategy.(int))
	}

	oldInterfaceNames, newInterfaceNames := data.GetChange("environment_names")
	if data.HasChange("environment_names") {

		apiId := newInterfaceNames.([]interface{})
		for _, v := range apiId {
			environmentNames = append(environmentNames, v.(string))
		}
	} else {
		apiId := oldInterfaceNames.([]interface{})
		for _, v := range apiId {
			environmentNames = append(environmentNames, v.(string))
		}
	}

	_, err = throttlingService.ModifyServiceEnvironmentStrategy(ctx, serviceId, strategy, environmentNames)
	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudAPIGatewayThrottlingServiceDelete(data *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_throttling_service.delete")()

	var (
		logId = getLogId(contextNil)
		ctx   = context.WithValue(context.TODO(), logIdKey, logId)
		err   error

		serviceId               = data.Id()
		throttlingService       = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		strategy          int64 = 5000
		environmentNames  []string
	)

	environmentList, err := throttlingService.DescribeServiceEnvironmentStrategyList(ctx, serviceId)
	if err != nil {
		return err
	}
	for i := range environmentList {
		if environmentList[i] == nil || environmentList[i].EnvironmentName == nil {
			continue
		}
		environmentNames = append(environmentNames, *environmentList[i].EnvironmentName)
	}

	if len(environmentNames) == 0 {
		return nil
	}

	_, err = throttlingService.ModifyServiceEnvironmentStrategy(ctx, serviceId, strategy, environmentNames)
	if err != nil {
		return err
	}

	return nil
}
