// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// Inspect tests displaying the detailed information of image or container.
func Inspect(o *option.Option) {
	ginkgo.Describe("inspect a container", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.It("should display the detailed information of a container", func() {
			command.Run(o, "run", "--name", testContainerName, localImages["defaultImage"])
			image := command.StdoutStr(o, "inspect", "--format", "{{.Image}}", testContainerName)
			gomega.Expect(image).To(gomega.Equal(localImages["defaultImage"]))
			containerName := command.StdoutStr(o, "inspect", "--format", "{{.Name}}", testContainerName)
			gomega.Expect(containerName).To(gomega.Equal(testContainerName))
			gomega.Expect(command.StdoutStr(o, "inspect", "--format", "{{.State.Status}}", testContainerName)).To(gomega.Equal("exited"))
			gomega.Expect(command.StdoutStr(o, "inspect", "--format", "{{.State.Error}}", testContainerName)).To(gomega.Equal(""))
		})

		ginkgo.It("should display multiple container image with --format flag", func() {
			const oldContainerName = "ctr-old"
			command.Run(o, "run", "--name", testContainerName, localImages["defaultImage"])
			command.Run(o, "run", "--name", oldContainerName, localImages["olderAlpineImage"])
			images := command.StdoutAsLines(o, "inspect", "--format", "{{.Image}}", testContainerName, oldContainerName)
			gomega.Expect(images).Should(gomega.ConsistOf(localImages["defaultImage"], localImages["olderAlpineImage"]))
		})

		ginkgo.It("should have an error if inspect a non-existing container", func() {
			command.RunWithoutSuccessfulExit(o, "inspect", nonexistentContainerName)
		})

		ginkgo.It("should show the information of a container with --type=container flag", func() {
			command.Run(o, "run", "--name", testContainerName, localImages["defaultImage"])
			image := command.StdoutStr(o, "inspect", "--type", "container", testContainerName, "--format", "{{.Image}}")
			gomega.Expect(image).Should(gomega.Equal(localImages["defaultImage"]))
			containerName := command.StdoutStr(o, "inspect", "--format", "{{.Name}}", testContainerName)
			gomega.Expect(containerName).Should(gomega.Equal(testContainerName))
		})

		ginkgo.It("should show the information of an image with --type=image flag", func() {
			pullImage(o, localImages["defaultImage"])
			image := command.StdoutStr(o, "inspect", "--type", "image", localImages["defaultImage"], "--format", "{{(index .RepoTags 0)}}")
			gomega.Expect(image).Should(gomega.Equal(localImages["defaultImage"]))
		})

		ginkgo.It("should have an error if specify the wrong object type", func() {
			command.Run(o, "run", "--name", testContainerName, localImages["defaultImage"])
			command.RunWithoutSuccessfulExit(o, "inspect", "--type", "image", testContainerName)
		})

		ginkgo.It("should have an error if inspect a non-existing image", func() {
			command.RunWithoutSuccessfulExit(o, "inspect", nonexistentImageName)
		})
	})
}
