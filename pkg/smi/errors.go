package smi

import "errors"

var (
	// ErrInvalidRuleKind is an error returned when the TrafficTarget has an unexpected rule kind.
	ErrInvalidRuleKind = errors.New("unsupported rule kind in TrafficTarget")
)
