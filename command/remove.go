// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"github.com/onsi/ginkgo/v2"

	"github.com/runfinch/common-tests/option"
)

// RemoveAll removes all containers and images in the testing environment specified by o.
func RemoveAll(o *option.Option) {
	RemoveContainers(o)
	RemoveImages(o)
	RemoveVolumes(o)
	RemoveNetworks(o)
}

// RemoveContainers removes all containers in the testing environment specified by o.
func RemoveContainers(o *option.Option) {
	allIDs := GetAllContainerIDs(o)
	var ids []string
	for _, id := range allIDs {
		if id != localRegistryContainerID {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		ginkgo.GinkgoWriter.Println("No containers to be removed")
		return
	}

	args := append([]string{"rm", "--force"}, ids...)
	Run(o, args...)
}

// RemoveVolumes removes all unused local volumes in the testing environment specified by o.
func RemoveVolumes(o *option.Option) {
	volumes := GetAllVolumeNames(o)
	if len(volumes) == 0 {
		ginkgo.GinkgoWriter.Println("No volumes to be removed")
		return
	}
	Run(o, "volume", "prune", "--force", "--all")
}

// RemoveImages removes all container images in the testing environment specified by o.
func RemoveImages(o *option.Option) {
	allIDs := GetAllImageIDs(o)
	var ids []string
	for _, id := range allIDs {
		if id != localRegistryImageID {
			ids = append(ids, id)
		}
	}
	if removedAllImages(ids) {
		return
	}
	args := append([]string{"rmi", "--force"}, ids...)
	Run(o, args...)

	allNames := GetAllImageNames(o)
	var names []string
	for _, name := range allNames {
		if name != localRegistryImageName {
			names = append(names, name)
		}
	}
	if removedAllImages(names) {
		return
	}
	args = append([]string{"rmi", "--force"}, names...)
	Run(o, args...)
}

// RemoveNetworks removes all networks in the testing environment specified by o.
// TODO: use "network prune" after upgrading nerdctl to v0.23.
func RemoveNetworks(o *option.Option) {
	defaultNetworks := []string{"bridge", "host", "none"}
	networks := GetAllNetworkNames(o)
	var customNetworks []string
	for _, n := range networks {
		if !contains(defaultNetworks, n) {
			customNetworks = append(customNetworks, n)
		}
	}

	if len(customNetworks) == 0 {
		ginkgo.GinkgoWriter.Println("No networks to be removed")
		return
	}

	args := append([]string{"network", "rm"}, customNetworks...)
	Run(o, args...)
}

func contains(strs []string, target string) bool {
	for _, str := range strs {
		if str == target {
			return true
		}
	}
	return false
}

func removedAllImages(images []string) bool {
	if len(images) == 0 {
		ginkgo.GinkgoWriter.Println("No images to be removed")
		return true
	}
	return false
}
