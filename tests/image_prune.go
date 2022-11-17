// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"fmt"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega/gbytes"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"
)

// ImagePrune tests "image prune" command that removes unused images.
func ImagePrune(o *option.Option) {
	// Currently, nerdctl image prune requires --all to be specified.
	// REF - https://github.com/containerd/nerdctl#whale-nerdctl-image-prune
	// TODO: Add a test case to only prune dangling images after `--all` is not required for `image prune`.
	ginkgo.Describe("Remove unused images", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
			pullImage(o, defaultImage)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.It("should remove all unused images with inputting y in prompt confirmation", func() {
			imageShouldExist(o, defaultImage)
			command.New(o, "image", "prune", "-a").WithStdin(gbytes.BufferWithBytes([]byte("y"))).Run()
			imageShouldNotExist(o, defaultImage)
		})

		ginkgo.It("should not remove any unused image with inputting n in prompt confirmation", func() {
			imageShouldExist(o, defaultImage)
			command.New(o, "image", "prune", "-a").WithStdin(gbytes.BufferWithBytes([]byte("n"))).Run()
			imageShouldExist(o, defaultImage)
		})

		for _, force := range []string{"-f", "--force"} {
			force := force
			ginkgo.It(fmt.Sprintf("with %s flag, should remove unused images without prompting a confirmation", force), func() {
				imageShouldExist(o, defaultImage)
				command.Run(o, "image", "prune", "-a", "-f")
				imageShouldNotExist(o, defaultImage)
			})
		}

		ginkgo.It("should not remove an image if it's used by a dead container", func() {
			command.Run(o, "run", defaultImage)
			imageShouldExist(o, defaultImage)
			command.Run(o, "image", "prune", "-a", "-f")
			imageShouldExist(o, defaultImage)
		})
	})
}
