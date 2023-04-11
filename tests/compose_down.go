// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/ffs"
	"github.com/runfinch/common-tests/option"
	"os"
)

// ComposeDown tests functionality of `compose down` command.
func ComposeDown(o *option.Option) {
	services := []string{"svc1_compose_down", "svc2_compose_down"}
	containerNames := []string{"container1_compose_down", "container2_compose_down"}

	ginkgo.Describe("Compose down command", func() {
		var composeContext string
		var composeFilePath string
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
			composeContext, composeFilePath = createComposeYmlForDownCmd(services, containerNames)
			ginkgo.DeferCleanup(os.RemoveAll, composeContext)
			command.Run(o, "compose", "up", "-d", "--file", composeFilePath)
			containerShouldExist(o, containerNames...)
		})

		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.It("should stop services defined in compose file without deleting volumes", func() {
			command.Run(o, "compose", "down", "--file", composeFilePath)
			containerShouldNotExist(o, containerNames...)
			volumeShouldExist(o, "compose_data_volume")
		})
	})
}

func createComposeYmlForDownCmd(serviceNames []string, containerNames []string) (string, string) {
	gomega.Expect(serviceNames).Should(gomega.HaveLen(2))
	gomega.Expect(containerNames).Should(gomega.HaveLen(2))

	composeYmlContent := fmt.Sprintf(
		`
services:
  %[1]s:
    image: "%[3]s"
    container_name: "%[4]s"
    command: sleep infinity
    volumes:
      - compose_data_volume:/usr/local/data
  %[2]s:
    image: "%[3]s"
    container_name: "%[5]s"
    command: sleep infinity
    volumes:
      - compose_data_volume:/usr/local/data
volumes:
  compose_data_volume:
`, serviceNames[0], serviceNames[1], defaultImage, containerNames[0], containerNames[1])
	return ffs.CreateComposeYmlContext(composeYmlContent)
}
