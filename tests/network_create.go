// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"fmt"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/option"
)

// NetworkCreate tests the "network create" command that creates a network.
func NetworkCreate(o *option.Option) {
	ginkgo.Describe("create a network", func() {
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})
		// TODO: add tests for --ipam-opt, --opt=parent=<INTERFACE>
		ginkgo.It("should create a bridge network", func() {
			command.Run(o, "network", "create", testNetwork)
			gomega.Expect(command.StdoutStr(o, "network", "inspect", testNetwork, "--format", "{{.Name}}")).To(gomega.Equal(testNetwork))
		})

		ginkgo.It("containers under the same network can communicate with each other", func() {
			command.Run(o, "run", "-d", "--name", testContainerName, defaultImage, "sh", "-c", "echo hello | nc -l -p 80")
			ipAddr := command.StdoutStr(o, "inspect", "--format", "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}", testContainerName)
			output := command.StdoutStr(o, "run", defaultImage, "nc", fmt.Sprintf("%s:80", ipAddr))
			gomega.Expect(output).Should(gomega.Equal("hello"))
		})

		ginkgo.It("should create a network with custom subnet using --subnet flag", func() {
			// Choosing 10.5.0.0/16 because is mentioned in the doc - https://github.com/containerd/nerdctl#whale-nerdctl-network-create,
			// so it shouldn't overlap the subnets of the default networks.
			const subnet = "10.5.0.0/16"
			command.Run(o, "network", "create", "--subnet", subnet, testNetwork)
			output := command.StdoutStr(o, "network", "inspect", testNetwork, "--format", "{{(index .IPAM.Config 0).Subnet}}")
			gomega.Expect(output).Should(gomega.Equal(subnet))
		})

		ginkgo.It("should create a network with custom gateway using --gateway flag", func() {
			const (
				subnet  = "10.5.0.0/16"
				gateway = "10.5.0.3"
			)
			command.Run(o, "network", "create", "--subnet", subnet, "--gateway", gateway, testNetwork)
			output := command.StdoutStr(o, "network", "inspect", testNetwork, "--format", "{{(index .IPAM.Config 0).Gateway}}")
			gomega.Expect(output).Should(gomega.Equal(gateway))
		})

		ginkgo.It("should create a network with custom ip range using --ip-range flag", func() {
			const (
				subnet  = "10.5.0.0/16"
				ipRange = "10.5.1.1/32"
			)
			command.Run(o, "network", "create", "--subnet", subnet, "--ip-range", ipRange, testNetwork)
			output := command.StdoutStr(o, "network", "inspect", testNetwork, "--format", "{{(index .IPAM.Config 0).IPRange}}")
			gomega.Expect(output).Should(gomega.Equal(ipRange))

			command.Run(o, "run", "-d", "--name", testContainerName, "--network", testNetwork, defaultImage, "sleep", "infinity")
			ipAddr := command.StdoutStr(o, "inspect", testContainerName, "--format", "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}")
			// Must be 10.5.1.1 because there is only one IP in the IP range.
			gomega.Expect(ipAddr).Should(gomega.Equal("10.5.1.1"))
			// Must fail because there is no available IP in the IP range now.
			command.RunWithoutSuccessfulExit(o, "run", "--name", "test-ctr2", "--network", testNetwork, defaultImage)
		})

		ginkgo.It("should create a network with label using --label flag", func() {
			command.Run(o, "network", "create", "--label", "key=val", testNetwork)
			output := command.StdoutStr(o, "network", "inspect", testNetwork, "--format", "{{.Labels.key}}")
			gomega.Expect(output).Should(gomega.Equal("val"))
		})

		for _, driverFlag := range []string{"-d", "--driver"} {
			driverFlag := driverFlag
			for _, driver := range []string{"macvlan", "ipvlan"} {
				driver := driver
				ginkgo.It(fmt.Sprintf("should create %s network with %s flag", driver, driverFlag), func() {
					command.Run(o, "network", "create", driverFlag, driver, testNetwork)
					netType := command.StdoutStr(o, "network", "inspect", testNetwork, "--mode=native",
						"--format", "{{(index .CNI.plugins 0).type}}")
					gomega.Expect(netType).Should(gomega.Equal(driver))
				})
			}
		}

		for _, opt := range []string{"-o", "--opt"} {
			opt := opt
			ginkgo.It(fmt.Sprintf("should set the containers network MTU with %s flag", opt), func() {
				command.Run(o, "network", "create", opt, "com.docker.network.driver.mtu=500", testNetwork)
				mtu := command.StdoutStr(o, "network", "inspect", testNetwork, "--mode=native", "--format", "{{(index .CNI.plugins 0).mtu}}")
				gomega.Expect(mtu).Should(gomega.Equal("500"))
			})

			ginkgo.It(fmt.Sprintf("should set macvlan network mode to bridge with %s flag", opt), func() {
				command.Run(o, "network", "create", opt, "macvlan_mode=bridge", "-d", "macvlan", testNetwork)
				mode := command.StdoutStr(o, "network", "inspect", testNetwork, "--mode=native", "--format", "{{(index .CNI.plugins 0).mode}}")
				gomega.Expect(mode).Should(gomega.Equal("bridge"))
			})

			ginkgo.It(fmt.Sprintf("should set ipvlan network mode to l3 with %s flag", opt), func() {
				command.Run(o, "network", "create", opt, "ipvlan_mode=l3", "-d", "ipvlan", testNetwork)
				mode := command.StdoutStr(o, "network", "inspect", testNetwork, "--mode=native", "--format", "{{(index .CNI.plugins 0).mode}}")
				gomega.Expect(mode).Should(gomega.Equal("l3"))
			})
		}

		ginkgo.It("should set IPAM driver with --ipam-driver flag", func() {
			command.Run(o, "network", "create", "--ipam-driver=default", testNetwork)
			driverType := command.StdoutStr(o, "network", "inspect", testNetwork, "--mode=native",
				"--format", "{{(index .CNI.plugins 0).ipam.type}}")
			// In unix, default driver type is host-local.
			// https://github.com/containerd/nerdctl/blob/817d6ec27c01986f9cd16a65380294087ef8905f/pkg/netutil/netutil_unix.go#L162
			gomega.Expect(driverType).Should(gomega.Equal("host-local"))
		})
	})
}
