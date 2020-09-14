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

resource "tencentcloud_api_gateway_usage_plan_attachment" "attach_service" {
	usage_plan_ids = [tencentcloud_api_gateway_usage_plan.plan.id]
	service_id     = "service-ke4o2arm"
	environment    = "test"
	bind_type      = "SERVICE"
}
```

Import

api gateway usage plan attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_usage_plan_attachment.attach_service '{"api_id":"","bind_type":"SERVICE","environment":"test","service_id":"service-pkegyqmc","usage_plan_ids":["usagePlan-26t0l0w3"]}'
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
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
			"usage_plan_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				ForceNew:    true,
				Description: "ID list of the usage plan.",
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
			"api_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				ForceNew:    true,
				Description: "API id list. This parameter will be required when `bind_type` is `API`.",
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

		serviceId   = data.Get("service_id").(string)
		environment = data.Get("environment").(string)
		bindType    = data.Get("bind_type").(string)

		apiIds []*string
		err    error
	)

	if v, ok := data.GetOk("api_ids"); ok {
		apiIdList := v.(*schema.Set).List()
		for _, v := range apiIdList {
			apiIds = append(apiIds, helper.String(v.(string)))
		}
	}

	if bindType == API_GATEWAY_TYPE_API && len(apiIds) == 0 {
		return fmt.Errorf("parameter `api_ids` is required when `bind_type` is `API`")
	}

	usagePlanIdList := data.Get("usage_plan_ids").(*schema.Set).List()
	usagePlanIds := make([]*string, 0, len(usagePlanIdList))
	for _, v := range usagePlanIdList {
		usagePlanId := v.(string)
		//check usage plan
		exist, err := checkUsagePlan(ctx, data, apiGatewayService, usagePlanId, false)
		if exist && err != nil {
			return err
		}

		usagePlanIds = append(usagePlanIds, helper.String(usagePlanId))
	}

	//check service
	exist, err := checkService(ctx, data, apiGatewayService, serviceId, false)
	if exist && err != nil {
		return err
	}

	// BindEnvironment
	if err = apiGatewayService.BindEnvironment(ctx, serviceId, environment, bindType, usagePlanIds, apiIds); err != nil {
		return err
	}

	idMap, err := json.Marshal(map[string]interface{}{
		"usage_plan_ids": usagePlanIds,
		"service_id":     serviceId,
		"environment":    environment,
		"bind_type":      bindType,
		"api_ids":        apiIds})
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

		idMap  = make(map[string]interface{})
		apiIds []interface{}
		err    error
	)
	if err = json.Unmarshal([]byte(data.Id()), &idMap); err != nil {
		return fmt.Errorf("id is broken,%s", err.Error())
	}

	var (
		usagePlanIds = idMap["usage_plan_ids"].([]interface{})
		serviceId    = idMap["service_id"].(string)
		environment  = idMap["environment"].(string)
		bindType     = idMap["bind_type"].(string)
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

	// check usage plan
	for _, usagePlanId := range usagePlanIds {
		exist, err := checkUsagePlan(ctx, data, apiGatewayService, usagePlanId.(string), true)
		if (exist && err != nil) || (!exist && err == nil) {
			return err
		}
	}

	//check service
	exist, err := checkService(ctx, data, apiGatewayService, serviceId, true)
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
			data.Set("usage_plan_ids", usagePlanIds),
			data.Set("service_id", serviceId),
			data.Set("environment", environment),
			data.Set("bind_type", bindType),
			data.Set("api_ids", apiIds),
		} {
			if err != nil {
				return err
			}
		}
		return nil
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
		usagePlanIds = idMap["usage_plan_ids"].([]interface{})
		serviceId    = idMap["service_id"].(string)
		environment  = idMap["environment"].(string)
		bindType     = idMap["bind_type"].(string)

		planIdArr = make([]*string, 0, len(usagePlanIds))
		apiIdArr  []*string
	)
	if v, ok := idMap["api_ids"]; ok {
		if v != nil {
			for _, apiId := range v.([]interface{}) {
				apiIdArr = append(apiIdArr, helper.String(apiId.(string)))
			}
		}
	}
	if len(usagePlanIds) == 0 || serviceId == "" || environment == "" || bindType == "" {
		return fmt.Errorf("id is broken")
	}
	if bindType == API_GATEWAY_TYPE_API && len(apiIdArr) == 0 {
		return fmt.Errorf("id is broken")
	}

	for _, usagePlanId := range usagePlanIds {
		planIdArr = append(planIdArr, helper.String(usagePlanId.(string)))
	}

	// BindEnvironment
	if err = apiGatewayService.UnBindEnvironment(ctx, serviceId, environment, bindType, planIdArr, apiIdArr); err != nil {
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
