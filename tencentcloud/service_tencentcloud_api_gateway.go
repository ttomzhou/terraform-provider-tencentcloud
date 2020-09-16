package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type APIGatewayService struct {
	client *connectivity.TencentCloudClient
}

func (me *APIGatewayService) CreateApiKey(ctx context.Context, secretName string) (accessKeyId string, errRet error) {
	request := apigateway.NewCreateApiKeyRequest()
	request.SecretName = &secretName
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().CreateApiKey(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response.Result == nil || response.Response.Result.AccessKeyId == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty AccessKeyId", request.GetAction())
		return
	}
	accessKeyId = *response.Response.Result.AccessKeyId
	return
}

func (me *APIGatewayService) EnableApiKey(ctx context.Context, accessKeyId string) (errRet error) {
	request := apigateway.NewEnableApiKeyRequest()
	request.AccessKeyId = &accessKeyId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().EnableApiKey(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	if *response.Response.Result {
		return
	}
	return fmt.Errorf("enable api key fail")
}

func (me *APIGatewayService) DisableApiKey(ctx context.Context, accessKeyId string) (errRet error) {
	request := apigateway.NewDisableApiKeyRequest()
	request.AccessKeyId = &accessKeyId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().DisableApiKey(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	if *response.Response.Result {
		return
	}
	return fmt.Errorf("disable api key fail")
}

func (me *APIGatewayService) DescribeApiKey(ctx context.Context,
	accessKeyId string) (apiKey *apigateway.ApiKey, has bool, errRet error) {
	apiKeySet, err := me.DescribeApiKeysStatus(ctx, "", accessKeyId)
	if err != nil {
		errRet = err
		return
	}
	if len(apiKeySet) == 0 {
		return
	}
	has = true
	apiKey = apiKeySet[0]
	return
}

func (me *APIGatewayService) DescribeApiKeysStatus(ctx context.Context, secretName, accessKeyId string) (apiKeySet []*apigateway.ApiKey, errRet error) {
	request := apigateway.NewDescribeApiKeysStatusRequest()
	if secretName != "" || accessKeyId != "" {
		request.Filters = make([]*apigateway.Filter, 0, 2)
		if secretName != "" {
			request.Filters = append(request.Filters, &apigateway.Filter{Name: helper.String("SecretName"),
				Values: []*string{
					&secretName,
				}})
		}
		if accessKeyId != "" {
			request.Filters = append(request.Filters, &apigateway.Filter{Name: helper.String("AccessKeyId"),
				Values: []*string{
					&accessKeyId,
				}})
		}
	}

	var limit int64 = 20
	var offset int64 = 0

	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseAPIGatewayClient().DescribeApiKeysStatus(request)
		if err != nil {
			errRet = err
			return
		}
		if response.Response.Result == nil {
			errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
			return
		}
		if len(response.Response.Result.ApiKeySet) > 0 {
			apiKeySet = append(apiKeySet, response.Response.Result.ApiKeySet...)
		}
		if len(response.Response.Result.ApiKeySet) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *APIGatewayService) DeleteApiKey(ctx context.Context, accessKeyId string) (errRet error) {
	request := apigateway.NewDeleteApiKeyRequest()
	request.AccessKeyId = &accessKeyId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().DeleteApiKey(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	if *response.Response.Result {
		return
	}
	return fmt.Errorf("delete api key fail")
}

func (me *APIGatewayService) CreateUsagePlan(ctx context.Context, usagePlanName string, usagePlanDesc *string,
	maxRequestNum, maxRequestNumPreSec int64) (usagePlanId string, errRet error) {

	logId := getLogId(ctx)

	request := apigateway.NewCreateUsagePlanRequest()
	request.UsagePlanName = &usagePlanName
	request.MaxRequestNum = &maxRequestNum
	request.MaxRequestNumPreSec = &maxRequestNumPreSec
	if nil != usagePlanDesc {
		request.UsagePlanDesc = usagePlanDesc
	}

	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := me.client.UseAPIGatewayClient().CreateUsagePlan(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			return resource.RetryableError(err)
		}
		if response.Response.Result == nil {
			return resource.RetryableError(fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction()))
		}
		usagePlanId = *response.Response.Result.UsagePlanId
		return nil
	})

	return
}

