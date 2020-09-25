/*
Use this resource to create api gateway service.

Example Usage

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "niceservice"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}
```

Import

api gateway service can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_service.service service-pg6ud8pa
```

*/
package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudAPIGatewayService() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayServiceCreate,
		Read:   resourceTencentCloudAPIGatewayServiceRead,
		Update: resourceTencentCloudAPIGatewayServiceUpdate,
		Delete: resourceTencentCloudAPIGatewayServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Custom service name.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(API_GATEWAY_SERVICE_PROTOCOLS),
				Description:  "Service frontend request type, such as `http`, `https`, and `http&https`.",
			},
			"service_desc": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom service description.",
			},
			"exclusive_set_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Self-deployed cluster name, which is used to specify the self-deployed cluster where the service is to be created.",
			},
			"net_type": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Network type list, which is used to specify the supported network types. " +
					"`INNER` indicates access over private network, and `OUTER` indicates access over public network.",
			},
			"ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "Pv4",
				ValidateFunc: validateAllowedStringValue(API_GATEWAY_NET_IP_VERSIONS),
				Description:  "IP version number. Valid values: `IPv4` (default value), `IPv6`.",
			},
			"set_server_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Cluster name, which is reserved and used by the tsf serverless type.",
			},
			"appid_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "User type, which is reserved and can be used by serverless users.",
			},

			// Computed values.
			"service_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service ID for query.",
			},
			"internal_sub_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private network access subdomain name.",
			},
			"outer_sub_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public network access subdomain name.",
			},
			"inner_http_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Port number for http access over private network.",
			},
			"inner_https_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Port number for https access over private network.",
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
			"usage_plan_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of attach usage plans. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"usage_plan_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the usage plan.",
						},
						"usage_plan_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the usage plan.",
						},
						"bind_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Binding type.",
						},
						"api_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the api.",
						},
					},
				},
			},
			"api_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of apis. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the api.",
						},
						"api_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the api.",
						},
						"api_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the api.",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path of the api.",
						},
						"method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Method of the api.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudAPIGatewayServiceCreate(d *schema.ResourceData, meta interface{}) error {

	defer logElapsed("resource.tencentcloud_api_gateway_service.create")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}

		serviceName      = d.Get("service_name").(string)
		protocol         = d.Get("protocol").(string)
		serviceDesc      = d.Get("service_desc").(string)
		exclusiveSetName = d.Get("exclusive_set_name").(string)
		ipVersion        = d.Get("ip_version").(string)
		setServerName    = d.Get("set_server_name").(string)
		appidType        = d.Get("appid_type").(string)
		netTypes         = helper.InterfacesStrings(d.Get("net_type").(*schema.Set).List())

		serviceId     string
		inErr, outErr error
	)

	for _, v := range netTypes {
		if v != API_GATEWAY_NET_TYPE_INNER && v != API_GATEWAY_NET_TYPE_OUTER {
			return fmt.Errorf("not support appid_type %s", v)
		}
	}
	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		serviceId, inErr = apiGatewayService.CreateService(ctx,
			serviceName,
			protocol,
			serviceDesc,
			exclusiveSetName,
			ipVersion,
			setServerName,
			appidType,
			netTypes)

		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	//wait service create ok
	if outErr := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, inErr := apiGatewayService.DescribeService(ctx, serviceId)
		if inErr != nil {
			return retryError(inErr, InternalError)
		}
		if has {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("service %s not found on server", serviceId))

	}); outErr != nil {
		return outErr
	}

	d.SetId(serviceId)
	return resourceTencentCloudAPIGatewayServiceRead(d, meta)

}

