// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"fmt"
	"os"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/ffs"
	"github.com/runfinch/common-tests/fnet"
	"github.com/runfinch/common-tests/option"
)

// Push tests pushing an image to a registry.
func Push(o *option.Option) {
	ginkgo.Describe("Push a container image to registry", func() {
		var buildContext string
		var port int

		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
			buildContext = ffs.CreateBuildContext(fmt.Sprintf(`FROM %s
		CMD ["echo", "bar"]
			`, localImages[defaultImage]))
			ginkgo.DeferCleanup(os.RemoveAll, buildContext)
			port = fnet.GetFreePort()
			command.Run(o, "run", "-dp", fmt.Sprintf("%d:5000", port), "--name", "registry", registryImage)
		})

		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.Context("Test push command without any flag", func() {
			ginkgo.It("should push an image with a valid tag to registry", func() {
				tag := fmt.Sprintf(`localhost:%d/test-push:tag`, port)
				command.Run(o, "build", "-t", tag, buildContext)
				command.Run(o, "push", tag)
				command.Run(o, "pull", tag)
			})

			ginkgo.It("should return an error when pushing a nonexistent tag", func() {
				nonexistentTag := fmt.Sprintf(`localhost:%d/nonexistent:tag`, port)
				stderr := command.RunWithoutSuccessfulExit(o, "push", nonexistentTag).Err.Contents()
				gomega.Expect(stderr).To(gomega.ContainSubstring("not found"))
			})
		})
	})
}
