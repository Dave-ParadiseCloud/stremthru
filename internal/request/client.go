package request

import (
	"net/http"
	"time"
)

var DefaultHTTPTransport = func() *http.Transport {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DisableKeepAlives = true
	return transport
}()

var DefaultHTTPClient = func() *http.Client {
	return &http.Client{
		Transport: DefaultHTTPTransport,
		Timeout:   90 * time.Second,
	}
}()

func GetHTTPClient(withProxy bool) *http.Client {
	transport := DefaultHTTPTransport.Clone()
	if !withProxy {
		transport.Proxy = nil
	}
	return &http.Client{
		Transport: transport,
		Timeout:   90 * time.Second,
	}
}

type APIResponse[T interface{}] struct {
	Header     http.Header
	StatusCode int
	Data       T
}

func NewAPIResponse[T interface{}](res *http.Response, data T) APIResponse[T] {
	apiResponse := APIResponse[T]{
		StatusCode: 503,
		Data:       data,
	}
	if res != nil {
		apiResponse.Header = res.Header
		apiResponse.StatusCode = res.StatusCode
	}
	return apiResponse
}
