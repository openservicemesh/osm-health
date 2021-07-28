package kubernetes

import "fmt"

// Pod is an identifier of a pod within a cluster that is being tested.
type Pod struct {
	Namespace Namespace
	Name      string
}

func (p Pod) String() string {
	return fmt.Sprintf("%s/%s", p.Namespace, p.Name)
}

// Namespace is a unique identifier of a namespace on a Kubernetes cluster.
type Namespace string

func (ns Namespace) String() string {
	return string(ns)
}
