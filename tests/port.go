// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"fmt"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/fnet"
	"github.com/runfinch/common-tests/option"
)

// Port tests listing port mappings or a specific mapping for a container.
func Port(o *option.Option) {
	ginkgo.Describe("list port mapping", func() {
		const containerPort = 4567
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.It("should output port mappings for a container", func() {
			hostPort := fnet.GetFreePort()
			command.Run(o, "run", "-p", fmt.Sprintf("%d:%d", hostPort, containerPort), "--name", testContainerName, localImages[defaultImage])

			output := command.StdoutStr(o, "port", testContainerName)
			gomega.Expect(output).Should(gomega.Equal(fmt.Sprintf("%d/tcp -> 0.0.0.0:%d", containerPort, hostPort)))
		})

		ginkgo.It("should output the host port according to container port", func() {
			hostPort := fnet.GetFreePort()
			command.Run(o, "run", "-p", fmt.Sprintf("%d:%d", hostPort, containerPort), "--name", testContainerName, localImages[defaultImage])

			output := command.StdoutStr(o, "port", testContainerName, fmt.Sprintf("%d/tcp", containerPort))
			gomega.Expect(output).Should(gomega.Equal(fmt.Sprintf("0.0.0.0:%d", hostPort)))
		})

		ginkgo.It("should have error if specifying wrong protocol", func() {
			hostPort := fnet.GetFreePort()
			command.Run(o,
				"run",
				"-p",
				fmt.Sprintf("%d:%d/udp", hostPort, containerPort),
				"--name",
				testContainerName,
				localImages[defaultImage],
			)

			command.RunWithoutSuccessfulExit(o, "port", testContainerName, fmt.Sprintf("%d/tcp", containerPort))
		})

		ginkgo.It("should still output the host port according to container port when no protocol is specified", func() {
			hostPort := fnet.GetFreePort()
			command.Run(o, "run", "-p", fmt.Sprintf("%d:%d", hostPort, containerPort), "--name", testContainerName, localImages[defaultImage])

			output := command.StdoutStr(o, "port", testContainerName, fmt.Sprint(containerPort))
			gomega.Expect(output).Should(gomega.Equal(fmt.Sprintf("0.0.0.0:%d", hostPort)))
		})

		ginkgo.It("should have error if trying to print container port which is not published to any host port", func() {
			hostPort := fnet.GetFreePort()
			command.Run(o, "run", "-p", fmt.Sprintf("%d:%d", hostPort, containerPort), "--name", testContainerName, localImages[defaultImage])

			command.RunWithoutSuccessfulExit(o, "port", testContainerName, "111/tcp")
		})
	})
}
