package namespace

import (
	"context"

	corev1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func getAnnotations(client kubernetes.Interface, namespace string) (map[string]string, error) {
	ns, err := client.CoreV1().Namespaces().Get(context.TODO(), namespace, corev1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return ns.Annotations, nil
}
