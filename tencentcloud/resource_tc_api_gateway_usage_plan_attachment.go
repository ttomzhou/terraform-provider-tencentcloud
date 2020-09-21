/*
Use this resource to attach api gateway usage plan to service.

Example Usage

```hcl
resource "tencentcloud_api_gateway_usage_plan" "plan" {
	usage_plan_name         = "my_plan"
	usage_plan_desc         = "nice plan"
	max_request_num         = 100
	max_request_num_pre_sec = 10
}

resource "tencentcloud_api_gateway_service" "service" {
	service_name = "niceservice"
	protocol     = "http&https"
	service_desc = "your nice service"
	net_type     = ["INNER", "OUTER"]
	ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_usage_plan_attachment" "attach_service" {
	usage_plan_id  = tencentcloud_api_gateway_usage_plan.plan.id
	service_id     = tencentcloud_api_gateway_service.service.id
	environment    = "test"
	bind_type      = "SERVICE"
}
```

Import

api gateway usage plan attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_usage_plan_attachment.attach_service '{"api_id":"","bind_type":"SERVICE","environment":"test","service_id":"service-pkegyqmc","usage_plan_id":"usagePlan-26t0l0w3"}'
```
*/
package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

func resourceTencentCloudAPIGatewayUsagePlanAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayUsagePlanAttachmentCreate,
		Read:   resourceTencentCloudAPIGatewayUsagePlanAttachmentRead,
		Delete: resourceTencentCloudAPIGatewayUsagePlanAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"usage_plan_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the usage plan.",
			},
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the service.",
			},
			"environment": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(API_GATEWAY_SERVICE_ENVS),
				Description:  "Environment to be bound `test`,`prepub` or `release`.",
			},
			"bind_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      API_GATEWAY_TYPE_SERVICE,
				ValidateFunc: validateAllowedStringValue(API_GATEWAY_TYPES),
				Description:  "Binding type. Valid values: `API`, `SERVICE` (default value).",
			},
			"api_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "API id. This parameter will be required when `bind_type` is `API`.",
			},
		},
	}
}

func resourceTencentCloudAPIGatewayUsagePlanAttachmentCreate(data *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_usage_plan_attachment.create")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}

		usagePlanId = data.Get("usage_plan_id").(string)
		serviceId   = data.Get("service_id").(string)
		environment = data.Get("environment").(string)
		bindType    = data.Get("bind_type").(string)

		apiId string
		err   error
	)

	if v, ok := data.GetOk("api_id"); ok {
		apiId = v.(string)
	}

	if bindType == API_GATEWAY_TYPE_API && apiId == "" {
		return fmt.Errorf("parameter `api_ids` is required when `bind_type` is `API`")
	}

	//check usage plan
	exist, err := checkUsagePlan(ctx, data, apiGatewayService, usagePlanId, false)
	if exist && err != nil {
		return err
	}

	//check service
	exist, err = checkService(ctx, data, apiGatewayService, serviceId, false)
	if exist && err != nil {
		return err
	}

	// BindEnvironment
	if err = apiGatewayService.BindEnvironment(ctx, serviceId, environment, bindType, usagePlanId, apiId); err != nil {
		return err
	}

	idMap, err := json.Marshal(map[string]interface{}{
		"usage_plan_id": usagePlanId,
		"service_id":    serviceId,
		"environment":   environment,
		"bind_type":     bindType,
		"api_id":        apiId})
	if err != nil {
		return fmt.Errorf("build id json fail,%s", err.Error())
	}

	data.SetId(string(idMap))

	return resourceTencentCloudAPIGatewayUsagePlanAttachmentRead(data, meta)
}

