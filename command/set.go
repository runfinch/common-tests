// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package command

var (
	localRegistryImageID     string
	localRegistryContainerID string
	localRegistryImageName   string
)

// SetLocalRegistryContainerID sets the ID for the local registry. Usually you don't need to invoke this function yourself.
// For more details, see tests.SetupLocalRegistry.
func SetLocalRegistryContainerID(id string) {
	localRegistryContainerID = id
}

// SetLocalRegistryImageID sets the ID for local registry image. Usually you don't need to invoke this function yourself.
// For more details, see tests.SetupLocalRegistry.
func SetLocalRegistryImageID(id string) {
	localRegistryImageID = id
}

// SetLocalRegistryImageName sets the local registry image name. Usually you don't need to invoke this function yourself.
// For more details, see tests.SetupLocalRegistry.
func SetLocalRegistryImageName(name string) {
	localRegistryImageName = name
}
