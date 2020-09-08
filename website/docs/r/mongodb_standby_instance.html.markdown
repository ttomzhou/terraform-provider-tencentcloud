---
subcategory: "MongoDB"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_standby_instance"
sidebar_current: "docs-tencentcloud-resource-mongodb_standby_instance"
description: |-
  Provide a resource to create a Mongodb standby instance.
---

# tencentcloud_mongodb_standby_instance

Provide a resource to create a Mongodb standby instance.

## Example Usage

```hcl
provider "tencentcloud" {
  region = "ap-guangzhou"
}

provider "tencentcloud" {
  alias  = "shanghai"
  region = "ap-shanghai"
}

resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-test"
  memory         = 4
  volume         = 100
  engine_version = "MONGO_40_WT"
  machine_type   = "HIO10G"
  available_zone = var.availability_zone
  project_id     = 0
  password       = "test1234"

  tags = {
    test = "test"
  }
}

resource "tencentcloud_mongodb_standby_instance" "mongodb" {
  provider               = tencentcloud.shanghai
  instance_name          = "tf-mongodb-standby-test"
  memory                 = 4
  volume                 = 100
  available_zone         = "ap-shanghai-2"
  project_id             = 0
  father_instance_id     = tencentcloud_mongodb_instance.mongodb.id
  father_instance_region = "ap-guangzhou"
  charge_type            = "PREPAID"
  prepaid_period         = 1

  tags = {
    test = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `available_zone` - (Required, ForceNew) The available zone of the Mongodb standby instance. NOTE: must not same with father instance's.
* `father_instance_id` - (Required, ForceNew) Indicates the master instance ID of standby instances.
* `father_instance_region` - (Required, ForceNew) Indicates the region of father instance.
* `instance_name` - (Required) Name of the Mongodb instance.
* `memory` - (Required) Memory size. The minimum value is 2, and unit is GB. Memory and volume must be upgraded or degraded simultaneously.
* `volume` - (Required) Disk size. The minimum value is 25, and unit is GB. Memory and volume must be upgraded or degraded simultaneously.
* `auto_renew_flag` - (Optional) Auto renew flag. Valid values are `0`(NOTIFY_AND_MANUAL_RENEW), `1`(NOTIFY_AND_AUTO_RENEW) and `2`(DISABLE_NOTIFY_AND_MANUAL_RENEW). Default value is `0`. Note: only works for PREPAID instance. Only supports`0` and `1` for creation.
* `charge_type` - (Optional, ForceNew) The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. Default value is `POSTPAID_BY_HOUR`. Note: TencentCloud International only supports `POSTPAID_BY_HOUR`. Caution that update operation on this field will delete old instances and create new one with new charge type.
* `prepaid_period` - (Optional) The tenancy (time unit is month) of the prepaid instance. Valid values are 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36. NOTE: it only works when charge_type is set to `PREPAID`.
* `project_id` - (Optional) ID of the project which the instance belongs.
* `security_groups` - (Optional, ForceNew) ID of the security group. NOTE: for instance which `engine_version` is `MONGO_40_WT`, `security_groups` is not supported.
* `subnet_id` - (Optional, ForceNew) ID of the subnet within this VPC. The value is required if `vpc_id` is set.
* `tags` - (Optional) The tags of the Mongodb. Key name `project` is system reserved and can't be used.
* `vpc_id` - (Optional, ForceNew) ID of the VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of the Mongodb instance.
* `engine_version` - Version of the Mongodb and must be same as the master's.
* `machine_type` - Type of Mongodb instance and must be same as the master's.
* `status` - Status of the Mongodb instance, and available values include pending initialization(expressed with 0),  processing(expressed with 1), running(expressed with 2) and expired(expressed with -2).
* `vip` - IP of the Mongodb instance.
* `vport` - IP port of the Mongodb instance.


## Import

Mongodb instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_mongodb_standby_instance.mongodb cmgo-41s6jwy4
```