func resourceTencentCloudAPIGatewayUsagePlanAttachmentRead(data *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_usage_plan_attachment.read")()
	defer inconsistentCheck(data, meta)()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}

		idMap = make(map[string]interface{})
		err   error
	)
	if err = json.Unmarshal([]byte(data.Id()), &idMap); err != nil {
		return fmt.Errorf("id is broken,%s", err.Error())
	}

	var (
		usagePlanId = idMap["usage_plan_id"].(string)
		serviceId   = idMap["service_id"].(string)
		environment = idMap["environment"].(string)
		bindType    = idMap["bind_type"].(string)
		apiId       = idMap["api_id"].(string)
	)

	if usagePlanId == "" || serviceId == "" || environment == "" || bindType == "" {
		return fmt.Errorf("id is broken")
	}
	if bindType == API_GATEWAY_TYPE_API && apiId == "" {
		return fmt.Errorf("id is broken")
	}

	// check usage plan
	exist, err := checkUsagePlan(ctx, data, apiGatewayService, usagePlanId, true)
	if (exist && err != nil) || (!exist && err == nil) {
		return err
	}

	//check service
	exist, err = checkService(ctx, data, apiGatewayService, serviceId, true)
	if (exist && err != nil) || (!exist && err == nil) {
		return err
	}

	plans := make([]*apigateway.ApiUsagePlan, 0)
	if bindType == API_GATEWAY_TYPE_API {
		if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			plans, err = apiGatewayService.DescribeApiUsagePlan(ctx, serviceId)
			if err != nil {
				return retryError(err, InternalError)
			}
			return nil
		}); err != nil {
			return err
		}
	} else {
		if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			plans, err = apiGatewayService.DescribeServiceUsagePlan(ctx, serviceId)
			if err != nil {
				return retryError(err, InternalError)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	var setData = func() error {
		for _, err := range []error{
			data.Set("usage_plan_id", usagePlanId),
			data.Set("service_id", serviceId),
			data.Set("environment", environment),
			data.Set("bind_type", bindType),
			data.Set("api_id", apiId),
		} {
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, plan := range plans {
		if *plan.UsagePlanId == usagePlanId && *plan.Environment == environment {
			if bindType == API_GATEWAY_TYPE_API {
				if plan.ApiId != nil && *plan.ApiId == apiId {
					return setData()
				}
			} else {
				return setData()
			}
		}
	}

	data.SetId("")
	return nil
}

func resourceTencentCloudAPIGatewayUsagePlanAttachmentDelete(data *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_usage_plan_attachment.delete")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}

		idMap = make(map[string]interface{})
		err   error
	)

	if err = json.Unmarshal([]byte(data.Id()), &idMap); err != nil {
		return fmt.Errorf("id is broken,%s", err.Error())
	}

	var (
		usagePlanId = idMap["usage_plan_id"].(string)
		serviceId   = idMap["service_id"].(string)
		environment = idMap["environment"].(string)
		bindType    = idMap["bind_type"].(string)
		apiId       = idMap["api_id"].(string)
	)

	if usagePlanId == "" || serviceId == "" || environment == "" || bindType == "" {
		return fmt.Errorf("id is broken")
	}
	if bindType == API_GATEWAY_TYPE_API && apiId == "" {
		return fmt.Errorf("id is broken")
	}

	// BindEnvironment
	if err = apiGatewayService.UnBindEnvironment(ctx, serviceId, environment, bindType, usagePlanId, apiId); err != nil {
		return err
	}

	return nil
}

func checkUsagePlan(ctx context.Context, data *schema.ResourceData, api APIGatewayService, usagePlanId string, isSetId bool) (bool, error) {
	var (
		err error
		has bool
	)
	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, err = api.DescribeUsagePlan(ctx, usagePlanId)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return true, err
	}

	if !has {
		if isSetId {
			data.SetId("")
			return false, nil
		} else {
			return true, fmt.Errorf("usage plan %s not exist", usagePlanId)
		}
	}

	return true, nil
}

func checkService(ctx context.Context, data *schema.ResourceData, api APIGatewayService, serviceId string, isSetId bool) (bool, error) {
	var (
		err error
		has bool
	)
	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, err = api.DescribeService(ctx, serviceId)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return true, err
	}

	if !has {
		if isSetId {
			data.SetId("")
			return false, nil
		} else {
			return true, fmt.Errorf("service %s not exist", serviceId)
		}
	}

	return true, nil
}
