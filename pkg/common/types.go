package common

import "fmt"

// Pod is an identifier of a pod within a cluster that is being tested.
type Pod struct {
	Cluster   string
	Namespace string
	Name      string
}

func (p Pod) String() string {
	return fmt.Sprintf("%s/%s/%s", p.Cluster, p.Namespace, p.Name)
}

// Result is the output of a test.
type Result struct {
	SMIPolicy

	Successful bool
}

// SMIPolicy contains context around how SMI affects a result of a test.
type SMIPolicy struct {
	HasPolicy                  bool
	ValidPolicy                bool
	SourceToDestinationAllowed bool
	// TODO: include actual SMI policy
}
