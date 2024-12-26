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

// NetworkLs tests the "network ls" command that list networks.
func NetworkLs(o *option.Option) {
	ginkgo.Describe("list networks", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.It("should list all the networks", func() {
			output := command.StdoutStr(o, "network", "ls")
			gomega.Expect(output).Should(gomega.ContainSubstring(bridgeNetwork))
		})

		ginkgo.It("should only display network name with --format flag", func() {
			lines := command.StdoutAsLines(o, "network", "ls", "--format", "{{.Name}}")
			gomega.Expect(lines).Should(gomega.ContainElement(bridgeNetwork))
		})

		for _, quiet := range []string{"-q", "--quiet"} {
			ginkgo.It(fmt.Sprintf("should only display network id with %s flag", quiet), func() {
				output := command.StdoutStr(o, "network", "ls", quiet)
				gomega.Expect(output).ShouldNot(gomega.BeEmpty())
				gomega.Expect(output).ShouldNot(gomega.ContainSubstring(bridgeNetwork))
			})
		}
	})
}
