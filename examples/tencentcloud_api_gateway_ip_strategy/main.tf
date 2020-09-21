resource "tencentcloud_api_gateway_ip_strategy" "test"{
	service_id 	  = "service-ohxqslqe"
	strategy_name = "tf_test"
	strategy_type = "BLACK"
	strategy_data = "9.9.9.9"
}

data "tencentcloud_api_gateway_ip_strategies" "id" {
	service_id = tencentcloud_api_gateway_ip_strategy.test.service_id
}