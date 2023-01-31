// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"github.com/runfinch/common-tests/option"
)

// GetAllContainerIDs returns all container IDs.
func GetAllContainerIDs(o *option.Option) []string {
	return StdoutAsLines(o, "ps", "--all", "--quiet", "--no-trunc")
}

// GetAllImageNames returns all image names.
func GetAllImageNames(o *option.Option) []string {
	return StdoutAsLines(o, "images", "--all", "--format", "{{.Repository}}:{{.Tag}}")
}

// GetAllVolumeNames returns all volume names.
func GetAllVolumeNames(o *option.Option) []string {
	return StdoutAsLines(o, "volume", "ls", "--quiet")
}

// GetAllNetworkNames returns all network names.
func GetAllNetworkNames(o *option.Option) []string {
	return StdoutAsLines(o, "network", "ls", "--format", "{{.Name}}")
}

// GetAllImageIDs returns all image IDs.
func GetAllImageIDs(o *option.Option) []string {
	return StdoutAsLines(o, "images", "--all", "--quiet")
}
