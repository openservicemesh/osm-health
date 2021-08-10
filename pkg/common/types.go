package common

// Runnable is a type of generic function that can be executed; it returns pass/fail on a given check.
type Runnable interface {
	// Run executes a check.
	Run() error

	// Info returns human-readable information on what check is being executed.
	Info() string

	// Suggestion returns a human-readable suggestion on how to fix the issue.
	Suggestion() string

	// FixIt attempts to fix the issue at hand.
	FixIt() error
}

// MeshName is the type for the name of a mesh.
type MeshName string

func (mn MeshName) String() string {
	return string(mn)
}

// MeshNamespace is the type for the namespace of a mesh.
type MeshNamespace string

func (mns MeshNamespace) String() string {
	return string(mns)
}

// Outcome is the context around a check (common.Runnable), which failed.
type Outcome struct {
	// RunnableInfo will hold context on the executed runnable - common.Runnable{}.Info()
	RunnableInfo string

	// Error is the error which common.Runnable{}.Run() returned.
	Error error
}
