package namespace

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	k8s "github.com/openservicemesh/osm-health/pkg/kubernetes"
)

func getAnnotations(client kubernetes.Interface, namespace k8s.Namespace) (map[string]string, error) {
	ns, err := client.CoreV1().Namespaces().Get(context.TODO(), namespace.String(), v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return ns.Annotations, nil
}
