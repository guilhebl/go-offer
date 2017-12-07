package model

type ServiceEndpoint struct {
	URL string
}

func NewServiceEndpoint(url string) *ServiceEndpoint {
	e := &ServiceEndpoint{
		URL: url,
	}
	return e
}
