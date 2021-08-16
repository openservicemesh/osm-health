package envoy

import (
	"testing"

	tassert "github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	bookstoreDestinationHost = "bookstore.bookstore"
)

func TestEnvoyOutboundRouteDomainPodChecker(t *testing.T) {
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
	routeDomainChecker := HasOutboundDynamicRouteConfigDomainPodCheck(configGetter, pod)
	checkError := routeDomainChecker.Run()
	assert.Nil(checkError)
}

func TestEnvoyOutboundRouteDomainPodCheckerEmptyConfig(t *testing.T) {
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
	routeDomainChecker := HasOutboundDynamicRouteConfigDomainPodCheck(configGetter, pod)
	checkError := routeDomainChecker.Run()
	assert.NotNil(checkError)
	assert.Equal("envoy config is empty", checkError.Error())
}

func TestEnvoyOutboundRouteDomainPodCheckerNoDomains(t *testing.T) {
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
	routeDomainChecker := HasOutboundDynamicRouteConfigDomainPodCheck(configGetter, pod)
	checkError := routeDomainChecker.Run()
	assert.NotNil(checkError)
	assert.Equal("no dynamic route config domains", checkError.Error())
}

func TestEnvoyOutboundRouteDomainPodCheckerDomainNotFound(t *testing.T) {
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
	routeDomainChecker := HasOutboundDynamicRouteConfigDomainPodCheck(configGetter, pod)
	checkError := routeDomainChecker.Run()
	assert.NotNil(checkError)
	assert.Equal("dynamic route config domain not found", checkError.Error())
}

func TestEnvoyInboundRouteDomainPodChecker(t *testing.T) {
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
	routeDomainChecker := HasInboundDynamicRouteConfigDomainPodCheck(configGetter, pod)
	checkError := routeDomainChecker.Run()
	assert.Nil(checkError)
}

func TestEnvoyInboundRouteDomainPodCheckerEmptyConfig(t *testing.T) {
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
	routeDomainChecker := HasInboundDynamicRouteConfigDomainPodCheck(configGetter, pod)
	checkError := routeDomainChecker.Run()
	assert.NotNil(checkError)
	assert.Equal("envoy config is empty", checkError.Error())
}

func TestEnvoyInboundRouteDomainPodCheckerNoDomains(t *testing.T) {
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
	routeDomainChecker := HasInboundDynamicRouteConfigDomainPodCheck(configGetter, pod)
	checkError := routeDomainChecker.Run()
	assert.NotNil(checkError)
	assert.Equal("no dynamic route config domains", checkError.Error())
}

func TestEnvoyInboundRouteDomainPodCheckerDomainNotFound(t *testing.T) {
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
	routeDomainChecker := HasInboundDynamicRouteConfigDomainPodCheck(configGetter, pod)
	checkError := routeDomainChecker.Run()
	assert.NotNil(checkError)
	assert.Equal("dynamic route config domain not found", checkError.Error())
}

func TestEnvoyOutboundRouteDomainHostChecker(t *testing.T) {
	assert := tassert.New(t)
	configGetter := mockConfigGetter{
		getter: createConfigGetterFunc("../../tests/sample-envoy-config-dump-bookbuyer.json"),
	}
	routeDomainChecker := HasOutboundDynamicRouteConfigDomainHostCheck(configGetter, bookstoreDestinationHost)
	checkError := routeDomainChecker.Run()
	assert.Nil(checkError)
}

func TestEnvoyOutboundRouteDomainHostCheckerEmptyConfig(t *testing.T) {
	assert := tassert.New(t)
	configGetter := mockConfigGetter{
		getter: func() (*Config, error) {
			return nil, nil
		},
	}
	routeDomainChecker := HasOutboundDynamicRouteConfigDomainHostCheck(configGetter, bookstoreDestinationHost)
	checkError := routeDomainChecker.Run()
	assert.NotNil(checkError)
	assert.Equal("envoy config is empty", checkError.Error())
}

func TestEnvoyOutboundRouteDomainHostCheckerNoDomains(t *testing.T) {
	assert := tassert.New(t)
	configGetter := mockConfigGetter{
		getter: createConfigGetterFunc("../../tests/sample-envoy-config-dump-bookbuyer-no-rds-dynamic-route-virtual-host-domains.json"),
	}
	routeDomainChecker := HasOutboundDynamicRouteConfigDomainHostCheck(configGetter, bookstoreDestinationHost)
	checkError := routeDomainChecker.Run()
	assert.NotNil(checkError)
	assert.Equal("no dynamic route config domains", checkError.Error())
}

func TestEnvoyOutboundRouteDomainHostCheckerDomainNotFound(t *testing.T) {
	assert := tassert.New(t)
	configGetter := mockConfigGetter{
		getter: createConfigGetterFunc("../../tests/sample-envoy-config-dump-bookbuyer-not-found-rds-dynamic-route-virtual-host-domain.json"),
	}
	routeDomainChecker := HasOutboundDynamicRouteConfigDomainHostCheck(configGetter, bookstoreDestinationHost)
	checkError := routeDomainChecker.Run()
	assert.NotNil(checkError)
	assert.Equal("dynamic route config domain not found", checkError.Error())
}
