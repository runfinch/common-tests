// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"fmt"
	"strings"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"
)

// Logs tests fetching logs of a container.
func Logs(o *option.Option) {
	ginkgo.Describe("fetch logs of a container", func() {
		const foo = "foo"
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.When("the container is not running and has one line of logs", func() {
			ginkgo.BeforeEach(func() {
				// Currently, only containers created with `run -d` are supported.
				// https://github.com/containerd/nerdctl#whale-nerdctl-logs
				command.Run(o, "run", "-d", "--name", testContainerName, localImages[defaultImage], "echo", foo)
			})

			ginkgo.It("should fetch the logs of a container", func() {
				output := command.StdoutStr(o, "logs", testContainerName)
				gomega.Expect(output).Should(gomega.Equal(foo))
			})

			for _, timestamps := range []string{"-t", "--timestamps"} {
				ginkgo.It(fmt.Sprintf("should include timestamp with %s flag", timestamps), func() {
					output := command.StdoutStr(o, "logs", timestamps, testContainerName)
					// `logs --timestamps` command will add an RFC3339Nano timestamp,
					// for example 2014-09-16T06:17:46.000000000Z, to each log entry.
					// "2006-01-02" is a golang common layout which specifies the format to be yyyy-MM-dd.
					gomega.Expect(output).Should(gomega.ContainSubstring(time.Now().UTC().Format("2006-01-02")))
					gomega.Expect(output).Should(gomega.ContainSubstring(foo))
				})
			}

			ginkgo.It("should show log message depending on a relative time with --since flag", func() {
				time.Sleep(2 * time.Second)
				output := command.StdoutStr(o, "logs", "--since", "1s", testContainerName)
				gomega.Expect(output).Should(gomega.BeEmpty())
				output = command.StdoutStr(o, "logs", "--since", "5s", testContainerName)
				gomega.Expect(output).Should(gomega.Equal(foo))
			})

			ginkgo.It("should show log message depending on a relative time with --until flag", func() {
				time.Sleep(2 * time.Second)
				output := command.StdoutStr(o, "logs", "--until", "1s", testContainerName)
				gomega.Expect(output).Should(gomega.Equal(foo))
				output = command.StdoutStr(o, "logs", "--until", "5s", testContainerName)
				gomega.Expect(output).Should(gomega.BeEmpty())
			})
		})

		ginkgo.When("the container is not running and has multiple lines of logs", func() {
			const bar = "bar"
			ginkgo.BeforeEach(func() {
				command.Run(o, "run", "-d", "--name", testContainerName, localImages[defaultImage],
					"sh", "-c", fmt.Sprintf("echo %s; echo %s", foo, bar))
			})

			for _, tail := range []string{"-n", "--tail"} {
				ginkgo.It(fmt.Sprintf("should show number of lines from end of the logs with %s flag", tail), func() {
					expectedOutput := fmt.Sprintf("%s\n%s", foo, bar)
					output := command.StdoutStr(o, "logs", tail, "1", testContainerName)
					gomega.Expect(output).Should(gomega.Equal(bar))
					output = command.StdoutStr(o, "logs", tail, "all", testContainerName)
					gomega.Expect(output).Should(gomega.Equal(expectedOutput))
				})
			}
		})

		ginkgo.When("the container is running", func() {
			ginkgo.BeforeEach(func() {
				command.Run(o, "run", "-d", "--name", testContainerName, localImages[defaultImage], "sleep", "infinity")
			})

			for _, follow := range []string{"-f", "--follow"} {
				ginkgo.It(fmt.Sprintf("should follow log output with %s flag", follow), func() {
					const newLog = "hello"
					session := command.RunWithoutWait(o, "logs", follow, testContainerName)
					defer session.Kill()
					gomega.Expect(session.Out.Contents()).Should(gomega.BeEmpty())
					command.Run(o, "exec", testContainerName, "sh", "-c", fmt.Sprintf("echo %s >> /proc/1/fd/1", newLog))
					// allow propagation time
					gomega.Eventually(func(session *gexec.Session) string {
						return strings.TrimSpace(string(session.Out.Contents()))
					}).WithArguments(session).
						WithTimeout(30 * time.Second).
						WithPolling(1 * time.Second).
						Should(gomega.Equal(newLog))
				})
			}
		})

		ginkgo.It("should fail if container doesn't exist", func() {
			command.RunWithoutSuccessfulExit(o, "logs", nonexistentContainerName)
		})
	})
}
