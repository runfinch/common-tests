// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"
)

// Info tests "info" command that displays system-wide information, synonyms to "system info" command.
func Info(o *option.Option) {
	ginkgo.Describe("display system-wide information", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.It("should display system-wide information", func() {
			gomega.Expect(command.StdoutStr(o, "system", "info", "--format", "{{.OSType}}")).ShouldNot(gomega.BeEmpty())
			gomega.Expect(command.StdoutStr(o, "system", "info", "--format", "{{.Architecture}}")).ShouldNot(gomega.BeEmpty())
		})
	})
}
