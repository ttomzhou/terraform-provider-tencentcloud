/*
Use this resource to create api gateway ip strategy.

Example Usage

```hcl
resource "tencentcloud_api_gateway_ip_strategy" "test"{
    service_id    = "service-ohxqslqe"
    strategy_name = "tf_test"
    strategy_type = "BLACK"
    strategy_data = "9.9.9.9"
}
```

Import

api gateway ip strategy can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_ip_strategy.test service-ohxqslqe#IPStrategy-q1lk8ud2
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

func resourceTencentCloudAPIGatewayIPStrategy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayIPStrategyCreate,
		Read:   resourceTencentCloudAPIGatewayIPStrategyRead,
		Update: resourceTencentCloudAPIGatewayIPStrategyUpdate,
		Delete: resourceTencentCloudAPIGatewayIPStrategyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "The ID of the API gateway service.",
			},

			"strategy_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "User defined strategy name.",
			},
			"strategy_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "Blacklist or whitelist.",
			},
			"strategy_data": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "IP address data.",
			},

			// Computed values.
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.",
			},
			"strategy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP policy ID.",
			},
		},
	}
}

func resourceTencentCloudAPIGatewayIPStrategyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_ip_strategy.create")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}

		serviceId    = d.Get("service_id").(string)
		strategyName = d.Get("strategy_name").(string)
		strategyType = d.Get("strategy_type").(string)
		strategyData = d.Get("strategy_data").(string)

		strategyId    string
		inErr, outErr error
	)
	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		strategyId, inErr = apiGatewayService.CreateIPStrategy(ctx, serviceId, strategyName, strategyType, strategyData)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	d.SetId(serviceId + FILED_SP + strategyId)

	//wait ip strategy create ok
	if outErr := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, inErr := apiGatewayService.DescribeIPStrategyStatus(ctx, serviceId, strategyId)
		if inErr != nil {
			return retryError(inErr, InternalError)
		}
		if has {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("strategyID %s not found on server", strategyId))

	}); outErr != nil {
		return outErr
	}
	return resourceTencentCloudAPIGatewayIPStrategyRead(d, meta)

}

func resourceTencentCloudAPIGatewayIPStrategyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_ip_strategy.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		serIp             = d.Id()

		inErr    error
		has      bool
		IpStatus *apigateway.IPStrategy
	)

	idSplit := strings.Split(serIp, FILED_SP)
	if len(idSplit) < 2 {
		d.SetId("")
		return nil
	}
	serviceId := idSplit[0]
	strategyId := idSplit[1]

	if outErr := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		IpStatus, has, inErr = apiGatewayService.DescribeIPStrategyStatus(ctx, serviceId, strategyId)
		if inErr != nil {
			return retryError(inErr, InternalError)
		}
		return nil
	}); outErr != nil {
		return outErr
	}

	if !has {
		d.SetId("")
		return nil
	}
	errs := []error{
		d.Set("service_id", *IpStatus.ServiceId),
		d.Set("strategy_name", *IpStatus.StrategyName),
		d.Set("strategy_type", *IpStatus.StrategyType),
		d.Set("strategy_data", *IpStatus.StrategyData),
		d.Set("strategy_id", *IpStatus.StrategyId),
		d.Set("create_time", *IpStatus.CreatedTime),
	}

	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceTencentCloudAPIGatewayIPStrategyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_ip_strategy.update")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		serIp             = d.Id()
	)

	idSplit := strings.Split(serIp, FILED_SP)
	if len(idSplit) < 2 {
		return fmt.Errorf("ip strategy is not create,can't update")
	}
	serviceId := idSplit[0]
	strategyId := idSplit[1]

	if d.HasChange("strategy_data") {
		var (
			strategyData = d.Get("strategy_data").(string)
			inErr        error
		)

		if outErr := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = apiGatewayService.UpdateIPStrategy(ctx, serviceId, strategyId, strategyData)

			if inErr != nil {
				return retryError(inErr, InternalError)
			}
			return nil
		}); outErr != nil {
			return outErr
		}
	}

	return resourceTencentCloudAPIGatewayIPStrategyRead(d, meta)
}

func resourceTencentCloudAPIGatewayIPStrategyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_ip_strategy.delete")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		serIp             = d.Id()
		outErr, inErr     error
	)

	idSplit := strings.Split(serIp, FILED_SP)
	if len(idSplit) < 2 {
		d.SetId("")
		return nil
	}
	serviceId := idSplit[0]
	strategyId := idSplit[1]

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		inErr = apiGatewayService.DeleteIPStrategy(ctx, serviceId, strategyId)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}
	return nil
}
