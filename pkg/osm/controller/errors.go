package controller

import "errors"

var (
	// ErrorNoControllerPodsExistInNamespace denotes when no osm-controller pods exist in the specified namespace.
	ErrorNoControllerPodsExistInNamespace = errors.New("no osm-controller pods exist in the specified namespace")
)
