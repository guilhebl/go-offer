package amazon

// Endpoints are the Amazon API endpoints by region
var Endpoints = map[string]string{
	"BR": "webservices.amazon.com.br",
	"CA": "webservices.amazon.ca",
	"CN": "webservices.amazon.cn",
	"DE": "webservices.amazon.de",
	"ES": "webservices.amazon.es",
	"FR": "webservices.amazon.fr",
	"IN": "webservices.amazon.in",
	"IT": "webservices.amazon.it",
	"JP": "webservices.amazon.co.jp",
	"MX": "webservices.amazon.com.mx",
	"UK": "webservices.amazon.co.uk",
	"US": "webservices.amazon.com",
}

// EndpointURI is the fixed request URI of the API
const EndpointURI = "/onca/xml"

// Config describes the service configuration
type Config struct {
	AccessKey    string
	AccessSecret string
	AssociateTag string
	Region       string
	Secure       bool
}

func NewConfig(accessKeyId, secretKey, associateTag, region string, secure bool) Config {
	cfg := Config{
		AccessKey:    accessKeyId,
		AccessSecret: secretKey,
		AssociateTag: associateTag,
		Region:       region,
		Secure:       secure,
	}
	return cfg
}
