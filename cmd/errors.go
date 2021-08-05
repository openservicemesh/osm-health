package main

import "errors"

var (
	// ErrNoSourcePodOrNoDestinationPod is an error when the user does not supply the SOURCE_POD or the DESTINATION_POD.
	ErrNoSourcePodOrNoDestinationPod = errors.New("provide both SOURCE_POD and DESTINATION_POD")

	// ErrInvalidSourcePod is an error when the supplied source pod is invalid.
	ErrInvalidSourcePod = errors.New("invalid SOURCE_POD")

	// ErrInvalidDestinationPod is an error when the supplied destination pod is invalid.
	ErrInvalidDestinationPod = errors.New("invalid DESTINATION_POD")

	// ErrNoSourcePodOrNoDestinationURL is an error when the user does not supply the SOURCE_POD or the DESTINATION_URL.
	ErrNoSourcePodOrNoDestinationURL = errors.New("provide both SOURCE_POD and DESTINATION_URL")

	// ErrInvalidDestinationURL is an error when the supplied destination URL is invalid.
	ErrInvalidDestinationURL = errors.New("invalid DESTINATION_URL")

	// ErrNoDestinationPod is an error when the user does not supply the DESTINATION_POD.
	ErrNoDestinationPod = errors.New("provide DESTINATION_POD")
)
