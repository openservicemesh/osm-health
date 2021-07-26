package namespace

import (
	"context"

	kubernetes2 "github.com/openservicemesh/osm-health/pkg/kubernetes"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func getAnnotations(client kubernetes.Interface, namespace kubernetes2.Namespace) (map[string]string, error) {
	ns, err := client.CoreV1().Namespaces().Get(context.TODO(), namespace.String(), v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return ns.Annotations, nil
}

func getLabels(client kubernetes.Interface, namespace kubernetes2.Namespace) (map[string]string, error) {
	ns, err := client.CoreV1().Namespaces().Get(context.TODO(), namespace.String(), v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return ns.Labels, nil
}
