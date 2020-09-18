provider "tencentcloud" {
  region = "ap-guangzhou"
}

resource "tencentcloud_api_gateway_usage_plan" "plan" {
  	usage_plan_name         = "my_plan"
  	usage_plan_desc         = "nice plan"
  	max_request_num         = 100
  	max_request_num_pre_sec = 10
}
