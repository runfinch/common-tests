// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"
)

// ImageInspect tests "image inspect" command that displays detailed information on one or more images.
func ImageInspect(o *option.Option) {
	ginkgo.Describe("display detailed information on one or more images", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
			pullImage(o, defaultImage)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.It("should display detailed information on an image", func() {
			gomega.Expect(command.StdoutStr(o, "image", "inspect", defaultImage)).ShouldNot(gomega.BeEmpty())
		})

		ginkgo.It("should display image RepoTags with --format flag", func() {
			image := command.StdoutStr(o, "image", "inspect", defaultImage, "--format", "{{(index .RepoTags 0)}}")
			gomega.Expect(image).Should(gomega.Equal(defaultImage))
		})

		ginkgo.It("should display multiple image RepoTags with --format flag", func() {
			pullImage(o, olderAlpineImage)
			lines := command.StdOutAsLines(o, "image", "inspect", defaultImage, olderAlpineImage, "--format", "{{(index .RepoTags 0)}}")
			gomega.Expect(lines).Should(gomega.ConsistOf(defaultImage, olderAlpineImage))
		})

		ginkgo.It("should not display information if image doesn't exist", func() {
			command.RunWithoutSuccessfulExit(o, "image", "inspect", nonexistentImageName)
		})
	})
}
