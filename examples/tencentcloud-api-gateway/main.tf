resource "tencentcloud_api_gateway_service" "service" {
  service_name = "ck"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_api" "api" {
  service_id            = tencentcloud_api_gateway_service.service.id
  api_name              = "hello"
  api_desc              = "my hello api"
  auth_type             = "NONE"
  protocol              = "HTTP"
  enable_cors           = true
  request_config_path   = "/user/info"
  request_config_method = "GET"

  request_parameters {
    name          = "name"
    position      = "QUERY"
    type          = "string"
    desc          = "who are you?"
    default_value = "tom"
    required      = true
  }
  service_config_type      = "HTTP"
  service_config_timeout   = 15
  service_config_url       = "http://www.qq.com"
  service_config_path      = "/user"
  service_config_method    = "GET"
  response_type            = "HTML"
  response_success_example = "success"
  response_fail_example    = "fail"
  response_error_codes {
    code           = 100
    msg            = "system error"
    desc           = "system error code"
    converted_code = -100
    need_convert   = true
  }

}

resource "tencentcloud_api_gateway_api" "scf" {
  service_id                            = tencentcloud_api_gateway_service.service.id
  api_name                              = "scf_hello"
  api_desc                              = "my scf hello api"
  auth_type                             = "NONE"
  protocol                              = "HTTP"
  request_config_path                   = "/user/info2"
  request_config_method                 = "GET"
  service_config_type                   = "SCF"
  service_config_scf_function_name      = "resource-bot"
  service_config_scf_function_namespace = "default"
  service_config_scf_function_qualifier = "$LATEST"
}


resource "tencentcloud_api_gateway_api" "mock" {
  service_id                         = tencentcloud_api_gateway_service.service.id
  api_name                           = "mock_hello"
  api_desc                           = "my mock hello api"
  auth_type                          = "NONE"
  protocol                           = "HTTP"
  request_config_path                = "/user/mock"
  request_config_method              = "POST"
  service_config_type                = "MOCK"
  service_config_mock_return_message = "guaguajiao"
}

resource "tencentcloud_api_gateway_api" "websock" {
  service_id            = tencentcloud_api_gateway_service.service.id
  api_name              = "websock_hello"
  api_desc              = "my websock hello api"
  auth_type             = "NONE"
  protocol              = "WEBSOCKET"
  request_config_path   = "/user/websock"
  request_config_method = "GET"

  service_config_type    = "WEBSOCKET"
  service_config_timeout = 15
  service_config_url     = "ws://www.qq.com"
  service_config_path    = "/user"
  service_config_method  = "GET"
}

data "tencentcloud_api_gateway_apis" "id" {
  service_id = tencentcloud_api_gateway_service.service.id
  api_id     = tencentcloud_api_gateway_api.api.id
}

data "tencentcloud_api_gateway_apis" "name" {
  service_id = tencentcloud_api_gateway_service.service.id
  api_name   = tencentcloud_api_gateway_api.api.api_name
}

resource "tencentcloud_api_gateway_service" "service" {
  service_name = "niceservice"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

data "tencentcloud_api_gateway_services" "name" {
    service_name = tencentcloud_api_gateway_service.service.service_name
}

data "tencentcloud_api_gateway_services" "ids" {
    service_id = tencentcloud_api_gateway_service.service.id
}

resource "tencentcloud_api_gateway_custom_domain" "service" {
	service_id 			   = "service-ohxqslqe"
	sub_domain 			   = "tic-test.dnsv1.com"
	protocol   			   = "http"
	net_type   			   = "OUTER"
	is_default_mapping = "false"
	default_domain 	 	 = "service-ohxqslqe-1259649581.gz.apigw.tencentcs.com"
	path_mappings 		 = ["/good#test","/root#release"]
}

resource "tencentcloud_api_gateway_throttling_api" "service" {
	service_id       = "service-4r4xrrz4"
	strategy         = "400"
	environment_name = "test"
	api_ids          = ["api-lukm33yk"]
}

data "tencentcloud_api_gateway_throttling_apis" "id" {
    service_id = tencentcloud_api_gateway_throttling_api.service.service_id
}

data "tencentcloud_api_gateway_throttling_apis" "foo" {
	service_id 		    = tencentcloud_api_gateway_throttling_api.service.service_id
	environment_names = ["release", "test"]
}

resource "tencentcloud_api_gateway_throttling_service" "service" {
	service_id        = "service-4r4xrrz4"
	strategy          = "400"
	environment_names = ["release"]
}

data "tencentcloud_api_gateway_throttling_services" "id" {
    service_id = tencentcloud_api_gateway_throttling_service.service.service_id
}