package common

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

// Runnable is a type of generic function that can be executed and it will return pass/fail on a given check.
type Runnable interface {
	Run() error
	Info() string
}

// MeshName is the type for the name of a mesh.
type MeshName string

func (mn MeshName) String() string {
	return string(mn)
}

// Outcome is the context around a check (common.Runnable), which failed.
type Outcome struct {
	// RunnableInfo will hold context on the executed runnable - common.Runnable{}.Info()
	RunnableInfo string

	// Error is the error which common.Runnable{}.Run() returned.
	Error error
}
