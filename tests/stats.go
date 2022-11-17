// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"
)

// Stats tests displaying container resource usage statistics.
func Stats(o *option.Option) {
	ginkgo.Describe("display a container", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})
		// TODO: add tests for -a flag
		// REF issue: https://github.com/containerd/nerdctl/issues/1415
		// TODO: add test for streaming data
		ginkgo.When("the container is running", func() {
			ginkgo.BeforeEach(func() {
				command.Run(o, "run", "-d", "--name", testContainerName, defaultImage, "sleep", "infinity")
			})

			ginkgo.It("should disable streaming usage stats and print result with --no-stream flag", func() {
				output := command.StdoutStr(o, "stats", "--no-stream", testContainerName, "--format", "{{.Name}}")
				gomega.Expect(output).Should(gomega.Equal(testContainerName))
			})

			ginkgo.It("should not truncate output with --no-trunc flag", func() {
				noTruncated := command.StdoutStr(o, "stats", "--no-stream", "--no-trunc", testContainerName, "--format", "{{.ID}}")
				truncated := command.StdoutStr(o, "stats", "--no-stream", testContainerName, "--format", "{{.ID}}")
				gomega.Expect(len(noTruncated) > len(truncated)).Should(gomega.BeTrue())
				gomega.Expect(noTruncated).Should(gomega.ContainSubstring(truncated))
			})
		})

		ginkgo.It("should not print usage stats if container doesn't exist", func() {
			command.RunWithoutSuccessfulExit(o, "stats", nonexistentContainerName)
		})
	})
}
