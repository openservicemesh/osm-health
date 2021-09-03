package common

import (
	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
)

// Runnable is a type of generic function that can be executed; it returns pass/fail on a given check.
type Runnable interface {
	// Run executes a check and returns an outcome (outcomes.Outcome).
	Run() outcomes.Outcome

	// Description returns human-readable information on what check is being executed.
	Description() string

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

// Printable is the printable context around a check (common.Runnable).
type Printable struct {
	// CheckDescription holds the description of a check, such as describing what the check does (common.Runnable)
	CheckDescription string

	// Type holds the type of the check outcome, such as success, fail, info or unknown
	Type string

	// Diagnostics holds detailed diagnostics that were dynamically-generated during the check
	Diagnostics string

	// Error is the error which common.Runnable{}.Run() may return
	Error error
}
