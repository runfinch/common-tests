// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"
)

// Events tests "events" command that gets real time events from server, synonyms to "system events" command.
func Events(o *option.Option) {
	ginkgo.Describe("get real time events from the server", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.It("should get real time events from command", func() {
			session := command.RunWithoutWait(o, "system", "events")
			defer session.Kill()
			gomega.Expect(session.Out.Contents()).Should(gomega.BeEmpty())
			command.Run(o, "pull", defaultImage)
			gomega.Expect(session.Out.Contents()).Should(gomega.ContainSubstring(defaultImage))
		})
	})
}