func (me *APIGatewayService) DescribeUsagePlan(ctx context.Context, usagePlanId string) (info apigateway.UsagePlanInfo, has bool, errRet error) {

	request := apigateway.NewDescribeUsagePlanRequest()
	request.UsagePlanId = &usagePlanId

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAPIGatewayClient().DescribeUsagePlan(request)
	if err != nil {
		if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok && sdkErr.GetCode() == "ResourceNotFound.InvalidUsagePlan" {
			return
		}
		errRet = err
		return
	}
	if response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	has = true
	info = *response.Response.Result
	return
}

func (me *APIGatewayService) DeleteUsagePlan(ctx context.Context, usagePlanId string) (errRet error) {

	request := apigateway.NewDeleteUsagePlanRequest()
	request.UsagePlanId = &usagePlanId

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAPIGatewayClient().DeleteUsagePlan(request)

	if err != nil {
		return err
	}
	if response.Response.Result == nil {
		return fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
	}

	if !*response.Response.Result {
		return fmt.Errorf("delete usage plan fail")
	}

	return
}

func (me *APIGatewayService) ModifyUsagePlan(ctx context.Context,
	usagePlanId string,
	usagePlanName string,
	usagePlanDesc *string,
	maxRequestNum,
	maxRequestNumPreSec int64) (errRet error) {

	request := apigateway.NewModifyUsagePlanRequest()
	request.UsagePlanId = &usagePlanId

	ratelimit.Check(request.GetAction())
	request.UsagePlanName = &usagePlanName
	if usagePlanDesc != nil {
		request.UsagePlanDesc = usagePlanDesc
	}
	request.MaxRequestNum = &maxRequestNum
	request.MaxRequestNumPreSec = &maxRequestNumPreSec

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAPIGatewayClient().ModifyUsagePlan(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}

	return nil
}

