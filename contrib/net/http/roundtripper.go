package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/DataDog/dd-trace-go/tracer"
	"github.com/DataDog/dd-trace-go/tracer/ext"
)

const defaultResourceName = "http.request"

type roundTripper struct {
	base           http.RoundTripper
	service        string
	resourceMapper func(*http.Request) string
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	t := tracer.DefaultTracer
	if !t.Enabled() {
		return rt.base.RoundTrip(req)
	}

	span := tracer.NewChildSpanFromContext("http.request", req.Context())
	defer span.Finish()
	span.Resource = rt.resourceMapper(req)
	if span.Resource == "" {
		span.Resource = defaultResourceName
	}
	span.Service = rt.service
	span.Type = ext.HTTPType
	span.SetMeta(ext.HTTPMethod, req.Method)
	span.SetMeta(ext.HTTPURL, req.URL.Path)
	res, err := rt.base.RoundTrip(req)
	if err != nil {
		span.FinishWithErr(err)
		return res, err
	}
	span.SetMeta(ext.HTTPCode, strconv.Itoa(res.StatusCode))
	if res.StatusCode/100 == 5 {
		span.FinishWithErr(errors.New(res.Status))
	}
	return res, err
}

// WrapRoundTripper returns a new RoundTripper which traces all requests sent
func WrapRoundTripper(
	rt http.RoundTripper,
	service string,
	resourceMapper func(*http.Request) string,
) http.RoundTripper {
	if resourceMapper == nil {
		resourceMapper = func(*http.Request) string {
			return defaultResourceName
		}
	}
	return &roundTripper{
		base:           rt,
		service:        service,
		resourceMapper: resourceMapper,
	}
}
