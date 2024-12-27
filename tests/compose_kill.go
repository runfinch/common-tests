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

// ComposeKill tests functionality of `compose` command.
func ComposeKill(o *option.Option) {
	services := []string{"svc1_compose_kill", "svc2_compose_kill"}
	containerNames := []string{"container1_compose_kill", "container2_compose_kill"}

	ginkgo.Describe("Compose kill command", func() {
		var composeContext string
		var composeFilePath string
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
			composeContext, composeFilePath = createComposeYmlForKillCmd(services, containerNames)
			ginkgo.DeferCleanup(os.RemoveAll, composeContext)
			command.Run(o, "compose", "up", "-d", "--file", composeFilePath)

			containerShouldExist(o, containerNames...)
		})

		ginkgo.AfterEach(func() {
			command.Run(o, "compose", "down", "--file", composeFilePath)
			command.RemoveAll(o)
		})
		ginkgo.It("should kill all service", func() {
			command.Run(o, "compose", "kill", "--file", composeFilePath)
			containerShouldNotBeRunning(o, containerNames...)
		})

		for _, signal := range []string{"-s", "--signal"} {
			for _, term := range []string{"SIGTERM", "TERM"} {
				ginkgo.It(fmt.Sprintf("should send %s to containers when using %s", term, signal), func() {
					command.Run(o, "compose", "kill", signal, term, "--file", composeFilePath)
					containerShouldNotBeRunning(o, containerNames...)
				})
			}
		}
	})
}

func createComposeYmlForKillCmd(serviceNames []string, containerNames []string) (string, string) {
	gomega.Expect(serviceNames).Should(gomega.HaveLen(2))
	gomega.Expect(containerNames).Should(gomega.HaveLen(2))

	// Service commands implement SIGTERM handler to test compose kill
	// can send non-default signals.
	//
	// With PID=1, `sleep infinity` would only exit when receiving SIGKILL.
	// https://stackoverflow.com/questions/45148381/why-cant-i-ctrl-c-a-sleep-infinity-in-docker-when-it-runs-as-pid-1
	composeYmlContent := fmt.Sprintf(
		`
services:
  %[1]s:
    image: "%[3]s"
    container_name: "%[4]s"
    command: |
      sh -c "
        trap 'echo shutting down; exit 0' SIGTERM
        sleep infinity &
        wait
      "
  %[2]s:
    image: "%[3]s"
    container_name: "%[5]s"
    command: |
      sh -c "
        trap 'echo shutting down; exit 0' SIGTERM
        sleep infinity &
        wait
      "
`, serviceNames[0], serviceNames[1], localImages[defaultImage], containerNames[0], containerNames[1])
	return ffs.CreateComposeYmlContext(composeYmlContent)
}
