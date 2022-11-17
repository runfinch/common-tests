// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"
)

// NetworkRm tests the "network rm" command that removes one or more networks.
func NetworkRm(o *option.Option) {
	ginkgo.Describe("remove one or more networks", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
			command.Run(o, "network", "create", testNetwork)
		})

		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.It("should remove a network", func() {
			gomega.Expect(command.StdOutAsLines(o, "network", "ls", "--format", "{{.Name}}")).Should(gomega.ContainElement(testNetwork))
			command.Run(o, "network", "rm", testNetwork)
			gomega.Expect(command.StdOutAsLines(o, "network", "ls", "--format", "{{.Name}}")).ShouldNot(gomega.ContainElement(testNetwork))
		})

		ginkgo.It("should remove multiple networks", func() {
			const testNetwork2 = "test-network2"
			command.Run(o, "network", "create", testNetwork2)
			command.Run(o, "network", "rm", testNetwork, testNetwork2)
			lines := command.StdOutAsLines(o, "network", "ls", "--format", "{{.Name}}")
			gomega.Expect(lines).ShouldNot(gomega.ContainElement(testNetwork))
			gomega.Expect(lines).ShouldNot(gomega.ContainElement(testNetwork2))
		})
	})
}
