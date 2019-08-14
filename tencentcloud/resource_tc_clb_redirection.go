/*
Provides a resource to create a CLB redirection.

Example Usage

```hcl
resource "tencentcloud_clb_redirection" "foo" {
  clb_id             = "lb-p7olt9e5"
  source_listener_id = "lbl-jc1dx6ju"
  target_listener_id = "lbl-asj1hzuo"
  source_rule_id     = "loc-ft8fmngv"
  target_rule_id     = "loc-4xxr2cy7"
}
```

Import

CLB redirection can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_redirection.foo loc-ft8fmngv#loc-4xxr2cy7#lbl-jc1dx6ju#lbl-asj1hzuo#lb-p7olt9e5
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
)

func resourceTencentCloudClbRedirection() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbRedirectionCreate,
		Read:   resourceTencentCloudClbRedirectionRead,
		Delete: resourceTencentCloudClbRedirectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of CLB instance.",
			},
			"source_listener_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Id of source listener.",
			},
			"target_listener_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Id of source listener.",
			},
			"source_rule_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Rule id of source listener.",
			},
			"target_rule_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Rule id of target listener.",
			},
		},
	}
}

func resourceTencentCloudClbRedirectionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_redirection.create")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)

	clbId := d.Get("clb_id").(string)
	sourceListenerId := d.Get("source_listener_id").(string)
	targertListenerId := d.Get("target_listener_id").(string)
	sourceLocId := d.Get("source_rule_id").(string)
	targetLocId := d.Get("target_rule_id").(string)

	request := clb.NewManualRewriteRequest()

	request.LoadBalancerId = stringToPointer(clbId)
	request.SourceListenerId = stringToPointer(sourceListenerId)
	request.TargetListenerId = stringToPointer(targertListenerId)

	var rewriteInfo clb.RewriteLocationMap
	rewriteInfo.SourceLocationId = stringToPointer(sourceLocId)
	rewriteInfo.TargetLocationId = stringToPointer(targetLocId)
	request.RewriteInfos = []*clb.RewriteLocationMap{&rewriteInfo}
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().ManualRewrite(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			requestId := *response.Response.RequestId
			retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
			if retryErr != nil {
				return resource.NonRetryableError(retryErr)
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create clb redirection failed, reason:%s\n ", logId, err.Error())
		return err
	}
	d.SetId(sourceLocId + "#" + targetLocId + "#" + sourceListenerId + "#" + targertListenerId + "#" + clbId)

	return resourceTencentCloudClbRedirectionRead(d, meta)
}

func resourceTencentCloudClbRedirectionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_redirection.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	rewriteId := d.Id()

	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var instance *map[string]string
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := clbService.DescribeRedirectionById(ctx, rewriteId)
		if e != nil {
			return retryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read clb redirection failed, reason:%s\n ", logId, err.Error())
		return err
	}
	d.Set("clb_id", (*instance)["clb_id"])
	d.Set("source_listener_id", (*instance)["source_listener_id"])
	d.Set("target_listener_id", (*instance)["target_listener_id"])
	d.Set("source_rule_id", (*instance)["source_rule_id"])
	d.Set("target_rule_id", (*instance)["target_rule_id"])

	return nil
}

func resourceTencentCloudClbRedirectionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_redirection.delete")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	clbId := d.Id()
	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := clbService.DeleteRedirectionById(ctx, clbId)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete clb redirection failed, reason:%s\n ", logId, err.Error())
		return err
	}
	return nil
}