resource "tencentcloud_api_gateway_usage_plan" "plan" {
  	usage_plan_name         = "my_plan"
  	usage_plan_desc         = "nice plan"
  	max_request_num         = 100
  	max_request_num_pre_sec = 10
}

resource "tencentcloud_api_gateway_usage_plan_attachment" "attach_service" {
  	usage_plan_id  = tencentcloud_api_gateway_usage_plan.plan.id
  	service_id     = "service-ke4o2arm"
  	environment    = "test"
  	bind_type      = "SERVICE"
}