func (me *APIGatewayService) DescribeUsagePlanEnvironments(ctx context.Context,
	usagePlanId string, bindType string) (list []*apigateway.UsagePlanEnvironment, errRet error) {

	request := apigateway.NewDescribeUsagePlanEnvironmentsRequest()
	request.UsagePlanId = &usagePlanId
	request.BindType = &bindType

	var limit int64 = 20
	var offset int64 = 0

	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseAPIGatewayClient().DescribeUsagePlanEnvironments(request)
		if err != nil {
			errRet = err
			return
		}
		if response.Response.Result == nil {
			errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
			return
		}
		if len(response.Response.Result.EnvironmentList) > 0 {
			list = append(list, response.Response.Result.EnvironmentList...)
		}
		if len(response.Response.Result.EnvironmentList) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *APIGatewayService) DescribeUsagePlansStatus(ctx context.Context,
	usagePlanId string, usagePlanName string) (infos []*apigateway.UsagePlanStatusInfo, errRet error) {

	request := apigateway.NewDescribeUsagePlansStatusRequest()

	if usagePlanId != "" || usagePlanName != "" {
		request.Filters = make([]*apigateway.Filter, 0, 2)
		if usagePlanId != "" {
			request.Filters = append(request.Filters, &apigateway.Filter{Name: helper.String("UsagePlanId"),
				Values: []*string{
					&usagePlanId,
				}})
		}
		if usagePlanName != "" {
			request.Filters = append(request.Filters, &apigateway.Filter{Name: helper.String("UsagePlanName"),
				Values: []*string{
					&usagePlanName,
				}})
		}
	}

	var limit int64 = 20
	var offset int64 = 0

	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseAPIGatewayClient().DescribeUsagePlansStatus(request)
		if err != nil {
			errRet = err
			return
		}
		if response.Response.Result == nil {
			errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
			return
		}
		if len(response.Response.Result.UsagePlanStatusSet) > 0 {
			infos = append(infos, response.Response.Result.UsagePlanStatusSet...)
		}
		if len(response.Response.Result.UsagePlanStatusSet) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *APIGatewayService) DescribeIPStrategysStatus(ctx context.Context,
	serviceId, strategyName string) (infos []*apigateway.IPStrategy, errRet error) {

	request := apigateway.NewDescribeIPStrategysStatusRequest()
	request.ServiceId = &serviceId

	if strategyName != "" {
		request.Filters = make([]*apigateway.Filter, 0, 1)
		if strategyName != "" {
			request.Filters = append(request.Filters, &apigateway.Filter{Name: helper.String("StrategyName"),
				Values: []*string{
					&strategyName,
				}})
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().DescribeIPStrategysStatus(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	if len(response.Response.Result.StrategySet) > 0 {
		infos = append(infos, response.Response.Result.StrategySet...)
	}
	return
}

func (me *APIGatewayService) DescribeIPStrategies(ctx context.Context, serviceId, strategyId, environmentName string) (ipStrategies *apigateway.IPStrategy, errRet error) {
	request := apigateway.NewDescribeIPStrategyRequest()

	request.ServiceId = &serviceId
	request.StrategyId = &strategyId
	request.EnvironmentName = &environmentName

	var (
		limit   int64 = 20
		offset  int64 = 0
		apiList []*apigateway.DesApisStatus
	)

	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseAPIGatewayClient().DescribeIPStrategy(request)
		if err != nil {
			errRet = err
			return
		}
		if response.Response.Result == nil {
			errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
			return
		}
		if len(response.Response.Result.BindApis) > 0 {
			apiList = append(apiList, response.Response.Result.BindApis...)
		}
		if len(response.Response.Result.BindApis) < int(limit) {
			ipStrategies = response.Response.Result
			ipStrategies.BindApis = apiList
			return
		}
		offset += limit
	}
}

func (me *APIGatewayService) BindSecretId(ctx context.Context,
	usagePlanId string, apiKeyId string) (errRet error) {

	request := apigateway.NewBindSecretIdsRequest()
	request.UsagePlanId = &usagePlanId
	request.AccessKeyIds = []*string{&apiKeyId}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().BindSecretIds(request)

	if err != nil {
		return err
	}
	if response.Response.Result == nil {
		return fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
	}

	if !*response.Response.Result {
		return fmt.Errorf("bind api key to usage plan fail")
	}

	return
}

func flattenOauthConfigMappings(v *apigateway.OauthConfig) map[string]interface{} {
	if v != nil {
		return map[string]interface{}{
			"login_redirect_url": *v.LoginRedirectUrl,
			"public_key":         *v.PublicKey,
			"token_location":     *v.TokenLocation,
		}
	}
	return nil
}

func (me *APIGatewayService) UnBindSecretId(ctx context.Context,
	usagePlanId string,
	apiKeyId string) (errRet error) {
	request := apigateway.NewUnBindSecretIdsRequest()
	request.UsagePlanId = &usagePlanId
	request.AccessKeyIds = []*string{&apiKeyId}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().UnBindSecretIds(request)

	if err != nil {
		return err
	}
	if response.Response.Result == nil {
		return fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
	}

	if !*response.Response.Result {
		return fmt.Errorf("unbind api key to usage plan fail")
	}

	return
}

func (me *APIGatewayService) CreateService(ctx context.Context,
	serviceName,
	protocol,
	serviceDesc,
	exclusiveSetName,
	ipVersion,
	setServerName,
	appidType string,
	netTypes []string) (serviceId string, errRet error) {

	request := apigateway.NewCreateServiceRequest()
	request.ServiceName = &serviceName
	request.Protocol = &protocol
	if serviceDesc != "" {
		request.ServiceDesc = &serviceDesc
	}
	if exclusiveSetName != "" {
		request.ExclusiveSetName = &exclusiveSetName
	}
	if ipVersion != "" {
		request.IpVersion = &ipVersion
	}
	if appidType != "" {
		request.AppIdType = &appidType
	}
	if setServerName != "" {
		request.SetServerName = &setServerName
	}
	request.NetTypes = helper.Strings(netTypes)

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().CreateService(request)

	if err != nil {
		errRet = err
		return
	}
	serviceId = *response.Response.ServiceId
	return
}

func (me *APIGatewayService) DescribeService(ctx context.Context, serviceId string) (info apigateway.DescribeServiceResponse, has bool, errRet error) {

	request := apigateway.NewDescribeServiceRequest()
	request.ServiceId = &serviceId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().DescribeService(request)
	if err != nil {
		if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok && sdkErr.GetCode() == "ResourceNotFound.InvalidService" {
			return
		}
		errRet = err
		return
	}
	info = *response
	has = true
	return
}

func (me *APIGatewayService) ModifyService(ctx context.Context,
	serviceId,
	serviceName,
	protocol,
	serviceDesc string,
	netTypes []string) (errRet error) {

	request := apigateway.NewModifyServiceRequest()
	request.ServiceId = &serviceId
	request.ServiceName = &serviceName
	request.Protocol = &protocol
	request.ServiceDesc = &serviceDesc
	request.NetTypes = helper.Strings(netTypes)

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseAPIGatewayClient().ModifyService(request)
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *APIGatewayService) DeleteService(ctx context.Context,
	serviceId string) (errRet error) {

	request := apigateway.NewDeleteServiceRequest()
	request.ServiceId = &serviceId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().DeleteService(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response.Result == nil {
		return fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
	}

	if !*response.Response.Result {
		return fmt.Errorf("delete service fail")
	}

	return
}

func (me *APIGatewayService) UnReleaseService(ctx context.Context,
	serviceId string,
	environment string) (errRet error) {

	request := apigateway.NewUnReleaseServiceRequest()
	request.ServiceId = &serviceId
	request.EnvironmentName = &environment

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().UnReleaseService(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response.Result == nil {
		return fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
	}

	if !*response.Response.Result {
		return fmt.Errorf("unrelease service %s.%s fail", serviceId, environment)
	}
	return
}

func (me *APIGatewayService) DescribeServiceUsagePlan(ctx context.Context,
	serviceId string) (list []*apigateway.ApiUsagePlan, errRet error) {

	request := apigateway.NewDescribeServiceUsagePlanRequest()
	request.ServiceId = &serviceId

	var limit int64 = 20
	var offset int64 = 0

	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseAPIGatewayClient().DescribeServiceUsagePlan(request)
		if err != nil {
			errRet = err
			return
		}
		if response.Response.Result == nil {
			errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
			return
		}
		if len(response.Response.Result.ServiceUsagePlanList) > 0 {
			list = append(list, response.Response.Result.ServiceUsagePlanList...)
		}
		if len(response.Response.Result.ServiceUsagePlanList) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *APIGatewayService) DescribeApiUsagePlan(ctx context.Context,
	serviceId string) (list []*apigateway.ApiUsagePlan, errRet error) {

	request := apigateway.NewDescribeApiUsagePlanRequest()
	request.ServiceId = &serviceId

	var limit int64 = 20
	var offset int64 = 0

	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseAPIGatewayClient().DescribeApiUsagePlan(request)
		if err != nil {
			errRet = err
			return
		}
		if response.Response.Result == nil {
			errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
			return
		}
		if len(response.Response.Result.ApiUsagePlanList) > 0 {
			list = append(list, response.Response.Result.ApiUsagePlanList...)
		}
		if len(response.Response.Result.ApiUsagePlanList) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *APIGatewayService) BindEnvironment(ctx context.Context,
	serviceId, environment, bindType, usagePlanId, apiId string) (errRet error) {

	logId := getLogId(ctx)
	request := apigateway.NewBindEnvironmentRequest()
	request.ServiceId = &serviceId
	request.UsagePlanIds = []*string{&usagePlanId}
	request.Environment = &environment
	request.BindType = &bindType

	if bindType == API_GATEWAY_TYPE_API {
		request.ApiIds = []*string{&apiId}
	}

	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := me.client.UseAPIGatewayClient().BindEnvironment(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			return retryError(err)
		}

		if response.Response.Result == nil {
			return resource.RetryableError(fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction()))
		}

		if !*response.Response.Result {
			return resource.RetryableError(fmt.Errorf("%s attach to %s.%s fail", usagePlanId, serviceId, apiId))
		}

		return nil
	})

	return
}

func (me *APIGatewayService) UnBindEnvironment(ctx context.Context,
	serviceId, environment, bindType, usagePlanId, apiId string) (errRet error) {

	request := apigateway.NewUnBindEnvironmentRequest()
	request.ServiceId = &serviceId
	request.UsagePlanIds = []*string{&usagePlanId}
	request.Environment = &environment
	request.BindType = &bindType

	if bindType == API_GATEWAY_TYPE_API {
		request.ApiIds = []*string{&apiId}
	}

	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, errRet := me.client.UseAPIGatewayClient().UnBindEnvironment(request)
		if errRet != nil {
			return retryError(errRet)
		}

		if response.Response.Result == nil {
			return resource.RetryableError(fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction()))
		}

		if !*response.Response.Result {
			return resource.RetryableError(fmt.Errorf("%s unattach to %s.%s fail", usagePlanId, serviceId, apiId))
		}

		return nil
	})

	return
}

