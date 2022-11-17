// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"
)

// Tag tests tagging a container image.
func Tag(o *option.Option) {
	ginkgo.Describe("tag a container image", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.It("should tag an image when the image exists", func() {
			pullImage(o, defaultImage)

			command.Run(o, "tag", defaultImage, testImageName)
			defaultImageID := command.StdOut(o, "images", "--quiet", "--no-trunc", defaultImage)
			taggedImageID := command.StdOut(o, "images", "--quiet", "--no-trunc", testImageName)
			gomega.Expect(taggedImageID).ShouldNot(gomega.BeEmpty())
			gomega.Expect(taggedImageID).To(gomega.Equal(defaultImageID))
		})
		ginkgo.It("should not tag an image when the image doesn't exist", func() {
			command.RunWithoutSuccessfulExit(o, "tag", nonexistentImageName, testImageName)
			imageShouldNotExist(o, testImageName)
		})
	})
}
