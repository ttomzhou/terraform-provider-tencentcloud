/*
Use this resource to create api gateway ip strategy attachment.

Example Usage

```hcl
resource "tencentcloud_api_gateway_strategy_attachment" "test"{
   service_id       = "service-ohxqslqe"
   strategy_id      = tencentcloud_api_gateway_ip_strategy.test.strategy_id
   environment_name = "release"
   bind_api_id      = "api-jbtpu758"
}
```

Import

api gateway ip strategy attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_strategy_attachment.test service-ohxqslqe#IPStrategy-nbxqk56k#api-jbtpu758#release
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudAPIGatewayStrategyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayStrategyAttachmentCreate,
		Read:   resourceTencentCloudAPIGatewayStrategyAttachmentRead,
		Delete: resourceTencentCloudAPIGatewayStrategyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "The id of the API gateway service.",
			},

			"strategy_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "The id of the API gateway strategy.",
			},
			"environment_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "The environment of the strategy association.",
			},
			"bind_api_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "The API that needs to be bound.",
			},
		},
	}
}

func resourceTencentCloudAPIGatewayStrategyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_strategy_attachment.create")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}

		serviceId  = d.Get("service_id").(string)
		strategyId = d.Get("strategy_id").(string)
		envName    = d.Get("environment_name").(string)
		bindApiId  = d.Get("bind_api_id").(string)

		inErr, outErr error
	)

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, inErr = apiGatewayService.CreateStrategyAttachment(ctx, serviceId, strategyId, envName, bindApiId)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	d.SetId(serviceId + "#" + strategyId + "#" + bindApiId + "#" + envName)
	//wait ip strategy create ok
	if outErr := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		nothas, inErr := apiGatewayService.DescribeStrategyAttachment(ctx, serviceId, strategyId, bindApiId)
		if inErr != nil {
			return retryError(inErr, InternalError)
		}
		if !nothas {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("IP strategy attachment %s not found on server", strategyId+"#"+bindApiId))

	}); outErr != nil {
		return outErr
	}
	return resourceTencentCloudAPIGatewayStrategyAttachmentRead(d, meta)
}

func resourceTencentCloudAPIGatewayStrategyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_strategy_attachment.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		attachmentId      = d.Id()

		inErr  error
		notHas bool
	)

	idSplit := strings.Split(attachmentId, FILED_SP)
	if len(idSplit) < 3 {
		d.SetId("")
		return nil
	}
	serviceId := idSplit[0]
	strategyId := idSplit[1]
	bindApiId := idSplit[2]
	envname := idSplit[3]

	if outErr := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		notHas, inErr = apiGatewayService.DescribeStrategyAttachment(ctx, serviceId, strategyId, bindApiId)
		if inErr != nil {
			return retryError(inErr, InternalError)
		}
		return nil
	}); outErr != nil {
		return outErr
	}

	if notHas {
		d.SetId("")
		return nil
	}

	errs := []error{
		d.Set("service_id", serviceId),
		d.Set("strategy_id", strategyId),
		d.Set("bind_api_id", bindApiId),
		d.Set("environment_name", envname),
	}

	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil

}

func resourceTencentCloudAPIGatewayStrategyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_strategy_attachment.delete")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		serviceId         = d.Get("service_id").(string)
		strategyId        = d.Get("strategy_id").(string)
		envName           = d.Get("environment_name").(string)
		bindApiId         = d.Get("bind_api_id").(string)
	)

	has, err := apiGatewayService.DeleteStrategyAttachment(ctx, serviceId, strategyId, envName, bindApiId)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("delete ip strategy is err")
	}

	return nil
}
