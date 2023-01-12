// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"
)

// VolumeInspect tests "volume inspect" command that displays detailed information on one or more volumes.
func VolumeInspect(o *option.Option) {
	ginkgo.Describe("display detailed volume on a volume", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.It("should display the detailed information of volume", func() {
			command.Run(o, "volume", "create", testVolumeName)
			name := command.StdoutStr(o, "volume", "inspect", testVolumeName, "--format", "{{.Name}}")
			gomega.Expect(name).Should(gomega.Equal(testVolumeName))
			mp := command.StdoutStr(o, "volume", "inspect", testVolumeName, "--format", "{{.Mountpoint}}")
			gomega.Expect(mp).ShouldNot(gomega.BeEmpty())
		})

		ginkgo.It("should display detailed information of multiple volumes", func() {
			const testVol2 = "testVol2"
			command.Run(o, "volume", "create", testVolumeName)
			command.Run(o, "volume", "create", testVol2)
			lines := command.StdoutAsLines(o, "volume", "inspect", testVolumeName, "testVol2", "--format", "{{.Name}}")
			gomega.Expect(lines).Should(gomega.ContainElements(testVolumeName, testVol2))
		})

		ginkgo.It("should have error if inspect a nonexistent volume", func() {
			command.RunWithoutSuccessfulExit(o, "volume", "inspect", "ne-volume")
		})
	})
}