func resourceTencentCloudAPIGatewayServiceRead(d *schema.ResourceData, meta interface{}) error {

	defer logElapsed("resource.tencentcloud_api_gateway_service.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}

		info apigateway.DescribeServiceResponse
		has  bool

		serviceId     = d.Id()
		inErr, outErr error
	)

	if outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		info, has, inErr = apiGatewayService.DescribeService(ctx, serviceId)
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

	var apiList = make([]map[string]interface{}, 0, len(info.Response.ApiIdStatusSet))

	for _, item := range info.Response.ApiIdStatusSet {
		apiList = append(
			apiList, map[string]interface{}{
				"api_id":   item.ApiId,
				"api_name": item.ApiName,
				"api_desc": item.ApiDesc,
				"path":     item.Path,
				"method":   item.Method,
			})
	}

	var plans []*apigateway.ApiUsagePlan

	var planList = make([]map[string]interface{}, 0, len(info.Response.ApiIdStatusSet))
	var hasContains = make(map[string]bool)

	//from service
	if outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		plans, inErr = apiGatewayService.DescribeServiceUsagePlan(ctx, serviceId)
		if inErr != nil {
			return retryError(inErr, InternalError)
		}
		return nil
	}); outErr != nil {
		return outErr
	}

	for _, item := range plans {
		if hasContains[*item.UsagePlanId] {
			continue
		}
		hasContains[*item.UsagePlanId] = true
		planList = append(
			planList, map[string]interface{}{
				"usage_plan_id":   item.UsagePlanId,
				"usage_plan_name": item.UsagePlanName,
				"bind_type":       API_GATEWAY_TYPE_SERVICE,
				"api_id":          "",
			})
	}

	//from api
	if outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		plans, inErr = apiGatewayService.DescribeApiUsagePlan(ctx, serviceId)
		if inErr != nil {
			return retryError(inErr, InternalError)
		}
		return nil
	}); outErr != nil {
		return outErr
	}

	for _, item := range plans {
		planList = append(
			planList, map[string]interface{}{
				"usage_plan_id":   item.UsagePlanId,
				"usage_plan_name": item.UsagePlanName,
				"bind_type":       API_GATEWAY_TYPE_API,
				"api_id":          item.ApiId,
			})
	}

	d.Set("service_id", serviceId)
	d.Set("service_name", info.Response.ServiceName)
	d.Set("protocol", info.Response.Protocol)
	d.Set("service_desc", info.Response.ServiceDesc)
	d.Set("exclusive_set_name", info.Response.ExclusiveSetName)
	d.Set("ip_version", info.Response.IpVersion)
	d.Set("net_type", info.Response.NetTypes)
	d.Set("internal_sub_domain", info.Response.InternalSubDomain)
	d.Set("outer_sub_domain", info.Response.OuterSubDomain)
	d.Set("inner_http_port", info.Response.InnerHttpPort)
	d.Set("inner_https_port", info.Response.InnerHttpsPort)
	d.Set("modify_time", info.Response.ModifiedTime)
	d.Set("create_time", info.Response.CreatedTime)
	d.Set("api_list", apiList)
	d.Set("usage_plan_list", planList)

	return nil
}

func resourceTencentCloudAPIGatewayServiceUpdate(d *schema.ResourceData, meta interface{}) error {

	defer logElapsed("resource.tencentcloud_api_gateway_service.update")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}

		serviceName = d.Get("service_name").(string)
		protocol    = d.Get("protocol").(string)
		serviceDesc = d.Get("service_desc").(string)
		netTypes    = helper.InterfacesStrings(d.Get("net_type").(*schema.Set).List())

		serviceId     = d.Id()
		inErr, outErr error
	)

	for _, v := range netTypes {
		if v != API_GATEWAY_NET_TYPE_INNER && v != API_GATEWAY_NET_TYPE_OUTER {
			return fmt.Errorf("not support appid_type %s", v)
		}
	}

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if inErr = apiGatewayService.ModifyService(ctx,
			serviceId,
			serviceName,
			protocol,
			serviceDesc,
			netTypes); inErr != nil {
			return retryError(inErr, InternalError)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}
	return resourceTencentCloudAPIGatewayServiceRead(d, meta)
}

func resourceTencentCloudAPIGatewayServiceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_service.delete")()
	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		serviceId         = d.Id()
		err               error
	)

	for _, env := range API_GATEWAY_SERVICE_ENVS {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			if err = apiGatewayService.UnReleaseService(ctx,
				serviceId,
				env); err != nil {
				return retryError(err, InternalError)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if err = apiGatewayService.DeleteService(ctx,
			serviceId); err != nil {
			return retryError(err, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
