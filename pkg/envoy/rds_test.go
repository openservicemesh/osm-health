package envoy

import (
	"testing"

	tassert "github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
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
	client := fake.NewSimpleClientset(pod)
	routeDomainChecker := NewOutboundRouteDomainPodCheck(client, configGetter, pod)
	outcome := routeDomainChecker.Run()
	assert.Nil(outcome.GetError())
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
	client := fake.NewSimpleClientset(pod)
	routeDomainChecker := NewOutboundRouteDomainPodCheck(client, configGetter, pod)
	outcome := routeDomainChecker.Run()
	assert.NotNil(outcome.GetError())
	assert.Equal(ErrEnvoyConfigEmpty.Error(), outcome.GetError().Error())
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
	client := fake.NewSimpleClientset(pod)
	routeDomainChecker := NewOutboundRouteDomainPodCheck(client, configGetter, pod)
	outcome := routeDomainChecker.Run()
	assert.NotNil(outcome.GetError())
	assert.Equal(ErrNoDynamicRouteConfigDomains.Error(), outcome.GetError().Error())
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
			Labels: map[string]string{
				"mykey": "myval",
			},
		},
	}
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bookstore",
			Namespace: "bookstore",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"mykey": "myval",
			},
		},
	}
	client := fake.NewSimpleClientset(pod, svc)
	routeDomainChecker := NewOutboundRouteDomainPodCheck(client, configGetter, pod)
	outcome := routeDomainChecker.Run()
	assert.NotNil(outcome.GetError())
	assert.Equal(ErrDynamicRouteConfigDomainNotFound.Error(), outcome.GetError().Error())
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
	client := fake.NewSimpleClientset(pod)
	routeDomainChecker := NewInboundRouteDomainPodCheck(client, configGetter, pod)
	outcome := routeDomainChecker.Run()
	assert.Nil(outcome.GetError())
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
	client := fake.NewSimpleClientset(pod)
	routeDomainChecker := NewInboundRouteDomainPodCheck(client, configGetter, pod)
	outcome := routeDomainChecker.Run()
	assert.NotNil(outcome.GetError())
	assert.Equal(ErrEnvoyConfigEmpty.Error(), outcome.GetError().Error())
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
	client := fake.NewSimpleClientset(pod)
	routeDomainChecker := NewInboundRouteDomainPodCheck(client, configGetter, pod)
	outcome := routeDomainChecker.Run()
	assert.NotNil(outcome.GetError())
	assert.Equal(ErrNoDynamicRouteConfigDomains.Error(), outcome.GetError().Error())
}

func TestEnvoyInboundRouteDomainPodCheckerDomainNotFound(t *testing.T) {
	assert := tassert.New(t)
	configGetter := mockConfigGetter{
		getter: createConfigGetterFunc("../../tests/sample-envoy-config-dump-bookstore-not-found-rds-dynamic-route-virtual-host-domain.json"),
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bookstore",
			Namespace: "bookstore",
			Labels: map[string]string{
				"mykey": "myval",
			},
		},
	}
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bookstore",
			Namespace: "bookstore",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"mykey": "myval",
			},
		},
	}
	client := fake.NewSimpleClientset(pod, svc)
	routeDomainChecker := NewInboundRouteDomainPodCheck(client, configGetter, pod)
	outcome := routeDomainChecker.Run()
	assert.NotNil(outcome.GetError())
	assert.Equal(ErrDynamicRouteConfigDomainNotFound.Error(), outcome.GetError().Error())
}

func TestEnvoyOutboundRouteDomainHostChecker(t *testing.T) {
	assert := tassert.New(t)
	configGetter := mockConfigGetter{
		getter: createConfigGetterFunc("../../tests/sample-envoy-config-dump-bookbuyer.json"),
	}
	routeDomainChecker := NewOutboundRouteDomainHostCheck(configGetter, bookstoreDestinationHost)
	outcome := routeDomainChecker.Run()
	assert.Nil(outcome.GetError())
}

func TestEnvoyOutboundRouteDomainHostCheckerEmptyConfig(t *testing.T) {
	assert := tassert.New(t)
	configGetter := mockConfigGetter{
		getter: func() (*Config, error) {
			return nil, nil
		},
	}
	routeDomainChecker := NewOutboundRouteDomainHostCheck(configGetter, bookstoreDestinationHost)
	outcome := routeDomainChecker.Run()
	assert.NotNil(outcome.GetError())
	assert.Equal(ErrEnvoyConfigEmpty.Error(), outcome.GetError().Error())
}

func TestEnvoyOutboundRouteDomainHostCheckerNoDomains(t *testing.T) {
	assert := tassert.New(t)
	configGetter := mockConfigGetter{
		getter: createConfigGetterFunc("../../tests/sample-envoy-config-dump-bookbuyer-no-rds-dynamic-route-virtual-host-domains.json"),
	}
	routeDomainChecker := NewOutboundRouteDomainHostCheck(configGetter, bookstoreDestinationHost)
	outcome := routeDomainChecker.Run()
	assert.NotNil(outcome.GetError())
	assert.Equal(ErrNoDynamicRouteConfigDomains.Error(), outcome.GetError().Error())
}

func TestEnvoyOutboundRouteDomainHostCheckerDomainNotFound(t *testing.T) {
	assert := tassert.New(t)
	configGetter := mockConfigGetter{
		getter: createConfigGetterFunc("../../tests/sample-envoy-config-dump-bookbuyer-not-found-rds-dynamic-route-virtual-host-domain.json"),
	}
	routeDomainChecker := NewOutboundRouteDomainHostCheck(configGetter, bookstoreDestinationHost)
	outcome := routeDomainChecker.Run()
	assert.NotNil(outcome.GetError())
	assert.Equal(ErrDynamicRouteConfigDomainNotFound.Error(), outcome.GetError().Error())
}
