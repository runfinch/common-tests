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
		ginkgo.It("should prune the builder cache successfully", func() {
			command.Run(o, "build", "--output=type=docker", buildContext)
			// There is no interface to validate the current builder cache size.
			// To validate in Buildkit, run `buildctl du -v`.
			args := []string{"builder", "prune"}

			// TODO(macedonv): remove after nerdctlv2 is supported across all platforms.
			if o.IsNerdctlV2() {
				// Do not prompt for user response during automated testing.
				args = append(args, "--force")
			}

			command.Run(o, args...)
		})
	})
}
