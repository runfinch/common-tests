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

// ComposeLogs tests functionality of `compose logs` command.
func ComposeLogs(o *option.Option) {
	services := []string{"svc1_compose_logs", "svc2_compose_logs"}
	containerNames := []string{"container1_compose_logs", "container2_compose_logs"}
	imageNames := []string{defaultImage, defaultImage}

	ginkgo.Describe("Compose logs command", func() {
		var buildContext string
		var composeFilePath string
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
			buildContext, composeFilePath = createComposeYmlForLogsCmd(services, imageNames, containerNames)
			ginkgo.DeferCleanup(os.RemoveAll, buildContext)
			command.Run(o, "compose", "up", "-d", "--file", composeFilePath)
			containerShouldExist(o, containerNames...)
		})

		ginkgo.AfterEach(func() {
			command.Run(o, "compose", "down", "--file", composeFilePath)
			command.RemoveAll(o)
		})
		ginkgo.It("should show the logs with prefixes", func() {
			output := command.StdOutAsLines(o, "compose", "logs", "--file", composeFilePath)
			// Log format: container_name |log_msg
			// example: container1_composelogs |hello from service 1
			gomega.Expect(output).Should(gomega.ContainElements(
				gomega.HavePrefix(containerNames[0]),
				gomega.HavePrefix(containerNames[1])))
		})
		ginkgo.It("should show the logs without prefixes", func() {
			output := command.StdOutAsLines(o, "compose", "logs", "--no-log-prefix", "--file", composeFilePath)
			// Log format: log_msg
			// example: hello from service 1
			gomega.Expect(output).ShouldNot(gomega.ContainElements(
				gomega.HavePrefix(containerNames[0]),
				gomega.HavePrefix(containerNames[1])))
		})
		ginkgo.It("should show the logs with no color", func() {
			output := command.StdoutStr(o, "compose", "logs", "--no-color", "--file", composeFilePath)
			// The asci color code has prefix \x1b[3 e.g. Black: \u001b[30m, Red: \u001b[31m
			gomega.Expect(output).ShouldNot(gomega.ContainSubstring("\x1b[3"))
		})
		ginkgo.It("should show the last line of the logs", func() {
			output := command.StdOutAsLines(o, "compose", "logs", services[0], "--tail", "1", "--file", composeFilePath)
			gomega.Expect(output).Should(gomega.HaveLen(1))
		})

		for _, tFlag := range []string{"-t", "--timestamps"} {
			tFlag := tFlag
			ginkgo.It(fmt.Sprintf("should show the logs with timestamp with no prefixes and no color [flag %s]", tFlag), func() {
				// Log format: YYYY-MM-DDThh:mm:ss.000000000Z LOG MSG
				timestampMatcher := gomega.MatchRegexp("^[0-9]{1,4}-[0-9]{1,2}-[0-9]{1,2}T[0-9]{1,2}:[0-9]{1,2}:[0-9]{1,2}.*")
				output := command.StdOutAsLines(o, "compose", "logs", tFlag, "--no-log-prefix", "--no-color", "--file", composeFilePath)
				gomega.Expect(output).Should(gomega.HaveEach(timestampMatcher))
			})
		}
	})
}

func createComposeYmlForLogsCmd(serviceNames []string, imageNames []string, containerNames []string) (string, string) {
	gomega.Expect(serviceNames).Should(gomega.HaveLen(2))
	gomega.Expect(imageNames).Should(gomega.HaveLen(2))
	gomega.Expect(containerNames).Should(gomega.HaveLen(2))

	composeYmlContent := fmt.Sprintf(
		`
services:
  %[1]s:
    image: "%[3]s"
    container_name: "%[5]s"
    command: sh -c 'echo "hello from service 1"; echo "again hello"; sleep infinity'
  %[2]s:
    image: "%[4]s"
    container_name: "%[6]s"
    command: sh -c 'echo "hello from service 2"; echo "again hello"; sleep infinity'
`, serviceNames[0], serviceNames[1], imageNames[0], imageNames[1], containerNames[0], containerNames[1])
	return ffs.CreateComposeYmlContext(composeYmlContent)
}
