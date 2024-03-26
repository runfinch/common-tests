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
	"github.com/runfinch/common-tests/option"
)

// ComposePull tests functionality of `compose pull` command.
func ComposePull(o *option.Option) {
	services := []string{"svc1_compose_pull", "svc2_compose_pull"}
	ginkgo.Describe("Compose pull command", func() {
		var composeContext string
		var composeFilePath string
		var imageNames []string

		ginkgo.BeforeEach(func() {
			imageNames = []string{localImages[defaultImage], localImages[olderAlpineImage]}
			command.RemoveAll(o)
			composeContext, composeFilePath = createComposeYmlForPullCmd(services, imageNames)
			ginkgo.DeferCleanup(os.RemoveAll, composeContext)
		})

		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.It("should pull images for all services", func() {
			command.Run(o, "compose", "pull", "--file", composeFilePath)
			imageList := command.GetAllImageNames(o)
			gomega.Expect(imageList).Should(gomega.ContainElements(imageNames))
		})

		ginkgo.It("should pull the image for the first service only", func() {
			command.Run(o, "compose", "pull", services[0], "--file", composeFilePath)
			imageList := command.GetAllImageNames(o)
			gomega.Expect(imageList).Should(gomega.ContainElement(imageNames[0]))
			gomega.Expect(imageList).ShouldNot(gomega.ContainElement(imageNames[1]))
		})

		for _, qFlag := range []string{"-q", "--quiet"} {
			ginkgo.It(fmt.Sprintf("should pull the images without printing progress information with %s flag", qFlag), func() {
				qFlag := qFlag
				command.Run(o, "compose", "pull", qFlag, "--file", composeFilePath)
				imageList := command.GetAllImageNames(o)
				gomega.Expect(imageList).Should(gomega.ContainElements(imageNames))
			})
		}
	})
}

func createComposeYmlForPullCmd(serviceNames []string, imageNames []string) (string, string) {
	gomega.Expect(serviceNames).Should(gomega.HaveLen(2))
	gomega.Expect(imageNames).Should(gomega.HaveLen(2))

	composeYmlContent := fmt.Sprintf(
		`
services:
  %[1]s:
    image: "%[3]s"
  %[2]s:
    image: "%[4]s"
`, serviceNames[0], serviceNames[1], imageNames[0], imageNames[1])
	return ffs.CreateComposeYmlContext(composeYmlContent)
}
