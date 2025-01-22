// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"fmt"
	"os"

	"github.com/onsi/ginkgo/v2"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/ffs"
	"github.com/runfinch/common-tests/option"
)

// BuilderPrune tests the "builder prune" command that prunes the builder cache.
func BuilderPrune(o *option.Option) {
	ginkgo.Describe("prune the builder cache", func() {
		var buildContext string
		ginkgo.BeforeEach(func() {
			buildContext = ffs.CreateBuildContext(fmt.Sprintf(`FROM %s
			CMD ["echo", "finch-test-dummy-output"]
			`, localImages[defaultImage]))
			ginkgo.DeferCleanup(os.RemoveAll, buildContext)
			command.RemoveAll(o)
		})

		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.Describe("with nerdctl v1.x command", func() {
			if o.IsNerdctlV2() {
				ginkgo.Skip("nerdctl runtime dependency is v2")
			}
			ginkgo.It("should prune the builder cache successfully", func() {
				// There is no interface to validate the current builder cache size.
				// To validate in Buildkit, run `buildctl du -v`.
				command.Run(o, "build", "--output=type=docker", buildContext)
				command.Run(o, "builder", "prune")
			})
		})

		ginkgo.Describe("with nerdctl v2.x command", func() {
			if o.IsNerdctlV1() {
				ginkgo.Skip("nerdctl runtime dependency is v1")
			}
			ginkgo.DescribeTable("should prune the builder cache successfully",
				func(args ...string) {
					// There is no interface to validate the current builder cache size.
					// To validate in Buildkit, run `buildctl du -v`.
					args = append([]string{"builder", "prune"}, args...)
					command.Run(o, "build", "--output=type=docker", buildContext)
					command.Run(o, args...)
				},
				ginkgo.Entry("with '-f -a' flags", "-f", "-a"),
				ginkgo.Entry("with '--force -a' flags", "--force", "-a"),
				ginkgo.Entry("with '-f --all' flags", "-f", "--all"),
				ginkgo.Entry("with '--force --all' flags", "--force", "--all"),
			)
		})
	})
}
