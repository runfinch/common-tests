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

// VolumeLs tests "volume ls" command that lists volumes.
func VolumeLs(o *option.Option) {
	ginkgo.Describe("list volumes", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})
		// TODO: add test for --filter after upgrading to nerdctl v0.23
		ginkgo.It("should display all the volumes", func() {
			const testVol2 = "testVol2"
			command.Run(o, "volume", "create", testVolumeName)
			command.Run(o, "volume", "create", testVol2)
			lines := command.StdOutAsLines(o, "volume", "ls", "--format", "{{.Name}}")
			gomega.Expect(lines).Should(gomega.ContainElements(testVolumeName, testVol2))
		})

		for _, quiet := range []string{"--quiet", "-q"} {
			quiet := quiet
			ginkgo.It(fmt.Sprintf("should only display volume names with %s flag", quiet), func() {
				command.Run(o, "volume", "create", testVolumeName)
				gomega.Expect(command.StdOutAsLines(o, "volume", "ls", quiet)).Should(gomega.ContainElement(testVolumeName))
			})
		}
	})
}
