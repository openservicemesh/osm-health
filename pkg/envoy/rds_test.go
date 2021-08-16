package envoy

import (
	"testing"

	tassert "github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestEnvoyOutboundRouteDomainChecker(t *testing.T) {
	assert := tassert.New(t)
	configGetter := mockConfigGetter{
		getter: createConfigGetterFunc("../../tests/sample-envoy-config-dump-bookbuyer.json"),
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bookstore",
			Namespace: "bookstore",
		},
	}
	routeDomainChecker := HasOutboundDynamicRouteConfigDomainCheck(configGetter, pod)
	checkError := routeDomainChecker.Run()
	assert.Nil(checkError)
}

func TestEnvoyOutboundRouteDomainCheckerEmptyConfig(t *testing.T) {
	assert := tassert.New(t)
	configGetter := mockConfigGetter{
		getter: func() (*Config, error) {
			return nil, nil
		},
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bookstore",
			Namespace: "bookstore",
		},
	}
	routeDomainChecker := HasOutboundDynamicRouteConfigDomainCheck(configGetter, pod)
	checkError := routeDomainChecker.Run()
	assert.NotNil(checkError)
	assert.Equal("envoy config is empty", checkError.Error())
}

func TestEnvoyOutboundRouteDomainCheckerNoDomains(t *testing.T) {
	assert := tassert.New(t)
	configGetter := mockConfigGetter{
		getter: createConfigGetterFunc("../../tests/sample-envoy-config-dump-bookbuyer-no-rds-dynamic-route-virtual-host-domains.json"),
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bookstore",
			Namespace: "bookstore",
		},
	}
	routeDomainChecker := HasOutboundDynamicRouteConfigDomainCheck(configGetter, pod)
	checkError := routeDomainChecker.Run()
	assert.NotNil(checkError)
	assert.Equal("no dynamic route config domains", checkError.Error())
}

func TestEnvoyOutboundRouteDomainCheckerDomainNotFound(t *testing.T) {
	assert := tassert.New(t)
	configGetter := mockConfigGetter{
		getter: createConfigGetterFunc("../../tests/sample-envoy-config-dump-bookbuyer-not-found-rds-dynamic-route-virtual-host-domain.json"),
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bookstore",
			Namespace: "bookstore",
		},
	}
	routeDomainChecker := HasOutboundDynamicRouteConfigDomainCheck(configGetter, pod)
	checkError := routeDomainChecker.Run()
	assert.NotNil(checkError)
	assert.Equal("dynamic route config domain not found", checkError.Error())
}

func TestEnvoyInboundRouteDomainChecker(t *testing.T) {
	assert := tassert.New(t)
	configGetter := mockConfigGetter{
		getter: createConfigGetterFunc("../../tests/sample-envoy-config-dump-bookstore.json"),
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bookstore-v1",
			Namespace: "bookstore",
		},
	}
	routeDomainChecker := HasInboundDynamicRouteConfigDomainCheck(configGetter, pod)
	checkError := routeDomainChecker.Run()
	assert.Nil(checkError)
}

func TestEnvoyInboundRouteDomainCheckerEmptyConfig(t *testing.T) {
	assert := tassert.New(t)
	configGetter := mockConfigGetter{
		getter: func() (*Config, error) {
			return nil, nil
		},
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bookstore-v1",
			Namespace: "bookstore",
		},
	}
	routeDomainChecker := HasInboundDynamicRouteConfigDomainCheck(configGetter, pod)
	checkError := routeDomainChecker.Run()
	assert.NotNil(checkError)
	assert.Equal("envoy config is empty", checkError.Error())
}

func TestEnvoyInboundRouteDomainCheckerNoDomains(t *testing.T) {
	assert := tassert.New(t)
	configGetter := mockConfigGetter{
		getter: createConfigGetterFunc("../../tests/sample-envoy-config-dump-bookstore-no-rds-dynamic-route-virtual-host-domains.json"),
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bookstore-v1",
			Namespace: "bookstore",
		},
	}
	routeDomainChecker := HasInboundDynamicRouteConfigDomainCheck(configGetter, pod)
	checkError := routeDomainChecker.Run()
	assert.NotNil(checkError)
	assert.Equal("no dynamic route config domains", checkError.Error())
}

func TestEnvoyInboundRouteDomainCheckerDomainNotFound(t *testing.T) {
	assert := tassert.New(t)
	configGetter := mockConfigGetter{
		getter: createConfigGetterFunc("../../tests/sample-envoy-config-dump-bookstore-not-found-rds-dynamic-route-virtual-host-domain.json"),
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bookstore-v1",
			Namespace: "bookstore",
		},
	}
	routeDomainChecker := HasInboundDynamicRouteConfigDomainCheck(configGetter, pod)
	checkError := routeDomainChecker.Run()
	assert.NotNil(checkError)
	assert.Equal("dynamic route config domain not found", checkError.Error())
}
