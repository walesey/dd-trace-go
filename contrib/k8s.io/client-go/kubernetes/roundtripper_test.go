package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathToResource(t *testing.T) {
	expected := map[string]string{
		"/api/v1/componentstatuses":                                           "componentstatuses",
		"/api/v1/componentstatuses/NAME":                                      "componentstatuses/{name}",
		"/api/v1/configmaps":                                                  "configmaps",
		"/api/v1/namespaces/default/bindings":                                 "namespaces/{namespace}/bindings",
		"/api/v1/namespaces/someothernamespace/configmaps":                    "namespaces/{namespace}/configmaps",
		"/api/v1/namespaces/default/configmaps/some-config-map":               "namespaces/{namespace}/configmaps/{name}",
		"/api/v1/namespaces/default/persistentvolumeclaims/pvc-abcd/status":   "namespaces/{namespace}/persistentvolumeclaims/{name}/status",
		"/api/v1/namespaces/default/pods/pod-1234/proxy":                      "namespaces/{namespace}/pods/{name}/proxy",
		"/api/v1/namespaces/default/pods/pod-5678/proxy/some-path":            "namespaces/{namespace}/pods/{name}/proxy/{path}",
		"/api/v1/watch/configmaps":                                            "watch/configmaps",
		"/api/v1/watch/namespaces":                                            "watch/namespaces",
		"/api/v1/watch/namespaces/default/configmaps":                         "watch/namespaces/{namespace}/configmaps",
		"/api/v1/watch/namespaces/someothernamespace/configmaps/another-name": "watch/namespaces/{namespace}/configmaps/{name}",
	}

	for path, expectedResource := range expected {
		assert.Equal(t, expectedResource, pathToResource(path), "mapping %v", path)
	}
}