func (me *APIGatewayService) DescribeApi(ctx context.Context,
	serviceId,
	apiId string) (info apigateway.ApiInfo, has bool, errRet error) {

	request := apigateway.NewDescribeApiRequest()
	request.ServiceId = &serviceId
	request.ApiId = &apiId

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAPIGatewayClient().DescribeApi(request)
	if err != nil {
		if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok &&
			(sdkErr.GetCode() == "ResourceNotFound.InvalidService" || sdkErr.GetCode() == "ResourceNotFound.InvalidApi") {
			return
		}
		errRet = err
		return
	}

	if response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}

	has = true
	info = *response.Response.Result
	return
}

func (me *APIGatewayService) DeleteApi(ctx context.Context, serviceId,
	apiId string) (errRet error) {
	request := apigateway.NewDeleteApiRequest()
	request.ServiceId = &serviceId
	request.ApiId = &apiId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().DeleteApi(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	if *response.Response.Result {
		return
	}
	return fmt.Errorf("delete api fail")
}

func (me *APIGatewayService) DescribeServicesStatus(ctx context.Context,
	serviceId,
	serviceName string) (infos []*apigateway.Service, errRet error) {

	request := apigateway.NewDescribeServicesStatusRequest()

	if serviceId != "" || serviceName != "" {
		request.Filters = make([]*apigateway.Filter, 0, 2)
		if serviceId != "" {
			request.Filters = append(request.Filters, &apigateway.Filter{Name: helper.String("ServiceId"),
				Values: []*string{
					&serviceId,
				}})
		}
		if serviceName != "" {
			request.Filters = append(request.Filters, &apigateway.Filter{Name: helper.String("ServiceName"),
				Values: []*string{
					&serviceName,
				}})
		}
	}

	var limit int64 = 20
	var offset int64 = 0

	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseAPIGatewayClient().DescribeServicesStatus(request)
		if err != nil {
			errRet = err
			return
		}
		if response.Response.Result == nil {
			errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
			return
		}
		if len(response.Response.Result.ServiceSet) > 0 {
			infos = append(infos, response.Response.Result.ServiceSet...)
		}
		if len(response.Response.Result.ServiceSet) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *APIGatewayService) DescribeApisStatus(ctx context.Context,
	serviceId, apiName, apiId string) (infos []*apigateway.DesApisStatus, errRet error) {

	request := apigateway.NewDescribeApisStatusRequest()
	request.ServiceId = &serviceId

	if apiId != "" || apiName != "" {
		request.Filters = make([]*apigateway.Filter, 0, 2)
		if apiId != "" {
			request.Filters = append(request.Filters, &apigateway.Filter{Name: helper.String("ApiId"),
				Values: []*string{
					&apiId,
				}})
		}
		if apiName != "" {
			request.Filters = append(request.Filters, &apigateway.Filter{Name: helper.String("ApiName"),
				Values: []*string{
					&apiName,
				}})
		}
	}

	var limit int64 = 20
	var offset int64 = 0

	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseAPIGatewayClient().DescribeApisStatus(request)
		if err != nil {
			errRet = err
			return
		}
		if response.Response.Result == nil {
			errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
			return
		}
		if len(response.Response.Result.ApiIdStatusSet) > 0 {
			infos = append(infos, response.Response.Result.ApiIdStatusSet...)
		}
		if len(response.Response.Result.ApiIdStatusSet) < int(limit) {
			return
		}
		offset += limit
	}
}
