/*
Use this data source to query detailed information of audits.

Example Usage

```hcl
data "tencentcloud_audits" "audits" {
  name       = "test"
}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	audit "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func dataSourceTencentCloudAudits() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAuditsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the audits.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values
			"audit_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information list of the dedicated audits.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the audit.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the audit.",
						},
						"cos_bucket": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cos bucket name where audit save logs.",
						},
						"log_file_prefix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Prefix of the log file of the audit.",
						},
						"audit_switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicate whether audit start logging or not.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAuditsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_audits.read")()

	logId := getLogId(contextNil)

	name := d.Get("name").(string)
	request := audit.NewListAuditsRequest()

	var response *audit.ListAuditsResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAuditClient().ListAudits(request)
		if e != nil {
			log.Printf("[CRITAL]%s %s fail, reason:%s\n", logId, request.GetAction(), e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		return err
	}

	result := response.Response.AuditSummarys
	ids := make([]string, 0, len(result))
	auditList := make([]map[string]interface{}, 0, len(result))
	for _, audit := range result {
		if name != "" && name != *audit.AuditName {
			continue
		}
		mapping := map[string]interface{}{
			"id":              audit.AuditName,
			"name":            audit.AuditName,
			"audit_switch":    int(*audit.AuditStatus) > 0,
			"log_file_prefix": *audit.LogFilePrefix,
			"cos_bucket":      audit.CosBucketName,
		}

		auditList = append(auditList, mapping)
		ids = append(ids, *audit.AuditName)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("audit_list", auditList); e != nil {
		log.Printf("[CRITAL]%s provider set audit list fail, reason:%s\n", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), auditList); e != nil {
			return e
		}
	}

	return nil

}
