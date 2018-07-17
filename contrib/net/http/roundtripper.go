package http

import (
	"net/http"
	"strconv"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const defaultResourceName = "http.request"

type roundTripper struct {
	base           http.RoundTripper
	service        string
	resourceMapper func(*http.Request) string
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	span, _ := tracer.StartSpanFromContext(req.Context(), "http.request",
		tracer.ServiceName(rt.service),
		tracer.ResourceName(rt.resourceMapper(req)),
		tracer.Tag(ext.HTTPMethod, req.Method),
		tracer.Tag(ext.HTTPURL, req.URL.Path),
	)
	defer span.Finish()
	res, err := rt.base.RoundTrip(req)
	if err != nil {
		span.SetTag("http.errors", err.Error())
		span.SetTag(ext.Error, err.Error())
		return res, err
	}
	span.SetTag(ext.HTTPCode, strconv.Itoa(res.StatusCode))
	if res.StatusCode/100 == 5 {
		span.SetTag("http.errors", err.Error())
		span.SetTag(ext.Error, err.Error())
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
