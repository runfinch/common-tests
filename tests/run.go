// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/ffs"
	"github.com/runfinch/common-tests/fnet"
	"github.com/runfinch/common-tests/option"
)

// RunOption is the custom option to run tests in run.go.
type RunOption struct {
	// BaseOpt instructs how to run the test subject.
	BaseOpt *option.Option
	// CGMode is the cgroup mode that the host uses.
	CGMode CGMode
	// DefaultHostGatewayIP is the IP that the test subject will resolve special IP `host-gateway` to.
	DefaultHostGatewayIP string
}

// Run tests running a container image.
// TODO: test cases for bind propagation options.
func Run(o *RunOption) {
	ginkgo.Describe("Run a container image", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o.BaseOpt)
		})

		ginkgo.AfterEach(func() {
			command.RemoveAll(o.BaseOpt)
		})

		ginkgo.When("running a container that echos dummy output", func() {
			ginkgo.BeforeEach(func() {
				dockerfile := fmt.Sprintf(`FROM %s
			CMD ["echo", "finch-test-dummy-output"]
			`, defaultImage)
				buildContext := ffs.CreateBuildContext(dockerfile)
				ginkgo.DeferCleanup(os.RemoveAll, buildContext)
				command.Run(o.BaseOpt, "build", "-q", "-t", testImageName, buildContext)
			})

			ginkgo.It("should echo dummy output", func() {
				output := command.StdoutStr(o.BaseOpt, "run", testImageName)
				gomega.Expect(output).Should(gomega.Equal("finch-test-dummy-output"))
			})

			ginkgo.It("should not echo dummy output if running with -d flag", func() {
				output := command.Stdout(o.BaseOpt, "run", "-d", testImageName)
				gomega.Expect(output).ShouldNot(gomega.ContainSubstring("finch-test-dummy-output"))
			})
		})

		ginkgo.It("with --rm flag, container should be removed when it exits", func() {
			command.Run(o.BaseOpt, "run", "--rm", "--name", testContainerName, defaultImage)
			containerShouldNotExist(o.BaseOpt, testContainerName)
		})

		ginkgo.When("running a container with metadata related flags", func() {
			for _, label := range []string{"-l", "--label"} {
				label := label
				ginkgo.It(fmt.Sprintf("should set meta data on a container with %s flag", label), func() {
					command.Run(o.BaseOpt, "run", "--name", testContainerName, label, "testKey=testValue", defaultImage)
					gomega.Expect(command.StdoutStr(o.BaseOpt, "inspect", testContainerName,
						"--format", "{{.Config.Labels.testKey}}")).To(gomega.Equal("testValue"))
				})
			}

			ginkgo.It("should write the container ID to file with --cidfile flag", func() {
				dir := ffs.CreateTempDir("finch-test-cid")
				ginkgo.DeferCleanup(os.RemoveAll, dir)
				path := filepath.Join(dir, "test.cid")
				containerID := command.StdoutStr(o.BaseOpt, "run", "-d", "--cidfile", path, defaultImage)
				output, err := os.ReadFile(filepath.Clean(path))
				gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
				gomega.Expect(strings.TrimSpace(string(output))).Should(gomega.Equal(containerID))
			})

			ginkgo.It("should read labels from file with --label-file flag", func() {
				path := ffs.CreateTempFile("label-file", "key=value")
				ginkgo.DeferCleanup(os.RemoveAll, filepath.Dir(path))
				command.Run(o.BaseOpt, "run", "--name", testContainerName, "--label-file", path, defaultImage)
				gomega.Expect(command.StdoutStr(o.BaseOpt, "inspect", testContainerName,
					"--format", "{{.Config.Labels.key}}")).To(gomega.Equal("value"))
			})
		})

		ginkgo.When("running a container with environment related flags", func() {
			ginkgo.It("with --entrypoint flag, ENTRYPOINT in dockerfile should not be executed", func() {
				dockerfile := fmt.Sprintf(`FROM %s
		ENTRYPOINT ["echo", "foo"]
		CMD ["echo", "bar"]
			`, defaultImage)
				buildContext := ffs.CreateBuildContext(dockerfile)
				defer func() {
					gomega.Expect(os.RemoveAll(buildContext)).To(gomega.Succeed())
				}()
				command.Run(o.BaseOpt, "build", "-q", "-t", testImageName, buildContext)

				envOutput := command.Stdout(o.BaseOpt, "run", "--rm", "--entrypoint", "time", testImageName, "echo", "blah")
				gomega.Expect(envOutput).NotTo(gomega.ContainSubstring("foo"))
				gomega.Expect(envOutput).NotTo(gomega.ContainSubstring("bar"))
				gomega.Expect(envOutput).To(gomega.ContainSubstring("blah"))
			})

			for _, workDir := range []string{"--workdir", "-w"} {
				workDir := workDir
				ginkgo.It(fmt.Sprintf("should set working directory inside the container specified by %s flag", workDir), func() {
					dir := "/tmp"
					gomega.Expect(command.StdoutStr(o.BaseOpt, "run", workDir, dir, defaultImage, "pwd")).Should(gomega.Equal(dir))
				})
			}

			for _, env := range []string{"-e", "--env"} {
				env := env
				ginkgo.It(fmt.Sprintf("with %s flag, environment variables should be set in the container", env), func() {
					envOutput := command.Stdout(o.BaseOpt, "run", "--rm",
						env, "FOO=BAR", env, "FOO1", env, "ENV1=1", env, "ENV1=2",
						defaultImage, "env")
					gomega.Expect(envOutput).To(gomega.ContainSubstring("FOO=BAR"))
					gomega.Expect(envOutput).ToNot(gomega.ContainSubstring("FOO1"))
					gomega.Expect(envOutput).To(gomega.ContainSubstring("ENV1=2"))
				})
			}

			ginkgo.It("with -e flag passing env variables without a value, only host set vars should be set in the container", func() {
				gomega.Expect(os.Setenv("AVAR1", "avalue")).To(gomega.Succeed())
				envOutput := command.Stdout(o.BaseOpt, "run", "--rm",
					"-e", "AVAR1", "-e", "AVAR2", defaultImage, "env")
				gomega.Expect(envOutput).To(gomega.ContainSubstring("AVAR1=avalue"))
				gomega.Expect(envOutput).ToNot(gomega.ContainSubstring("AVAR2"))
			})

			ginkgo.It("with --env-file flag, environment variables in container should pick up those in environment file", func() {
				const envPair = "ENVKEY=ENVVAL"
				envPath := ffs.CreateTempFile("env", envPair)
				ginkgo.DeferCleanup(os.RemoveAll, filepath.Dir(envPath))

				envOutput := command.Stdout(o.BaseOpt, "run", "--rm", "--env-file", envPath, defaultImage, "env")
				gomega.Expect(envOutput).To(gomega.ContainSubstring(envPair))
			})

			ginkgo.It("using an env var file, env vars without values should only be set in the container if they are set on the host", func() {
				const envPair = "ENVKEY=ENVVAL\nAVAR1\nAVAR2\n"
				envPath := ffs.CreateTempFile("env", envPair)
				ginkgo.DeferCleanup(os.RemoveAll, filepath.Dir(envPath))
				gomega.Expect(os.Setenv("AVAR1", "avalue")).To(gomega.Succeed())
				envOutput := command.Stdout(o.BaseOpt, "run", "--rm", "--env-file", envPath, defaultImage, "env")
				gomega.Expect(envOutput).To(gomega.ContainSubstring("ENVKEY=ENVVAL"))
				gomega.Expect(envOutput).To(gomega.ContainSubstring("AVAR1=avalue"))
				gomega.Expect(envOutput).ToNot(gomega.ContainSubstring("AVAR2"))
			})

			ginkgo.It("using a file with the --env-file flag, comments and whitespace should be ignored properly", func() {
				const envPair = "ENVKEY=ENVVAL   \n# this is a comment\n\n AVAR1\nAVAR1\nAVAR2\n  # comment 2\n"
				envPath := ffs.CreateTempFile("env", envPair)
				ginkgo.DeferCleanup(os.RemoveAll, filepath.Dir(envPath))
				gomega.Expect(os.Setenv("AVAR1", "avalue")).To(gomega.Succeed())
				envOutput := command.Stdout(o.BaseOpt, "run", "--rm", "--env-file", envPath, defaultImage, "env")
				gomega.Expect(envOutput).To(gomega.ContainSubstring("ENVKEY=ENVVAL"))
				gomega.Expect(envOutput).To(gomega.ContainSubstring("AVAR1=avalue"))
				gomega.Expect(envOutput).ToNot(gomega.ContainSubstring("AVAR2"))
			})
		})

		ginkgo.When("running an image with --pull flag", func() {
			ginkgo.It("should have an error if set --pull=never and the image doesn't exist", func() {
				command.RunWithoutSuccessfulExit(o.BaseOpt, "run", "--pull", "never", defaultImage)
				imageShouldNotExist(o.BaseOpt, defaultImage)
			})

			ginkgo.It("should be able to run the container if set --pull=never and the image exists", func() {
				const containerName = "test-container"
				pullImage(o.BaseOpt, defaultImage)
				command.Run(o.BaseOpt, "run", "--name", containerName, "--pull", "never", defaultImage)
				containerShouldExist(o.BaseOpt, containerName)
			})

			ginkgo.It("should be able to run the container if set --pull=missing and the image doesn't exist", func() {
				command.Run(o.BaseOpt, "run", "--name", testContainerName, "--pull", "missing", defaultImage)
				containerShouldExist(o.BaseOpt, testContainerName)
				imageShouldExist(o.BaseOpt, defaultImage)
			})

			ginkgo.It("should be able to run the container if set --pull=missing and the image exists", func() {
				pullImage(o.BaseOpt, defaultImage)
				command.Run(o.BaseOpt, "run", "--name", testContainerName, "--pull", "missing", defaultImage)
				containerShouldExist(o.BaseOpt, testContainerName)
			})

			ginkgo.It("should be able to run the container if set --pull=always and the image doesn't exist", func() {
				command.Run(o.BaseOpt, "run", "--name", testContainerName, "--pull", "always", defaultImage)
				containerShouldExist(o.BaseOpt, testContainerName)
				imageShouldExist(o.BaseOpt, defaultImage)
			})
			ginkgo.It("should be able to run the container if set --pull=always and the image exists", func() {
				pullImage(o.BaseOpt, defaultImage)
				command.Run(o.BaseOpt, "run", "--name", testContainerName, "--pull", "always", defaultImage)
			})
		})

		for _, interactive := range []string{"-i", "--interactive"} {
			interactive := interactive
			ginkgo.It(fmt.Sprintf("should output string if %s flag keeps STDIN open", interactive), func() {
				want := []byte("hello")
				got := command.New(o.BaseOpt, "run", interactive, defaultImage, "cat").
					WithStdin(gbytes.BufferWithBytes(want)).Run().Out.Contents()
				gomega.Expect(got).Should(gomega.Equal(want))
			})
		}

		ginkgo.It("should stop running container within specified time by --stop-timeout flag", func() {
			// With PID=1, `sleep infinity` does not exit due to receiving a SIGTERM, which is sent by the stop command.
			// Ref. https://superuser.com/a/1299463/730265
			command.Run(o.BaseOpt, "run", "-d", "--name", testContainerName, "--stop-timeout", "1", defaultImage, "sleep", "infinity")
			gomega.Expect(command.StdoutStr(o.BaseOpt, "exec", testContainerName, "echo", "foo")).To(gomega.Equal("foo"))
			startTime := time.Now()
			command.Run(o.BaseOpt, "stop", testContainerName)
			// assert the container to be stopped within 1.5 seconds
			gomega.Expect(time.Since(startTime)).To(gomega.BeNumerically("~", 0*time.Millisecond, 1500*time.Millisecond))
			command.RunWithoutSuccessfulExit(o.BaseOpt, "exec", testContainerName, "echo", "foo")
		})

		ginkgo.It("should immediately stop the container with --stop-signal=SIGKILL", func() {
			// With PID=1, `sleep infinity` will only exit when receiving SIGKILL, while the signal sent by stop command is SIGTERM.
			command.Run(o.BaseOpt, "run", "-d", "--name", testContainerName, "--stop-signal", "SIGKILL", defaultImage, "sleep", "infinity")
			containerShouldBeRunning(o.BaseOpt, testContainerName)
			startTime := time.Now()
			command.Run(o.BaseOpt, "stop", testContainerName)
			gomega.Expect(time.Since(startTime)).To(gomega.BeNumerically("~", 0*time.Millisecond, 500*time.Millisecond))
			status := command.StdoutStr(o.BaseOpt, "inspect", "--format", "{{.State.Status}}", testContainerName)
			gomega.Expect(status).Should(gomega.Equal("exited"))
			command.RunWithoutSuccessfulExit(o.BaseOpt, "exec", testContainerName, "echo", "foo")
		})

		ginkgo.It("should share PID namespace with host with --pid=host", func() {
			command.Run(o.BaseOpt, "run", "-d", "--name", testContainerName, "--pid=host", defaultImage, "sleep", "infinity")
			pid := command.StdoutStr(o.BaseOpt, "inspect", "--format", "{{.State.Pid}}", testContainerName)
			command.Run(o.BaseOpt, "exec", testContainerName, "sh", "-c", fmt.Sprintf("ps -o pid,comm | grep '%s sleep'", pid))
		})

		ginkgo.It("should share PID namespace with a container with --pid=container:<container>", func() {
			command.Run(o.BaseOpt, "run", "-d", "--name", testContainerName, defaultImage, "sleep", "infinity")
			// We are joining the pid namespace that was "created" by testContainerName,
			// so the pid=1 process will be the main process of testContainerName, which is `sleep`.
			command.Run(o.BaseOpt, "exec", testContainerName, "sh", "-c", "ps -o pid,comm | grep '1 sleep'")
		})

		ginkgo.When("running a container with network related flags", func() {
			// TODO: add tests for --ip, --mac-address flags
			for _, network := range []string{"--net", "--network"} {
				network := network
				ginkgo.It(fmt.Sprintf("should connect a container to a network with %s flag", network), func() {
					command.Run(o.BaseOpt, "run", "-d", network, "bridge", "--name", testContainerName,
						defaultImage, "sh", "-c", "echo hello | nc -l -p 80")
					ipAddr := command.StdoutStr(o.BaseOpt, "inspect", "--format",
						"{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}", testContainerName)
					output := command.StdoutStr(o.BaseOpt, "run", network, "bridge", defaultImage, "nc", fmt.Sprintf("%s:80", ipAddr))
					gomega.Expect(output).Should(gomega.Equal("hello"))
				})

				ginkgo.It(fmt.Sprintf("should use the same network with container specified by %s=container:<name>", network), func() {
					command.Run(o.BaseOpt, "run", "-d", network, "bridge", "--name", testContainerName,
						defaultImage, "sh", "-c", "echo hello | nc -l -p 80")
					ipAddr := command.StdoutStr(o.BaseOpt, "inspect", "--format",
						"{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}", testContainerName)
					output := command.StdoutStr(o.BaseOpt, "run", fmt.Sprintf("%s=container:%s", network, testContainerName),
						defaultImage, "nc", fmt.Sprintf("%s:80", ipAddr))
					gomega.Expect(output).Should(gomega.Equal("hello"))
				})

				ginkgo.It(fmt.Sprintf("should use the same network with container specified by %s=container:<id>", network), func() {
					id := command.StdoutStr(o.BaseOpt, "run", "-d", network, "bridge", "--name", testContainerName,
						defaultImage, "sh", "-c", "echo hello | nc -l -p 80")
					ipAddr := command.StdoutStr(o.BaseOpt, "inspect", "--format",
						"{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}", testContainerName)
					output := command.StdoutStr(o.BaseOpt, "run", fmt.Sprintf("%s=container:%s", network, id),
						defaultImage, "nc", fmt.Sprintf("%s:80", ipAddr))
					gomega.Expect(output).Should(gomega.Equal("hello"))
				})
			}

			ginkgo.It("should be able to set custom DNS servers with --dns flag", func() {
				const nameserver = "10.10.10.10"
				lines := command.StdoutAsLines(o.BaseOpt, "run", "--dns", nameserver, "--name", testContainerName,
					defaultImage, "cat", "/etc/resolv.conf")
				gomega.Expect(lines).Should(gomega.ContainElement(fmt.Sprintf("nameserver %s", nameserver)))
			})

			ginkgo.It("should be able to set custom DNS search domains with --dns-search flag", func() {
				lines := command.StdoutAsLines(o.BaseOpt, "run", "--dns-search", "test", "--name", testContainerName,
					defaultImage, "cat", "/etc/resolv.conf")
				gomega.Expect(lines).Should(gomega.ContainElement("search test"))
			})

			for _, dnsOption := range []string{"--dns-opt", "--dns-option"} {
				dnsOption := dnsOption
				ginkgo.It(fmt.Sprintf("should be able to set DNS option with %s flag", dnsOption), func() {
					lines := command.StdoutAsLines(o.BaseOpt, "run", dnsOption, "debug", "--name", testContainerName,
						defaultImage, "cat", "/etc/resolv.conf")
					gomega.Expect(lines).Should(gomega.ContainElement("options debug"))
				})
			}

			for _, hostname := range []string{"--hostname", "-h"} {
				hostname := hostname
				ginkgo.It(fmt.Sprintf("should be able to set container host name with %s flag", hostname), func() {
					name := command.StdoutStr(o.BaseOpt, "run", hostname, "foo", defaultImage, "hostname")
					gomega.Expect(name).Should(gomega.Equal("foo"))
				})
			}

			ginkgo.It("should add a custom host-to-IP mapping with --add-host flag", func() {
				mapping := command.StdoutStr(o.BaseOpt, "run", "--add-host", "test-host:6.6.6.6", defaultImage, "cat", "/etc/hosts")
				gomega.Expect(mapping).Should(gomega.ContainSubstring("6.6.6.6"))
				gomega.Expect(mapping).Should(gomega.ContainSubstring("test-host"))
			})

			ginkgo.It("should add a custom host-to-IP mapping with --add-host flag with special IP", func() {
				response := "This is the expected response for --add-host special IP test."
				http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					io.WriteString(w, response) //nolint:errcheck,gosec // Function call in server handler for testing only.
				})
				hostPort := fnet.GetFreePort()
				s := http.Server{Addr: fmt.Sprintf(":%d", hostPort), Handler: nil, ReadTimeout: 30 * time.Second}
				go s.ListenAndServe() //nolint:errcheck // Asynchronously starting server for testing only.
				ginkgo.DeferCleanup(s.Shutdown, context.Background())
				command.Run(o.BaseOpt, "run", "-d", "--name", testContainerName, "--add-host", "test-host:host-gateway",
					amazonLinux2Image, "sleep", "infinity")
				mapping := command.StdoutStr(o.BaseOpt, "exec", testContainerName, "cat", "/etc/hosts")
				gomega.Expect(mapping).Should(gomega.ContainSubstring(o.DefaultHostGatewayIP))
				gomega.Expect(mapping).Should(gomega.ContainSubstring("test-host"))
				gomega.Expect(command.StdoutStr(o.BaseOpt, "exec", testContainerName, "curl",
					fmt.Sprintf("test-host:%d", hostPort))).Should(gomega.Equal(response))
				command.Run(o.BaseOpt, "run", "-d", "--name", testContainerName2, "--add-host=test-host:host-gateway",
					amazonLinux2Image, "sleep", "infinity")
				mapping = command.StdoutStr(o.BaseOpt, "exec", testContainerName2, "cat", "/etc/hosts")
				gomega.Expect(mapping).Should(gomega.ContainSubstring(o.DefaultHostGatewayIP))
				gomega.Expect(mapping).Should(gomega.ContainSubstring("test-host"))
				gomega.Expect(command.StdoutStr(o.BaseOpt, "exec", testContainerName2, "curl",
					fmt.Sprintf("test-host:%d", hostPort))).Should(gomega.Equal(response))
			})

			for _, publish := range []string{"-p", "--publish"} {
				publish := publish
				ginkgo.It(fmt.Sprintf("port of the container should be published to the host port with %s flag", publish), func() {
					const (
						containerPort = 80
						strEchoed     = "hello"
					)
					hostPort := fnet.GetFreePort()
					// The nc version used by alpine image is https://busybox.net/downloads/BusyBox.html#nc,
					// which is different from that in linux: https://linux.die.net/man/1/nc.
					command.Run(o.BaseOpt, "run", "-d", publish, fmt.Sprintf("%d:%d", hostPort, containerPort), defaultImage,
						"sh", "-c", fmt.Sprintf("echo %s | nc -l -p %d", strEchoed, containerPort))

					fnet.DialAndRead("tcp", fmt.Sprintf("localhost:%d", hostPort), strEchoed, 20, 200*time.Millisecond)
				})
			}
		})

		ginkgo.When("running a container with volume flags", func() {
			const (
				destDir    = "/tmp"
				volumeType = "volume"
				bindType   = "bind"
				tmpfsType  = "tmpfs"
			)
			for _, volume := range []string{"-v", "--volume"} {
				volume := volume
				ginkgo.It(fmt.Sprintf("should mount a volume when running a container with %s", volume), func() {
					command.Run(o.BaseOpt, "run", "--name", testContainerName, volume,
						fmt.Sprintf("%s:%s", testVolumeName, destDir), defaultImage, "sh", "-c", "echo foo > /tmp/test.txt")
					srcDir := command.StdoutStr(o.BaseOpt, "volume", "inspect", testVolumeName, "--format", "{{.Mountpoint}}")
					expectedMount := []MountJSON{makeMount(volumeType, srcDir, destDir, "", true)}
					actualMount := getContainerMounts(o.BaseOpt, testContainerName)
					verifyMountsInfo(actualMount, expectedMount)
					output := command.StdoutStr(o.BaseOpt, "run", "-v", fmt.Sprintf("%s:/tmp", testVolumeName), defaultImage,
						"cat", "/tmp/test.txt")
					gomega.Expect(output).Should(gomega.Equal("foo"))
				})

				ginkgo.It(fmt.Sprintf("should be able to set the volume options with %s testVol:%s:ro", volume, destDir), func() {
					command.Run(o.BaseOpt, "run", "-d", "--name", testContainerName, volume,
						fmt.Sprintf("%s:%s:ro", testVolumeName, destDir), defaultImage, "sleep", "infinity")
					srcDir := command.StdoutStr(o.BaseOpt, "volume", "inspect", "--format", "{{.Mountpoint}}", testVolumeName)
					expectedMount := []MountJSON{makeMount(volumeType, srcDir, destDir, "ro", false)}
					actualMount := getContainerMounts(o.BaseOpt, testContainerName)
					verifyMountsInfo(actualMount, expectedMount)
					// verify the volume is readonly
					command.RunWithoutSuccessfulExit(o.BaseOpt, "exec", testContainerName, "sh", "-c",
						fmt.Sprintf("echo foo > %s/test.txt", destDir))
				})
			}

			ginkgo.It("should create a tmpfs mount in a container", func() {
				const tmpfsContainerName = "tmpfs-ctr"
				command.Run(o.BaseOpt, "run", "-d", "--tmpfs", fmt.Sprintf("%s:size=64m,exec", destDir),
					"--name", tmpfsContainerName, defaultImage, "sleep", "infinity")
				expectedMount := []MountJSON{makeMount(tmpfsType, tmpfsType, destDir, "size=64m,exec", true)}
				actualMount := getContainerMounts(o.BaseOpt, tmpfsContainerName)
				verifyMountsInfo(actualMount, expectedMount)
				// create a file in tmpfs mount and verify it doesn't exist after stopping and restarting it
				command.Run(o.BaseOpt, "exec", tmpfsContainerName, "sh", "-c", fmt.Sprintf("echo foo > %s/bar.txt", destDir))
				command.Run(o.BaseOpt, "kill", tmpfsContainerName) // have to use kill to stop the container running with sleep infinity
				command.Run(o.BaseOpt, "start", tmpfsContainerName)
				command.RunWithoutSuccessfulExit(o.BaseOpt, "exec", tmpfsContainerName, "sh", "-c", fmt.Sprintf("cat %s/bar.txt", destDir))
			})

			ginkgo.It("should create a bind mount in a container", func() {
				file := ffs.CreateTempFile("bar.txt", "foo")
				fileDir := filepath.Dir(file)
				ginkgo.DeferCleanup(os.RemoveAll, fileDir)
				command.Run(o.BaseOpt, "run", "-d", "--name", testContainerName, "--mount",
					fmt.Sprintf("type=bind,source=%s,target=%s", fileDir, destDir),
					defaultImage, "sleep", "infinity")
				expectedMount := []MountJSON{makeMount(bindType, fileDir, destDir, "", true)}
				actualMount := getContainerMounts(o.BaseOpt, testContainerName)
				verifyMountsInfo(actualMount, expectedMount)
				output := command.StdoutStr(o.BaseOpt, "exec", testContainerName, "cat", fmt.Sprintf("%s/bar.txt", destDir))
				gomega.Expect(output).Should(gomega.Equal("foo"))
			})

			ginkgo.It("should set the bind mount as readonly with --mount <src>=/src,<target>=/target,ro", func() {
				file := ffs.CreateTempFile("bar.txt", "foo")
				fileDir := filepath.Dir(file)
				ginkgo.DeferCleanup(os.RemoveAll, fileDir)
				cmd := []byte(fmt.Sprintf("echo hello > %s/world.txt", destDir))
				// verify the bind mount is readonly by piping the command of creating a file in the interactive mode to the container
				command.New(o.BaseOpt, "run", "-i", "--name", testContainerName, "--mount",
					fmt.Sprintf("type=bind,source=%s,target=%s,ro", fileDir, destDir),
					defaultImage).WithStdin(gbytes.BufferWithBytes(cmd)).WithoutSuccessfulExit().Run()
				expectedMount := []MountJSON{makeMount(bindType, fileDir, destDir, "ro", false)}
				actualMount := getContainerMounts(o.BaseOpt, testContainerName)
				verifyMountsInfo(actualMount, expectedMount)
			})

			ginkgo.It("should create a tmpfs mount using --mount type=tmpfs flag", func() {
				tmpfsDir := "/tmpfsDir"
				command.Run(o.BaseOpt, "run", "-d", "--name", testContainerName, "--mount",
					fmt.Sprintf("type=tmpfs,destination=%s,tmpfs-mode=1770,tmpfs-size=64m", tmpfsDir),
					defaultImage, "sleep", "infinity")
				expectedMount := []MountJSON{makeMount(tmpfsType, tmpfsType, tmpfsDir, "mode=1770,size=64m", true)}
				actualMount := getContainerMounts(o.BaseOpt, testContainerName)
				verifyMountsInfo(actualMount, expectedMount)
				// create a file in tmpfs mount and verify it doesn't exist after stopping and restarting it
				command.Run(o.BaseOpt, "exec", testContainerName, "sh", "-c", fmt.Sprintf("echo foo > %s/bar.txt", tmpfsDir))
				command.Run(o.BaseOpt, "kill", testContainerName) // have to use kill to stop the container running with sleep infinity
				command.Run(o.BaseOpt, "start", testContainerName)
				command.RunWithoutSuccessfulExit(o.BaseOpt, "exec", testContainerName, "sh", "-c", fmt.Sprintf("cat %s/bar.txt", tmpfsDir))
			})

			ginkgo.It("should mount a volume using --mount type=volume flag", func() {
				command.Run(o.BaseOpt, "run", "--name", testContainerName, "--mount",
					fmt.Sprintf("type=volume,source=%s,target=%s", testVolumeName, destDir), defaultImage)
				srcDir := command.StdoutStr(o.BaseOpt, "volume", "inspect", testVolumeName, "--format", "{{.Mountpoint}}")
				expectedMount := []MountJSON{makeMount(volumeType, srcDir, destDir, "", true)}
				actualMount := getContainerMounts(o.BaseOpt, testContainerName)
				verifyMountsInfo(actualMount, expectedMount)
			})
		})

		// Cgroup version v2
		ginkgo.When("running a container with resource flags", func() {
			ginkgo.BeforeEach(func() {
				if o.CGMode != Unified {
					ginkgo.Skip("requires cgroup v2 to test resource flags")
				}
			})
			ginkgo.It("should set number of CPUs with --cpus flag", func() {
				cpuMax := command.StdoutStr(o.BaseOpt, "run", "--cpus", "0.42", "-w", "/sys/fs/cgroup", defaultImage, "cat", "cpu.max")
				gomega.Expect(cpuMax).To(gomega.Equal("42000 100000"))
			})

			ginkgo.It("should limit CPU CFS (Completely Fair Scheduler) quota and period with --cpu-quota and --cpu-period flags", func() {
				cpuMax := command.StdoutStr(o.BaseOpt, "run", "--cpu-quota", "42000",
					"--cpu-period", "100000", "-w", "/sys/fs/cgroup", defaultImage, "cat", "cpu.max")
				gomega.Expect(cpuMax).To(gomega.Equal("42000 100000"))
			})

			ginkgo.It("should set the CPU shares with --cpu-shares flag", func() {
				// CgroupV2 CPUShares => weight := 1 + ((shares-2)*9999)/262142
				//nolint: lll // the source of the CPUShares calculation formula
				// Ref. https://github.com/google/cadvisor/blob/ce07bb28eadc18183df15ca5346293af6b020b33/integration/tests/api/docker_test.go#L216-L222
				cpuWeight := command.StdoutStr(o.BaseOpt, "run", "--rm", "--cpu-shares", "2000",
					"-w", "/sys/fs/cgroup", defaultImage, "cat", "cpu.weight")
				gomega.Expect(cpuWeight).To(gomega.Equal("77"))
			})

			// assume the host has at least 2 CPUs
			ginkgo.When("--cpuset-cpus is used", func() {
				ginkgo.It("should set CPUs in which to allow execution as cpu 1 with --cpuset-cpus=1", func() {
					cpuSet := command.StdoutStr(o.BaseOpt, "run", "--cpuset-cpus", "1",
						"-w", "/sys/fs/cgroup", defaultImage, "cat", "cpuset.cpus")
					gomega.Expect(cpuSet).To(gomega.Equal("1"))
				})

				ginkgo.It("should set CPUs in which to allow execution as cpu 0-1 with --cpuset-cpus=0-1", func() {
					cpuSet := command.StdoutStr(o.BaseOpt, "run", "--cpuset-cpus", "0-1",
						"-w", "/sys/fs/cgroup", defaultImage, "cat", "cpuset.cpus")
					gomega.Expect(cpuSet).To(gomega.Equal("0-1"))
				})

				ginkgo.It("should set CPUs in which to allow execution as cpu 0-1 with --cpuset-cpus=0,1", func() {
					cpuSet := command.StdoutStr(o.BaseOpt, "run", "--cpuset-cpus", "0,1",
						"-w", "/sys/fs/cgroup", defaultImage, "cat", "cpuset.cpus")
					gomega.Expect(cpuSet).To(gomega.Equal("0-1"))
				})
			})

			// The range form (e.g., `--cpuset-mems 0-1`) is not tested because maybe there's only one memory node (i.e., `0`) on the host,
			// and if that's the case, any number other than 0 would incur an error:
			// https://man7.org/linux/man-pages/man7/cpuset.7.html#WARNINGS
			ginkgo.It("should set memory nodes (MEMs) in which to allow execution as memory node 0 with --cpuset-mems=0", func() {
				cpuMems := command.StdoutStr(o.BaseOpt, "run", "--cpuset-mems", "0",
					"-w", "/sys/fs/cgroup", defaultImage, "cat", "cpuset.mems")
				gomega.Expect(cpuMems).To(gomega.Equal("0"))
			})

			ginkgo.It("should set the memory limit with --memory", func() {
				mem := command.StdoutStr(o.BaseOpt, "run", "--memory", "42m",
					"-w", "/sys/fs/cgroup", defaultImage, "cat", "memory.max")
				gomega.Expect(mem).To(gomega.Equal("44040192"))
			})

			ginkgo.It("should set the memory soft limit with --memory-reservation", func() {
				mem := command.StdoutStr(o.BaseOpt, "run", "--memory-reservation", "6m",
					"-w", "/sys/fs/cgroup", defaultImage, "cat", "memory.low")
				gomega.Expect(mem).To(gomega.Equal("6291456"))
			})

			ginkgo.It("should set the amount of memory this container is allowed to swap to disk with --memory-swap", func() {
				mem := command.StdoutStr(o.BaseOpt, "run", "--memory", "42m", "--memory-swap", "100m",
					"-w", "/sys/fs/cgroup", defaultImage, "cat", "memory.max", "memory.swap.max")
				gomega.Expect(mem).To(gomega.Equal("44040192\n60817408"))
			})

			ginkgo.It("should set the container pids limit with --pids-limit", func() {
				pidsLimit := command.StdoutStr(o.BaseOpt, "run", "--pids-limit", "42",
					"-w", "/sys/fs/cgroup", defaultImage, "cat", "pids.max")
				gomega.Expect(pidsLimit).To(gomega.Equal("42"))
			})

			ginkgo.It("should set the container OOM score with --oom-score-adj", func() {
				// 100 is the minimum because finch VM runs rootless containerd inside.
				testcases := []string{"100", "1000"}
				for _, tc := range testcases {
					score := command.StdoutStr(o.BaseOpt, "run", "--oom-score-adj", tc, defaultImage, "cat", "/proc/self/oom_score_adj")
					gomega.Expect(score).To(gomega.Equal(tc))
				}
			})

			// TODO: --oom-kill-disable --blkio-weight --cgroupns
			// `--device` is not tested because we're not sure what host devices are available
			// as the tests here are OS-agnostic.
			// `--memory-swappiness` is not tested because /sys/fs/cgroup/memory/memory.swappiness is only available in cgroup v1
			// which we are not supporting now
			// `--oom-kill-disable` is not supported until the GitHub issue is resolved
			// https://github.com/containerd/nerdctl/issues/1520
			// `--blkio-weight` is not tested until next lima release https://github.com/containerd/nerdctl/issues/1514
		})

		for _, user := range []string{"-u", "--user"} {
			user := user
			ginkgo.It(fmt.Sprintf("should set the user of a container with %s flag", user), func() {
				// Ref: https://wiki.gentoo.org/wiki/UID_GID_Assignment_Table
				testCases := map[string]string{
					"65534":        "uid=65534(nobody) gid=65534(nobody)",
					"nobody":       "uid=65534(nobody) gid=65534(nobody)",
					"nobody:users": "uid=65534(nobody) gid=100(users)",
					"nobody:100":   "uid=65534(nobody) gid=100(users)",
				}
				for userStr, expected := range testCases {
					output := command.StdoutStr(o.BaseOpt, "run", user, userStr, defaultImage, "id")
					gomega.Expect(output).To(gomega.Equal(expected))
				}
			})

			ginkgo.It("should add additional groups for a specific user with --group-add flag", func() {
				testCases := []struct {
					groups   []string
					expected string
				}{
					{
						groups:   []string{"users"},
						expected: "uid=65534(nobody) gid=65534(nobody) groups=100(users)",
					},
					{
						groups:   []string{"100"},
						expected: "uid=65534(nobody) gid=65534(nobody) groups=100(users)",
					},
					{
						groups:   []string{"users", "nogroup"},
						expected: "uid=65534(nobody) gid=65534(nobody) groups=100(users),65533(nogroup)",
					},
				}

				for _, tc := range testCases {
					args := []string{"run", user, "nobody"}
					for _, group := range tc.groups {
						args = append(args, "--group-add", group)
					}
					args = append(args, defaultImage, "id")
					output := command.StdoutStr(o.BaseOpt, args...)
					gomega.Expect(output).To(gomega.Equal(tc.expected))
				}
			})
		}
	})
}

