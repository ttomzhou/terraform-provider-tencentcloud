/*
Use this resource to create custom domain of api gateway.

Example Usage

```hcl
resource "tencentcloud_api_gateway_custom_domain" "service" {
	service_id 			= "service-ohxqslqe"
	sub_domain 			= "tic-test.dnsv1.com"
	protocol   			= "http"
	net_type   			= "OUTER"
	is_default_mapping  = "false"
	default_domain 		= "service-ohxqslqe-1259649581.gz.apigw.tencentcs.com"
	path_mappings 		= ["/good#test","/root#release"]
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

func resourceTencentCloudAPIGatewayCustomDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayCustomDomainCreate,
		Read:   resourceTencentCloudAPIGatewayCustomDomainRead,
		Update: resourceTencentCloudAPIGatewayCustomDomainUpdate,
		Delete: resourceTencentCloudAPIGatewayCustomDomainDelete,

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateNotEmpty,
				ForceNew:     true,
				Description:  "Unique service ID.",
			},
			"sub_domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Custom domain name to be bound.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "Protocol supported by service. Valid values: http, https, http&https.",
			},
			"net_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "Network type. Valid values: OUTER, INNER.",
			},
			"is_default_mapping": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "Whether the default path mapping is used. The default value is true. If the value is false, the custom path mapping will be used and PathMappingSet will be required in this case.",
			},
			"default_domain": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "Default domain name.",
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Unique certificate ID of the custom domain name to be bound. The certificate can be uploaded if Protocol is https or http&https.",
			},
			"path_mappings": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Custom domain name path mapping. It can contain up to 3 Environment values which can be set to only test, prepub, and release, respectively.eg: path#environment.",
			},
			//compute
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Domain name resolution status. True: success; False: failure.",
			},
		},
	}
}

func resourceTencentCloudAPIGatewayCustomDomainCreate(data *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_custom_domain.create")()
	var (
		logId = getLogId(contextNil)
		ctx   = context.WithValue(context.TODO(), logIdKey, logId)

		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		serviceId         = data.Get("service_id").(string)
		subDomain         = data.Get("sub_domain").(string)
		protocol          = data.Get("protocol").(string)
		netType           = data.Get("net_type").(string)
		defaultDomain     = data.Get("default_domain").(string)
		isDefaultMapping  = data.Get("is_default_mapping").(bool)
		certificateId     string
		pathMappings      []string

		err error
	)

	if v, ok := data.GetOk("certificate_id"); ok {
		certificateId = v.(string)
	}
	if v, ok := data.GetOk("path_mappings"); ok {
		pathMappingTmps := v.(*schema.Set).List()
		for _, v := range pathMappingTmps {
			pathMappings = append(pathMappings, v.(string))
		}
	}

	err = apiGatewayService.BindSubDomainService(ctx, serviceId, subDomain, protocol, netType, defaultDomain, isDefaultMapping, certificateId, pathMappings)
	if err != nil {
		return err
	}

	data.SetId(strings.Join([]string{serviceId, subDomain}, FILED_SP))
	return resourceTencentCloudAPIGatewayCustomDomainRead(data, meta)
}

func resourceTencentCloudAPIGatewayCustomDomainRead(data *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_custom_domain.read")()
	var (
		logId = getLogId(contextNil)
		ctx   = context.WithValue(context.TODO(), logIdKey, logId)

		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		id                = data.Id()
		err               error
	)

	results := strings.Split(id, FILED_SP)
	if len(results) != 2 {
		data.SetId("")
		return fmt.Errorf("ids param is error. id:  %s", id)
	}
	serviceId := results[0]
	subDomain := results[1]
	resultList, err := apiGatewayService.DescribeServiceSubDomainsService(ctx, serviceId, subDomain)
	if err != nil {
		data.SetId("")
		return err
	}

	if len(resultList) == 0 {
		data.SetId("")
		return fmt.Errorf("custom domain: %s not found of service: %s.", subDomain, serviceId)
	}

	resultInfo := resultList[0]
	info, err := apiGatewayService.DescribeServiceSubDomainMappings(ctx, serviceId, *resultInfo.DomainName)
	if err != nil {
		data.SetId("")
		return fmt.Errorf("DescribeServiceSubDomainMappings err: %s", err.Error())
	}
	pathMap := make([]string, 0, len(info.PathMappingSet))
	for i := range info.PathMappingSet {
		pathMap = append(pathMap, *info.PathMappingSet[i].Path+FILED_SP+*info.PathMappingSet[i].Environment)
	}
	data.Set("path_mappings", pathMap)
	data.Set("domain_name", resultInfo.DomainName)
	data.Set("status", resultInfo.Status)
	data.Set("certificate_id", resultInfo.CertificateId)
	data.Set("is_default_mapping", resultInfo.IsDefaultMapping)
	data.Set("protocol", resultInfo.Protocol)
	data.Set("net_type", resultInfo.NetType)
	data.Set("service_id", serviceId)
	return nil
}

func resourceTencentCloudAPIGatewayCustomDomainUpdate(data *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_custom_domain.update")()
	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		id                = data.Id()

		subDomain        string
		isDefaultMapping bool
		certificateId    string
		protocol         string
		netType          string
		pathMappings     []string
	)
	results := strings.Split(id, FILED_SP)
	if len(results) != 2 {
		return fmt.Errorf("ids param is error. setId:  %s", id)
	}
	serviceId := results[0]

	oldInterfaceSubDomain, newInterfaceSubDomain := data.GetChange("sub_domain")
	if data.HasChange("sub_domain") {
		subDomain = newInterfaceSubDomain.(string)
	} else {
		subDomain = oldInterfaceSubDomain.(string)
	}

	oldInterfaceName, newInterfaceName := data.GetChange("is_default_mapping")
	if data.HasChange("is_default_mapping") {
		isDefaultMapping = newInterfaceName.(bool)
	} else {
		isDefaultMapping = oldInterfaceName.(bool)
	}

	_, newInterfaceCertificateId := data.GetChange("certificate_id")
	if data.HasChange("certificate_id") {
		certificateId = newInterfaceCertificateId.(string)
	}

	_, newInterfaceNetType := data.GetChange("net_type")
	if data.HasChange("net_type") {
		netType = newInterfaceNetType.(string)
	}

	_, newInterfaceProtocol := data.GetChange("protocol")
	if data.HasChange("protocol") {
		protocol = newInterfaceProtocol.(string)
	}

	_, newInterfacePathMappings := data.GetChange("path_mappings")
	if data.HasChange("path_mappings") {
		pathMappingsTmp := newInterfacePathMappings.([]interface{})
		for _, v := range pathMappingsTmp {
			pathMappings = append(pathMappings, v.(string))
		}
	}

	return apiGatewayService.ModifySubDomainService(ctx, serviceId, subDomain, isDefaultMapping, certificateId, protocol, netType, pathMappings)
}

func resourceTencentCloudAPIGatewayCustomDomainDelete(data *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_custom_domain.delete")()

	var (
		logId = getLogId(contextNil)
		ctx   = context.WithValue(context.TODO(), logIdKey, logId)

		id                = data.Id()
		apigatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	results := strings.Split(id, FILED_SP)
	if len(results) != 2 {
		return fmt.Errorf("ids param is error. setId:  %s", id)
	}
	serviceId := results[0]
	subDomain := results[1]

	return apigatewayService.UnBindSubDomainService(ctx, serviceId, subDomain)
}
