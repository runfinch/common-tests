// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"
)

// NetworkInspect tests the "network inspect" command that displays detailed information on one or more networks.
func NetworkInspect(o *option.Option) {
	ginkgo.Describe("display detailed information on network", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.It("should display detailed information about one network", func() {
			name := command.StdoutStr(o, "network", "inspect", bridgeNetwork, "--format", "{{.Name}}")
			gomega.Expect(name).Should(gomega.Equal(bridgeNetwork))
		})

		ginkgo.It("should display detailed information on multiple networks", func() {
			command.Run(o, "network", "create", testNetwork)
			lines := command.StdOutAsLines(o, "network", "inspect", bridgeNetwork, testNetwork, "--format", "{{.Name}}")
			gomega.Expect(lines).Should(gomega.ConsistOf(bridgeNetwork, testNetwork))
		})
	})
}
