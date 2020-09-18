package tencentcloud

const (
	API_GATEWAY_KEY_ENABLED  = "on"
	API_GATEWAY_KEY_DISABLED = "off"
)

var API_GATEWAY_KEYS = []string{
	API_GATEWAY_KEY_ENABLED,
	API_GATEWAY_KEY_DISABLED,
}
var API_GATEWAY_KEY_STR2INTS = map[string]int64{
	API_GATEWAY_KEY_ENABLED:  1,
	API_GATEWAY_KEY_DISABLED: 0,
}
var API_GATEWAY_KEY_INT2STRS = map[int64]string{
	1: API_GATEWAY_KEY_ENABLED,
	0: API_GATEWAY_KEY_DISABLED,
}

const (
	API_GATEWAY_TYPE_SERVICE = "SERVICE"
	API_GATEWAY_TYPE_API     = "API"
)

var API_GATEWAY_TYPES = []string{
	API_GATEWAY_TYPE_SERVICE,
	API_GATEWAY_TYPE_API,
}

const (
	API_GATEWAY_NET_TYPE_INNER = "INNER"
	API_GATEWAY_NET_TYPE_OUTER = "OUTER"
)

const (
	API_GATEWAY_NET_IP_VERSION4 = "IPv4"
	API_GATEWAY_NET_IP_VERSION6 = "IPv6"
)

var API_GATEWAY_NET_IP_VERSIONS = []string{
	API_GATEWAY_NET_IP_VERSION4,
	API_GATEWAY_NET_IP_VERSION6,
}

const (
	API_GATEWAY_SERVICE_PROTOCOL_HTTP  = "http"
	API_GATEWAY_SERVICE_PROTOCOL_HTTPS = "https"
	API_GATEWAY_SERVICE_PROTOCOL_ALL   = "http&https"
)

var API_GATEWAY_SERVICE_PROTOCOLS = []string{
	API_GATEWAY_SERVICE_PROTOCOL_HTTP,
	API_GATEWAY_SERVICE_PROTOCOL_HTTPS,
	API_GATEWAY_SERVICE_PROTOCOL_ALL,
}

const (
	API_GATEWAY_SERVICE_ENV_TEST    = "test"
	API_GATEWAY_SERVICE_ENV_RELEASE = "release"
	API_GATEWAY_SERVICE_ENV_PREPUB  = "prepub"
)

var API_GATEWAY_SERVICE_ENVS = []string{
	API_GATEWAY_SERVICE_ENV_TEST,
	API_GATEWAY_SERVICE_ENV_RELEASE,
	API_GATEWAY_SERVICE_ENV_PREPUB,
}

const (
	API_GATEWAY_SERVICE_TYPE_WEBSOCKET = "WEBSOCKET"
	API_GATEWAY_SERVICE_TYPE_HTTP      = "HTTP"
	API_GATEWAY_SERVICE_TYPE_SCF       = "SCF"
	API_GATEWAY_SERVICE_TYPE_MOCK      = "MOCK"
)

var API_GATEWAY_SERVICE_TYPES = []string{
	API_GATEWAY_SERVICE_TYPE_WEBSOCKET,
	API_GATEWAY_SERVICE_TYPE_HTTP,
	API_GATEWAY_SERVICE_TYPE_SCF,
	API_GATEWAY_SERVICE_TYPE_MOCK,
}

const (
	API_GATEWAY_AUTH_TYPE_SECRET = "SECRET"
	API_GATEWAY_AUTH_TYPE_NONE   = "NONE"
)

var API_GATEWAY_AUTH_TYPES = []string{
	API_GATEWAY_AUTH_TYPE_SECRET,
	API_GATEWAY_AUTH_TYPE_NONE,
}

const (
	API_GATEWAY_API_PROTOCOL_HTTP      = "HTTP"
	API_GATEWAY_API_PROTOCOL_WEBSOCKET = "WEBSOCKET"
)

var API_GATEWAY_API_PROTOCOLS = []string{
	API_GATEWAY_API_PROTOCOL_HTTP,
	API_GATEWAY_API_PROTOCOL_WEBSOCKET,
}

var API_GATEWAY_API_RESPONSE_TYPES = []string{"HTML", "JSON", "TEXT", "BINARY", "XML", ""}

const (
	CertificateIdExpired       = "FailedOperation.CertificateIdExpired"
	CertificateIdUnderVerify   = "FailedOperation.CertificateIdUnderVerify"
	DomainNeedBeian            = "FailedOperation.DomainNeedBeian"
	ExceededDefineMappingLimit = "LimitExceeded.ExceededDefineMappingLimit"
	DomainResolveError         = "FailedOperation.DomainResolveError"
)