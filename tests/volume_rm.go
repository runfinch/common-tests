// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"fmt"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"
)

// VolumeRm tests "volume rm" command that removes one or more volumes.
func VolumeRm(o *option.Option) {
	ginkgo.Describe("remove a volume", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.When("volumes are not used by any container", func() {
			ginkgo.BeforeEach(func() {
				command.Run(o, "volume", "create", testVolumeName)
			})
			ginkgo.It("should remove a volume", func() {
				volumeShouldExist(o, testVolumeName)
				command.Run(o, "volume", "rm", testVolumeName)
				volumeShouldNotExist(o, testVolumeName)
			})

			ginkgo.It("should remove multiple volumes", func() {
				const testVol2 = "testVol2"
				command.Run(o, "volume", "create", "testVol2")
				gomega.Expect(command.StdOutAsLines(o, "volume", "ls", "--quiet")).Should(gomega.ContainElements(testVolumeName, testVol2))
				command.Run(o, "volume", "rm", testVolumeName, testVol2)
				volumeShouldNotExist(o, testVolumeName)
			})
		})

		ginkgo.When("a volume is used by a container", func() {
			ginkgo.BeforeEach(func() {
				command.Run(o, "volume", "create", testVolumeName)
				command.Run(o, "run", "-v", fmt.Sprintf("%s:/tmp", testVolumeName), defaultImage)
			})

			// It's expected that `volume rm` can't remove the volume that is referenced to a container despite the container status.
			// REF - https://github.com/containerd/nerdctl/blob/657cf4be42f9e99ee0fd53103d4ded62d7137aa3/cmd/nerdctl/volume_rm.go#L36
			// TODO: add test for --force/-f after they are implemented.
			// REF - https://github.com/containerd/nerdctl/blob/657cf4be42f9e99ee0fd53103d4ded62d7137aa3/cmd/nerdctl/volume_rm.go#L43
			ginkgo.It("should not remove the volume that is referenced to a container", func() {
				command.RunWithoutSuccessfulExit(o, "volume", "rm", testVolumeName)
				volumeShouldExist(o, testVolumeName)
			})
		})
	})
}