// MountJSON is used to parse the json from the command images inspect --format "{{json .Mounts}}".
type MountJSON struct {
	MountType   string `json:"Type"`
	Source      string `json:"Source"`
	Destination string `json:"Destination"`
	Mode        string `json:"Mode"`
	Rw          bool   `json:"RW"`
}

func makeMount(mountType string, source string, destination string, mode string, rw bool) MountJSON {
	return MountJSON{
		MountType:   mountType,
		Source:      source,
		Destination: destination,
		Mode:        mode,
		Rw:          rw,
	}
}

// parse the Mounts[] section of the inspect command output JSON to an array []MountJSON.
func getContainerMounts(o *option.Option, containerName string) []MountJSON {
	mountsJSON := command.StdoutStr(o, "inspect", "--format", "{{json .Mounts}}", containerName)
	var mountArray []MountJSON
	err := json.Unmarshal([]byte(mountsJSON), &mountArray)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	return mountArray
}

func verifyMountsInfo(actual []MountJSON, want []MountJSON) {
	gomega.Expect(len(actual)).To(gomega.Equal(len(want)))
	for i, a := range actual {
		w := want[i]
		gomega.Expect(a.MountType).Should(gomega.Equal(w.MountType))
		gomega.Expect(a.Source).Should(gomega.Equal(w.Source))
		gomega.Expect(a.Destination).Should(gomega.Equal(w.Destination))
		gomega.Expect(a.Rw).Should(gomega.Equal(w.Rw))
		if len(w.Mode) > 0 {
			gomega.Expect(strings.Split(a.Mode, ",")).Should(gomega.ContainElements(strings.Split(w.Mode, ",")))
		}
	}
}
