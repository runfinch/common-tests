// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/ffs"
	"github.com/runfinch/common-tests/option"
)

// Exec tests executing a command in a running container.
func Exec(o *option.Option) {
	ginkgo.Describe("execute command in a container", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})
		// TODO: specifying -t flag will have error in test -> panic: provided file is not a console
		ginkgo.When("then container is running", func() {
			ginkgo.BeforeEach(func() {
				command.Run(o, "run", "-d", "--name", testContainerName, defaultImage, "sleep", "infinity")
			})

			ginkgo.It("should execute a command in a running container", func() {
				strEchoed := "hello"
				output := command.StdoutStr(o, "exec", testContainerName, "echo", strEchoed)
				gomega.Expect(output).Should(gomega.Equal(strEchoed))
			})

			for _, interactive := range []string{"-i", "--interactive"} {
				interactive := interactive
				ginkgo.It(fmt.Sprintf("should output string by piping if %s flag keeps STDIN open", interactive), func() {
					want := []byte("hello")
					got := command.New(o, "exec", interactive, testContainerName, "cat").
						WithStdin(gbytes.BufferWithBytes(want)).Run().Out.Contents()
					gomega.Expect(got).Should(gomega.Equal(want))
				})
			}

			for _, detach := range []string{"-d", "--detach"} {
				detach := detach
				ginkgo.It(fmt.Sprintf("should execute command in detached mode with %s flag", detach), func() {
					command.Run(o, "exec", detach, testContainerName, "nc", "-l")
					processes := command.StdoutStr(o, "exec", testContainerName, "ps", "aux")
					gomega.Expect(processes).Should(gomega.ContainSubstring("nc -l"))
				})
			}

			for _, workDir := range []string{"-w", "--workdir"} {
				workDir := workDir
				ginkgo.It(fmt.Sprintf("should execute command under directory specified by %s flag", workDir), func() {
					dir := "/tmp"
					output := command.StdoutStr(o, "exec", workDir, dir, testContainerName, "pwd")
					gomega.Expect(output).Should(gomega.Equal(dir))
				})
			}

			for _, env := range []string{"-e", "--env"} {
				env := env
				ginkgo.It(fmt.Sprintf("should set the environment variable with %s flag", env), func() {
					const envPair = "ENV=1"
					lines := command.StdoutAsLines(o, "exec", env, envPair, testContainerName, "env")
					gomega.Expect(lines).Should(gomega.ContainElement(envPair))
				})
			}

			ginkgo.It("should set environment variables from file with --env-file flag", func() {
				const envPair = "ENV=1"
				envPath := ffs.CreateTempFile("env", envPair)
				ginkgo.DeferCleanup(os.RemoveAll, filepath.Dir(envPath))

				envOutput := command.StdoutAsLines(o, "exec", "--env-file", envPath, testContainerName, "env")
				gomega.Expect(envOutput).Should(gomega.ContainElement(envPair))
			})

			ginkgo.It("should execute command in privileged mode with --privileged flag", func() {
				command.RunWithoutSuccessfulExit(o, "exec", testContainerName, "ip", "link", "add", "dummy1", "type", "dummy")
				command.Run(o, "exec", "--privileged", testContainerName, "ip", "link", "add", "dummy1", "type", "dummy")
				output := command.StdoutStr(o, "exec", "--privileged", testContainerName, "ip", "link")
				gomega.Expect(output).Should(gomega.ContainSubstring("dummy1"))
			})

			for _, user := range []string{"-u", "--user"} {
				user := user
				ginkgo.It(fmt.Sprintf("should output user id according to user name specified by %s flag", user), func() {
					testCases := map[string]string{
						"1000":       "uid=1000 gid=0(root)",
						"1000:users": "uid=1000 gid=100(users)",
					}

					for name, want := range testCases {
						output := command.StdoutStr(o, "exec", user, name, testContainerName, "id")
						gomega.Expect(output).Should(gomega.ContainSubstring(want))
					}
				})
			}
		})

		ginkgo.It("should not execute a command when the container is not running", func() {
			command.Run(o, "run", "--name", testContainerName, defaultImage)
			command.RunWithoutSuccessfulExit(o, "exec", testContainerName)
		})
	})
}
