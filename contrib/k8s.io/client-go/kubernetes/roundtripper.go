package kubernetes

import (
	"net/http"
	"strings"

	tracehttp "github.com/DataDog/dd-trace-go/contrib/net/http"
)

const (
	apiPrefix        = "/api/v1/"
	watchPrefix      = "watch/"
	namespacesPrefix = "namespaces/"
)

// WrapRoundTripper wraps a RoundTripper intended for interfacing with
// Kubernetes and traces all requests
func WrapRoundTripper(rt http.RoundTripper) http.RoundTripper {
	return tracehttp.WrapRoundTripper(rt, "kubernetes", func(req *http.Request) string {
		return pathToResource(req.URL.Path)
	})
}

func pathToResource(path string) string {
	resourceName := ""

	if !strings.HasPrefix(path, apiPrefix) {
		return resourceName
	}
	path = path[len(apiPrefix):]

	// strip out /watch
	if strings.HasPrefix(path, watchPrefix) {
		path = path[len(watchPrefix):]
		resourceName += watchPrefix
	}

	// {type}/{name}
	lastType := ""
	for i := 0; ; i++ {
		idx := strings.IndexByte(path, '/')
		if i%2 == 0 {
			// parse {type}
			if idx < 0 {
				lastType = path
			} else {
				lastType = path[:idx]
			}
			resourceName += lastType
		} else {
			// parse {name}
			resourceName += typeToPlaceholder(lastType)
		}
		if idx < 0 {
			break
		}
		path = path[idx+1:]
		resourceName += "/"
	}
	return resourceName
}

func typeToPlaceholder(typ string) string {
	switch typ {
	case "namespaces":
		return "{namespace}"
	case "proxy":
		return "{path}"
	default:
		return "{name}"
	}
}
