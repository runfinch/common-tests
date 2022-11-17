// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"encoding/json"
	"fmt"
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

// Run tests running a container image.
// TODO: test cases for bind propagation options.
func Run(o *option.Option) {
	ginkgo.Describe("Run a container image", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.When("running a container that echos dummy output", func() {
			ginkgo.BeforeEach(func() {
				dockerfile := fmt.Sprintf(`FROM %s
			CMD ["echo", "finch-test-dummy-output"]
			`, defaultImage)
				buildContext := ffs.CreateBuildContext(dockerfile)
				ginkgo.DeferCleanup(os.RemoveAll, buildContext)
				command.Run(o, "build", "-q", "-t", testImageName, buildContext)
			})

			ginkgo.It("should echo dummy output", func() {
				output := command.StdoutStr(o, "run", testImageName)
				gomega.Expect(output).Should(gomega.Equal("finch-test-dummy-output"))
			})

			ginkgo.It("should not echo dummy output if running with -d flag", func() {
				output := command.StdOut(o, "run", "-d", testImageName)
				gomega.Expect(output).ShouldNot(gomega.ContainSubstring("finch-test-dummy-output"))
			})
		})

		ginkgo.It("with --rm flag, container should be removed when it exits", func() {
			command.Run(o, "run", "--rm", "--name", testContainerName, defaultImage)
			containerShouldNotExist(o, testContainerName)
		})

		ginkgo.When("running a container with metadata related flags", func() {
			for _, label := range []string{"-l", "--label"} {
				label := label
				ginkgo.It(fmt.Sprintf("should set meta data on a container with %s flag", label), func() {
					command.Run(o, "run", "--name", testContainerName, label, "testKey=testValue", defaultImage)
					gomega.Expect(command.StdoutStr(o, "inspect", testContainerName,
						"--format", "{{.Config.Labels.testKey}}")).To(gomega.Equal("testValue"))
				})
			}

			ginkgo.It("should write the container ID to file with --cidfile flag", func() {
				dir := ffs.CreateTempDir("finch-test-cid")
				ginkgo.DeferCleanup(os.RemoveAll, dir)
				path := filepath.Join(dir, "test.cid")
				containerID := command.StdoutStr(o, "run", "-d", "--cidfile", path, defaultImage)
				output, err := os.ReadFile(filepath.Clean(path))
				gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
				gomega.Expect(strings.TrimSpace(string(output))).Should(gomega.Equal(containerID))
			})

			ginkgo.It("should read labels from file with --label-file flag", func() {
				path := ffs.CreateTempFile("label-file", "key=value")
				ginkgo.DeferCleanup(os.RemoveAll, filepath.Dir(path))
				command.Run(o, "run", "--name", testContainerName, "--label-file", path, defaultImage)
				gomega.Expect(command.StdoutStr(o, "inspect", testContainerName,
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
				command.Run(o, "build", "-q", "-t", testImageName, buildContext)

				envOutput := command.StdOut(o, "run", "--rm", "--entrypoint", "time", testImageName, "echo", "blah")
				gomega.Expect(envOutput).NotTo(gomega.ContainSubstring("foo"))
				gomega.Expect(envOutput).NotTo(gomega.ContainSubstring("bar"))
				gomega.Expect(envOutput).To(gomega.ContainSubstring("blah"))
			})

			for _, workDir := range []string{"--workdir", "-w"} {
				ginkgo.It(fmt.Sprintf("should set working directory inside the container specified by %s flag", workDir), func() {
					dir := "/tmp"
					gomega.Expect(command.StdoutStr(o, "run", workDir, dir, defaultImage, "pwd")).Should(gomega.Equal(dir))
				})
			}

			for _, env := range []string{"-e", "--env"} {
				ginkgo.It(fmt.Sprintf("with %s flag, environment variables should be set in the container", env), func() {
					envOutput := command.StdOut(o, "run", "--rm",
						"--env", "FOO=BAR", "--env", "FOO1", "-e", "ENV1=1", "-e", "ENV1=2",
						defaultImage, "env")
					gomega.Expect(envOutput).To(gomega.ContainSubstring("FOO=BAR"))
					gomega.Expect(envOutput).ToNot(gomega.ContainSubstring("FOO1"))
					gomega.Expect(envOutput).To(gomega.ContainSubstring("ENV1=2"))
				})
			}

			ginkgo.It("with --env-file flag, environment variables in container should pick up those in environment file", func() {
				const envPair = "ENVKEY=ENVVAL"
				envPath := ffs.CreateTempFile("env", envPair)
				ginkgo.DeferCleanup(os.RemoveAll, filepath.Dir(envPath))

				envOutput := command.StdOut(o, "run", "--rm", "--env-file", envPath, defaultImage, "env")
				gomega.Expect(envOutput).To(gomega.ContainSubstring(envPair))
			})
		})

		ginkgo.When("running an image with --pull flag", func() {
			ginkgo.It("should have an error if set --pull=never and the image doesn't exist", func() {
				command.RunWithoutSuccessfulExit(o, "run", "--pull", "never", defaultImage)
				imageShouldNotExist(o, defaultImage)
			})

			ginkgo.It("should be able to run the container if set --pull=never and the image exists", func() {
				const containerName = "test-container"
				pullImage(o, defaultImage)
				command.Run(o, "run", "--name", containerName, "--pull", "never", defaultImage)
				containerShouldExist(o, containerName)
			})

			ginkgo.It("should be able to run the container if set --pull=missing and the image doesn't exist", func() {
				command.Run(o, "run", "--name", testContainerName, "--pull", "missing", defaultImage)
				containerShouldExist(o, testContainerName)
				imageShouldExist(o, defaultImage)
			})

			ginkgo.It("should be able to run the container if set --pull=missing and the image exists", func() {
				pullImage(o, defaultImage)
				command.Run(o, "run", "--name", testContainerName, "--pull", "missing", defaultImage)
				containerShouldExist(o, testContainerName)
			})

			ginkgo.It("should be able to run the container if set --pull=always and the image doesn't exist", func() {
				command.Run(o, "run", "--name", testContainerName, "--pull", "always", defaultImage)
				containerShouldExist(o, testContainerName)
				imageShouldExist(o, defaultImage)
			})
			ginkgo.It("should be able to run the container if set --pull=always and the image exists", func() {
				pullImage(o, defaultImage)
				command.Run(o, "run", "--name", testContainerName, "--pull", "always", defaultImage)
			})
		})

		for _, interactive := range []string{"-i", "--interactive"} {
			interactive := interactive
			ginkgo.It(fmt.Sprintf("should output string if %s flag keeps STDIN open", interactive), func() {
				want := []byte("hello")
				got := command.New(o, "run", interactive, defaultImage, "cat").
					WithStdin(gbytes.BufferWithBytes(want)).Run().Out.Contents()
				gomega.Expect(got).Should(gomega.Equal(want))
			})
		}

		ginkgo.It("should stop running container within specified time by --stop-timeout flag", func() {
			// With PID=1, `sleep infinity` does not exit due to receiving a SIGTERM, which is sent by the stop command.
			// Ref. https://superuser.com/a/1299463/730265
			command.Run(o, "run", "-d", "--name", testContainerName, "--stop-timeout", "1", defaultImage, "sleep", "infinity")
			gomega.Expect(command.StdoutStr(o, "exec", testContainerName, "echo", "foo")).To(gomega.Equal("foo"))
			startTime := time.Now()
			command.Run(o, "stop", testContainerName)
			// assert the container to be stopped within 1.5 seconds
			gomega.Expect(time.Since(startTime)).To(gomega.BeNumerically("~", 0*time.Millisecond, 1500*time.Millisecond))
			command.RunWithoutSuccessfulExit(o, "exec", testContainerName, "echo", "foo")
		})

		ginkgo.It("should immediately stop the container with --stop-signal=SIGKILL", func() {
			// With PID=1, `sleep infinity` will only exit when receiving SIGKILL, while the signal sent by stop command is SIGTERM.
			command.Run(o, "run", "-d", "--name", testContainerName, "--stop-signal", "SIGKILL", defaultImage, "sleep", "infinity")
			containerShouldBeRunning(o, testContainerName)
			startTime := time.Now()
			command.Run(o, "stop", testContainerName)
			gomega.Expect(time.Since(startTime)).To(gomega.BeNumerically("~", 0*time.Millisecond, 500*time.Millisecond))
			status := command.StdoutStr(o, "inspect", "--format", "{{.State.Status}}", testContainerName)
			gomega.Expect(status).Should(gomega.Equal("exited"))
			command.RunWithoutSuccessfulExit(o, "exec", testContainerName, "echo", "foo")
		})

		ginkgo.It("should share PID namespace with host with --pid=host", func() {
			command.Run(o, "run", "-d", "--name", testContainerName, "--pid", "host", defaultImage, "sleep", "infinity")
			processes := command.StdoutStr(o, "top", testContainerName)
			pid := strings.Fields(processes)[9]
			output := command.StdOutAsLines(o, "exec", "-d", testContainerName, "ps", "-o", "pid")
			gomega.Expect(output).Should(gomega.ContainElement(pid))
		})

		ginkgo.When("running a container with network related flags", func() {
			// TODO: add tests for --ip, --mac-address flags
			for _, network := range []string{"--net", "--network"} {
				network := network
				ginkgo.It(fmt.Sprintf("should connect a container to a network with %s flag", network), func() {
					command.Run(o, "run", "-d", network, "bridge", "--name", testContainerName,
						defaultImage, "sh", "-c", "echo hello | nc -l -p 80")
					ipAddr := command.StdoutStr(o, "inspect", "--format",
						"{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}", testContainerName)
					output := command.StdoutStr(o, "run", network, "bridge", defaultImage, "nc", fmt.Sprintf("%s:80", ipAddr))
					gomega.Expect(output).Should(gomega.Equal("hello"))
				})

				ginkgo.It(fmt.Sprintf("should use the same network with container specified by %s=container:<name>", network), func() {
					command.Run(o, "run", "-d", network, "bridge", "--name", testContainerName,
						defaultImage, "sh", "-c", "echo hello | nc -l -p 80")
					ipAddr := command.StdoutStr(o, "inspect", "--format",
						"{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}", testContainerName)
					output := command.StdoutStr(o, "run", fmt.Sprintf("%s=container:%s", network, testContainerName),
						defaultImage, "nc", fmt.Sprintf("%s:80", ipAddr))
					gomega.Expect(output).Should(gomega.Equal("hello"))
				})

				ginkgo.It(fmt.Sprintf("should use the same network with container specified by %s=container:<id>", network), func() {
					id := command.StdoutStr(o, "run", "-d", network, "bridge", "--name", testContainerName,
						defaultImage, "sh", "-c", "echo hello | nc -l -p 80")
					ipAddr := command.StdoutStr(o, "inspect", "--format",
						"{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}", testContainerName)
					output := command.StdoutStr(o, "run", fmt.Sprintf("%s=container:%s", network, id),
						defaultImage, "nc", fmt.Sprintf("%s:80", ipAddr))
					gomega.Expect(output).Should(gomega.Equal("hello"))
				})
			}

			ginkgo.It("should be able to set custom DNS servers with --dns flag", func() {
				const nameserver = "10.10.10.10"
				lines := command.StdOutAsLines(o, "run", "--dns", nameserver, "--name", testContainerName,
					defaultImage, "cat", "/etc/resolv.conf")
				gomega.Expect(lines).Should(gomega.ContainElement(fmt.Sprintf("nameserver %s", nameserver)))
			})

			ginkgo.It("should be able to set custom DNS search domains with --dns-search flag", func() {
				lines := command.StdOutAsLines(o, "run", "--dns-search", "test", "--name", testContainerName,
					defaultImage, "cat", "/etc/resolv.conf")
				gomega.Expect(lines).Should(gomega.ContainElement("search test"))
			})

			for _, dnsOption := range []string{"--dns-opt", "--dns-option"} {
				dnsOption := dnsOption
				ginkgo.It(fmt.Sprintf("should be able to set DNS option with %s flag", dnsOption), func() {
					lines := command.StdOutAsLines(o, "run", dnsOption, "debug", "--name", testContainerName,
						defaultImage, "cat", "/etc/resolv.conf")
					gomega.Expect(lines).Should(gomega.ContainElement("options debug"))
				})
			}

			for _, hostname := range []string{"--hostname", "-h"} {
				hostname := hostname
				ginkgo.It(fmt.Sprintf("should be able to set container host name with %s flag", hostname), func() {
					name := command.StdoutStr(o, "run", hostname, "foo", defaultImage, "hostname")
					gomega.Expect(name).Should(gomega.Equal("foo"))
				})
			}

			ginkgo.It("should add a custom host-to-IP mapping with --add-host flag", func() {
				mapping := command.StdoutStr(o, "run", "--add-host", "test-host:6.6.6.6", defaultImage, "cat", "/etc/hosts")
				gomega.Expect(mapping).Should(gomega.ContainSubstring("6.6.6.6"))
				gomega.Expect(mapping).Should(gomega.ContainSubstring("test-host"))
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
					command.Run(o, "run", "-d", publish, fmt.Sprintf("%d:%d", hostPort, containerPort), defaultImage,
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
					command.Run(o, "run", "--name", testContainerName, volume,
						fmt.Sprintf("%s:%s", testVolumeName, destDir), defaultImage, "sh", "-c", "echo foo > /tmp/test.txt")
					srcDir := command.StdoutStr(o, "volume", "inspect", testVolumeName, "--format", "{{.Mountpoint}}")
					expectedMount := []MountJSON{makeMount(volumeType, srcDir, destDir, "", true)}
					actualMount := getContainerMounts(o, testContainerName)
					verifyMountsInfo(actualMount, expectedMount)
					output := command.StdoutStr(o, "run", "-v", fmt.Sprintf("%s:/tmp", testVolumeName), defaultImage, "cat", "/tmp/test.txt")
					gomega.Expect(output).Should(gomega.Equal("foo"))
				})

				ginkgo.It(fmt.Sprintf("should be able to set the volume options with %s testVol:%s:ro", volume, destDir), func() {
					command.Run(o, "run", "-d", "--name", testContainerName, volume,
						fmt.Sprintf("%s:%s:ro", testVolumeName, destDir), defaultImage, "sleep", "infinity")
					srcDir := command.StdoutStr(o, "volume", "inspect", "--format", "{{.Mountpoint}}", testVolumeName)
					expectedMount := []MountJSON{makeMount(volumeType, srcDir, destDir, "ro", false)}
					actualMount := getContainerMounts(o, testContainerName)
					verifyMountsInfo(actualMount, expectedMount)
					// verify the volume is readonly
					command.RunWithoutSuccessfulExit(o, "exec", testContainerName, "sh", "-c", fmt.Sprintf("echo foo > %s/test.txt", destDir))
				})
			}

			ginkgo.It("should create a tmpfs mount in a container", func() {
				const tmpfsContainerName = "tmpfs-ctr"
				command.Run(o, "run", "-d", "--tmpfs", fmt.Sprintf("%s:size=64m,exec", destDir),
					"--name", tmpfsContainerName, defaultImage, "sleep", "infinity")
				expectedMount := []MountJSON{makeMount(tmpfsType, tmpfsType, destDir, "size=64m,exec", true)}
				actualMount := getContainerMounts(o, tmpfsContainerName)
				verifyMountsInfo(actualMount, expectedMount)
				// create a file in tmpfs mount and verify it doesn't exist after stopping and restarting it
				command.Run(o, "exec", tmpfsContainerName, "sh", "-c", fmt.Sprintf("echo foo > %s/bar.txt", destDir))
				command.Run(o, "kill", tmpfsContainerName) // have to use kill to stop the container running with sleep infinity
				command.Run(o, "start", tmpfsContainerName)
				command.RunWithoutSuccessfulExit(o, "exec", tmpfsContainerName, "sh", "-c", fmt.Sprintf("cat %s/bar.txt", destDir))
			})

			ginkgo.It("should create a bind mount in a container", func() {
				file := ffs.CreateTempFile("bar.txt", "foo")
				fileDir := filepath.Dir(file)
				ginkgo.DeferCleanup(os.RemoveAll, fileDir)
				command.Run(o, "run", "-d", "--name", testContainerName, "--mount",
					fmt.Sprintf("type=bind,source=%s,target=%s", fileDir, destDir),
					defaultImage, "sleep", "infinity")
				expectedMount := []MountJSON{makeMount(bindType, fileDir, destDir, "", true)}
				actualMount := getContainerMounts(o, testContainerName)
				verifyMountsInfo(actualMount, expectedMount)
				output := command.StdoutStr(o, "exec", testContainerName, "cat", fmt.Sprintf("%s/bar.txt", destDir))
				gomega.Expect(output).Should(gomega.Equal("foo"))
			})

			ginkgo.It("should set the bind mount as readonly with --mount <src>=/src,<target>=/target,ro", func() {
				file := ffs.CreateTempFile("bar.txt", "foo")
				fileDir := filepath.Dir(file)
				ginkgo.DeferCleanup(os.RemoveAll, fileDir)
				cmd := []byte(fmt.Sprintf("echo hello > %s/world.txt", destDir))
				// verify the bind mount is readonly by piping the command of creating a file in the interactive mode to the container
				command.New(o, "run", "-i", "--name", testContainerName, "--mount",
					fmt.Sprintf("type=bind,source=%s,target=%s,ro", fileDir, destDir),
					defaultImage).WithStdin(gbytes.BufferWithBytes(cmd)).WithoutSuccessfulExit().Run()
				expectedMount := []MountJSON{makeMount(bindType, fileDir, destDir, "ro", false)}
				actualMount := getContainerMounts(o, testContainerName)
				verifyMountsInfo(actualMount, expectedMount)
			})

			ginkgo.It("should create a tmpfs mount using --mount type=tmpfs flag", func() {
				tmpfsDir := "/tmpfsDir"
				command.Run(o, "run", "-d", "--name", testContainerName, "--mount",
					fmt.Sprintf("type=tmpfs,destination=%s,tmpfs-mode=1770,tmpfs-size=64m", tmpfsDir),
					defaultImage, "sleep", "infinity")
				expectedMount := []MountJSON{makeMount(tmpfsType, tmpfsType, tmpfsDir, "mode=1770,size=64m", true)}
				actualMount := getContainerMounts(o, testContainerName)
				verifyMountsInfo(actualMount, expectedMount)
				// create a file in tmpfs mount and verify it doesn't exist after stopping and restarting it
				command.Run(o, "exec", testContainerName, "sh", "-c", fmt.Sprintf("echo foo > %s/bar.txt", tmpfsDir))
				command.Run(o, "kill", testContainerName) // have to use kill to stop the container running with sleep infinity
				command.Run(o, "start", testContainerName)
				command.RunWithoutSuccessfulExit(o, "exec", testContainerName, "sh", "-c", fmt.Sprintf("cat %s/bar.txt", tmpfsDir))
			})

			ginkgo.It("should mount a volume using --mount type=volume flag", func() {
				command.Run(o, "run", "--name", testContainerName, "--mount",
					fmt.Sprintf("type=volume,source=%s,target=%s", testVolumeName, destDir), defaultImage)
				srcDir := command.StdoutStr(o, "volume", "inspect", testVolumeName, "--format", "{{.Mountpoint}}")
				expectedMount := []MountJSON{makeMount(volumeType, srcDir, destDir, "", true)}
				actualMount := getContainerMounts(o, testContainerName)
				verifyMountsInfo(actualMount, expectedMount)
			})
		})
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
