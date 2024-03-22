// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"github.com/onsi/gomega"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"

	"github.com/onsi/ginkgo/v2"
)

// Restart tests "restart" command that will restart one or more running containers.
func Restart(o *option.Option) {
	// TODO: add tests for -t/--time flag
	// REF issue - https://github.com/containerd/nerdctl/issues/1485
	ginkgo.Describe("restart command", ginkgo.Ordered, func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
			// Functionality wise, we only need `sleep infinity` to keep the container running,
			// but with PID=1, `sleep infinity` will only exit when receiving SIGKILL,
			// which means that we'll have to wait for the default timeout (10 seconds for now) to restart the container,
			// so we use `nc -l` instead to save time.
			// TODO: Remove the above comment after we add a test case for -t/--time flag with `sleep infinity` because it's more obvious.
			command.Run(o, "run", "-d", "--name", testContainerName, localImages[defaultImage], "nc", "-l")
		})

		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.It("should restart a running container", func() {
			pid := getContainerPID(o, testContainerName)
			command.Run(o, "restart", testContainerName)
			newPid := getContainerPID(o, testContainerName)

			gomega.Expect(pid).NotTo(gomega.Equal(newPid))
		})

		ginkgo.It("should restart multiple running containers", func() {
			const ctrName = "ctr-name"
			command.Run(o, "run", "-d", "--name", ctrName, localImages[defaultImage], "nc", "-l")
			pid := getContainerPID(o, testContainerName)
			pid2 := getContainerPID(o, ctrName)
			command.Run(o, "restart", testContainerName, ctrName)
			newPid := getContainerPID(o, testContainerName)
			newPid2 := getContainerPID(o, ctrName)
			gomega.Expect(pid).NotTo(gomega.Equal(newPid))
			gomega.Expect(pid2).NotTo(gomega.Equal(newPid2))
		})

		ginkgo.It("should have error when restarting a nonexistent container", func() {
			command.RunWithoutSuccessfulExit(o, "restart", nonexistentContainerName)
		})
	})
}

func getContainerPID(o *option.Option, containerName string) string {
	return command.StdoutStr(o, "inspect", containerName, "--format", "{{.State.Pid}}")
}
