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

// ComposePs tests functionality of `compose ps` command.
func ComposePs(o *option.Option) {
	services := []string{"svc1_compose_ps", "svc2_compose_ps"}
	containerNames := []string{"container1_compose_ps", "container2_compose_ps"}
	imageNames := []string{localImages["defaultImage"], localImages["defaultImage"]}

	ginkgo.Describe("Compose ps command", func() {
		var composeContext string
		var composeFilePath string
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
			composeContext, composeFilePath = createComposeYmlForPsCmd(services, imageNames, containerNames)
			ginkgo.DeferCleanup(os.RemoveAll, composeContext)
			command.Run(o, "compose", "up", "-d", "--file", composeFilePath)
			containerShouldExist(o, containerNames...)
		})

		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.It("should list services defined in compose file", func() {
			psOutput := command.StdoutAsLines(o, "compose", "ps", "--file", composeFilePath)
			gomega.Expect(psOutput).Should(gomega.ContainElements(
				gomega.ContainSubstring(services[0]),
				gomega.ContainSubstring(services[1])))
			gomega.Expect(psOutput).Should(gomega.ContainElements(
				gomega.ContainSubstring(containerNames[0]),
				gomega.ContainSubstring(containerNames[1])))
			gomega.Expect(psOutput).Should(gomega.ContainElement(gomega.ContainSubstring("sleep infinity")))
			gomega.Expect(psOutput).Should(gomega.ContainElement(gomega.ContainSubstring("8080->8080/tcp")))
			gomega.Expect(psOutput).Should(gomega.ContainElement(gomega.ContainSubstring("running")))
		})
	})
}

func createComposeYmlForPsCmd(serviceNames []string, imageNames []string, containerNames []string) (string, string) {
	gomega.Expect(serviceNames).Should(gomega.HaveLen(2))
	gomega.Expect(imageNames).Should(gomega.HaveLen(2))
	gomega.Expect(containerNames).Should(gomega.HaveLen(2))

	composeYmlContent := fmt.Sprintf(
		`
services:
  %[1]s:
    image: "%[3]s"
    container_name: "%[5]s"
    command: sleep infinity
    ports:
      - 8080:8080
  %[2]s:
    image: "%[4]s"
    container_name: "%[6]s"
    command: sleep infinity
`, serviceNames[0], serviceNames[1], imageNames[0], imageNames[1], containerNames[0], containerNames[1])
	return ffs.CreateComposeYmlContext(composeYmlContent)
}
