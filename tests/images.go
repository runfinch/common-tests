// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// Images tests functionality of `images` command that lists container images.
func Images(o *option.Option) {
	const sha256RegexTruncated = `^[a-f0-9]{12}$`
	ginkgo.Describe("list container images", ginkgo.Ordered, func() {
		testImageName := "fn-test-images-cmd:latest"
		ginkgo.BeforeAll(func() {
			pullImage(o, defaultImage)
			buildImage(o, testImageName)
		})

		ginkgo.AfterAll(func() {
			removeImage(o, testImageName)
			removeImage(o, defaultImage)
		})

		ginkgo.It("should list all the images in a tabular format", func() {
			images := command.StdOutAsLines(o, "images")
			header, images := images[0], images[1:]
			gomega.Expect(images).ShouldNot(gomega.BeEmpty())
			gomega.Expect(header).Should(gomega.MatchRegexp(
				"REPOSITORY[\t ]+TAG[\t ]+IMAGE ID[\t ]+CREATED[\t ]+PLATFORM[\t ]+SIZE[\t ]+BLOB SIZE"))
			gomega.Expect(images).Should(gomega.HaveEach((gomega.MatchRegexp(`^(.+[\t ]+){6}.+$`))))
			// TODO: add more strict validation using output matcher.
		})
		ginkgo.It("should list all the images with image names in a tabular format ", func() {
			images := command.StdOutAsLines(o, "images", "--names")
			header, images := images[0], images[1:]
			gomega.Expect(images).ShouldNot(gomega.BeEmpty())
			gomega.Expect(header).Should(gomega.MatchRegexp("NAME[\t ]+IMAGE ID[\t ]+CREATED[\t ]+PLATFORM[\t ]+SIZE[\t ]+BLOB SIZE"))
			gomega.Expect(images).Should(gomega.HaveEach((gomega.MatchRegexp(`^(.+[\t ]+){5}.+$`))))
			// TODO: add more strict validation using output matcher.
		})
		ginkgo.It("should list all the images", func() {
			images := command.StdOutAsLines(o, "images", "--format", "{{.Repository}}:{{.Tag}}")
			gomega.Expect(images).ShouldNot(gomega.BeEmpty())
			gomega.Expect(images).Should(gomega.ContainElements(testImageName))
			gomega.Expect(images).Should(gomega.ContainElements(defaultImage))
		})
		ginkgo.It("should list truncated IMAGE IDs", func() {
			images := command.StdOutAsLines(o, "images", "--quiet")
			gomega.Expect(images).ShouldNot(gomega.BeEmpty())
			gomega.Expect(images).Should(gomega.HaveEach(gomega.MatchRegexp(sha256RegexTruncated)))
		})
		ginkgo.It("should list full IMAGE IDs", func() {
			images := command.StdOutAsLines(o, "images", "--quiet", "--no-trunc")
			gomega.Expect(images).ShouldNot(gomega.BeEmpty())
			gomega.Expect(images).Should(gomega.HaveEach(gomega.MatchRegexp(sha256RegexFull)))
		})
		ginkgo.It("should list IMAGE digests", func() {
			images := command.StdOutAsLines(o, "images", "--digests", "--format", "{{.Digest}}")
			gomega.Expect(images).ShouldNot(gomega.BeEmpty())
			gomega.Expect(images).Should(gomega.HaveEach(gomega.MatchRegexp(sha256RegexFull)))
		})
		// TODO: need to implement --filter functional test once we upgrade to nerdctl 0.23.
	})
}
