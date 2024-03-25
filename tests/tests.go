// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package tests contains the exported functions that are meant to be imported as test cases.
//
// It should not export any other thing except for a SubcommandOption struct (e.g., LoginOption) that may be added in the future.
//
// Each file contains one subcommand to test and is named after that subcommand.
// Note that the file names are not suffixed with _test so that they can appear in Go Doc.
package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/runfinch/common-tests/fnet"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/ffs"
	"github.com/runfinch/common-tests/option"
)

const (
	alpineImage              = "public.ecr.aws/docker/library/alpine:latest"
	testImageName            = "test:tag"
	nonexistentImageName     = "ne-repo:ne-tag"
	nonexistentContainerName = "ne-ctr"
	testContainerName        = "ctr-test"
	testContainerName2       = "ctr-test-2"
	testVolumeName           = "testVol"
	registryImage            = "public.ecr.aws/docker/library/registry:latest"
	localRegistryName        = "local-registry"
	testUser                 = "testUser"
	testPassword             = "testPassword"
	sha256RegexFull          = "^sha256:[a-f0-9]{64}$"
	bridgeNetwork            = "bridge"
	testNetwork              = "test-network"
)

type localImage string

const (
	defaultImage      localImage = "defaultImage"
	olderAlpineImage  localImage = "olderAlpineImage"
	amazonLinux2Image localImage = "amazonLinux2Image"
	nginxImage        localImage = "nginxImage"
)

var remoteImages = map[localImage]string{
	defaultImage:      alpineImage,
	olderAlpineImage:  "public.ecr.aws/docker/library/alpine:3.13",
	amazonLinux2Image: "public.ecr.aws/amazonlinux/amazonlinux:2",
	nginxImage:        "public.ecr.aws/docker/library/nginx:latest",
}

var localImages = map[localImage]string{}

// CGMode is the cgroups mode of the host system.
// We copy the struct from containerd/cgroups [1] instead of using it as a library
// because it only builds on linux,
// while we don't really need the functions that make it only build on linux
// (e.g., determine the cgroup version of the current host).
//
// [1] https://github.com/containerd/cgroups/blob/cc78c6c1e32dc5bde018d92999910fdace3cfa27/utils.go#L38-L50
type CGMode int

const (
	// Unavailable cgroup mountpoint.
	Unavailable CGMode = iota
	// Legacy cgroups v1.
	Legacy
	// Hybrid with cgroups v1 and v2 controllers mounted.
	Hybrid
	// Unified with only cgroups v2 mounted.
	Unified
)

// SetupLocalRegistry can be invoked before running the tests to save time when pulling images during tests.
//
// It spins up a local registry, tags all remoteImages, pushes the new tagged images to the local registry,
// and changes adds corresponding entries to localImages for all of the new tags pushed to local registry.
//
// After all the tests are done, invoke CleanupLocalRegistry to clean up the local registry.
func SetupLocalRegistry(o *option.Option) {
	command.RemoveAll(o)
	hostPort := fnet.GetFreePort()
	containerID := command.StdoutStr(o, "run", "-d", "-p",
		fmt.Sprintf("%d:5000", hostPort), "--name", localRegistryName, registryImage)
	imageID := command.StdoutStr(o, "images", "-q")
	command.SetLocalRegistryContainerID(containerID)
	command.SetLocalRegistryImageID(imageID)
	command.SetLocalRegistryImageName(registryImage)

	for k, ref := range remoteImages {
		// split image tag according to spec
		// https://github.com/distribution/distribution/blob/d0deff9cd6c2b8c82c6f3d1c713af51df099d07b/reference/reference.go
		_, name, _ := strings.Cut(ref, "/")
		command.Run(o, "pull", ref)
		localRef := fmt.Sprintf("localhost:%d/%s", hostPort, name)
		command.Run(o, "tag", ref, localRef)
		command.Run(o, "push", localRef)
		command.Run(o, "rmi", ref)
		localImages[k] = localRef
	}
}

// CleanupLocalRegistry removes the local registry container and image. It's used together with SetupLocalRegistry,
// and should be invoked after running all the tests.
func CleanupLocalRegistry(o *option.Option) {
	containerID := command.StdoutStr(o, "inspect", localRegistryName, "--format", "{{.ID}}")
	command.Run(o, "rm", "-f", containerID)
	imageID := command.StdoutStr(o, "images", "-q")
	command.Run(o, "rmi", "-f", imageID)
	localImages = map[localImage]string{}
}

func pullImage(o *option.Option, imageName string) {
	command.Run(o, "pull", "-q", imageName)
	imageID := command.Stdout(o, "images", "--quiet", imageName)
	gomega.Expect(imageID).ShouldNot(gomega.BeEmpty())
}

func removeImage(o *option.Option, imageName string) {
	command.Run(o, "rmi", "--force", imageName)
	imageID := command.Stdout(o, "images", "--quiet", imageName)
	gomega.Expect(string(imageID)).Should(gomega.BeEmpty())
}

func containerShouldBeRunning(o *option.Option, containerNames ...string) {
	for _, containerName := range containerNames {
		gomega.Expect(command.Stdout(o, "ps", "-q", "--filter",
			fmt.Sprintf("name=%s", containerName))).NotTo(gomega.BeEmpty())
	}
}

func containerShouldNotBeRunning(o *option.Option, containerNames ...string) {
	for _, containerName := range containerNames {
		gomega.Expect(command.Stdout(o, "ps", "-q", "--filter",
			fmt.Sprintf("name=%s", containerName))).To(gomega.BeEmpty())
	}
}

func containerShouldExist(o *option.Option, containerNames ...string) {
	for _, containerName := range containerNames {
		gomega.Expect(command.Stdout(o, "ps", "-a", "-q", "--filter",
			fmt.Sprintf("name=%s", containerName))).NotTo(gomega.BeEmpty())
	}
}

func containerShouldNotExist(o *option.Option, containerNames ...string) error {
	for _, containerName := range containerNames {
		containerExists := command.Stdout(o, "ps", "-a", "-q", "--filter",
			fmt.Sprintf("name=%s", containerName))
		if len(containerExists) > 0 {
			return fmt.Errorf("containerd '%s' exists but should not", containerName)
		}
	}
	return nil
}

func imageShouldExist(o *option.Option, imageName string) {
	gomega.Expect(command.Stdout(o, "images", "-q", imageName)).NotTo(gomega.BeEmpty())
}

func imageShouldNotExist(o *option.Option, imageName string) {
	gomega.Expect(command.Stdout(o, "images", "-q", imageName)).To(gomega.BeEmpty())
}

func volumeShouldExist(o *option.Option, volumeName string) {
	gomega.Expect(command.Stdout(o, "volume", "ls", "-q", "--filter",
		fmt.Sprintf("name=%s", volumeName))).NotTo(gomega.BeEmpty())
}

func volumeShouldNotExist(o *option.Option, volumeName string) {
	gomega.Expect(command.Stdout(o, "volume", "ls", "-q", "--filter",
		fmt.Sprintf("name=%s", volumeName))).To(gomega.BeEmpty())
}

func fileShouldExist(path, content string) {
	gomega.Expect(path).To(gomega.BeARegularFile())
	actualContent, err := os.ReadFile(filepath.Clean(path))
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(string(actualContent)).To(gomega.Equal(content))
}

func fileShouldNotExist(path string) {
	gomega.Expect(path).ToNot(gomega.BeAnExistingFile())
}

func fileShouldExistInContainer(o *option.Option, containerName, path, content string) {
	gomega.Expect(command.StdoutStr(o, "exec", containerName, "cat", path)).To(gomega.Equal(content))
}

func fileShouldNotExistInContainer(o *option.Option, containerName, path string) {
	cmdOut := command.RunWithoutSuccessfulExit(o, "exec", containerName, "cat", path)
	gomega.Expect(cmdOut.Err.Contents()).To(gomega.ContainSubstring("No such file or directory"))
}

func buildImage(o *option.Option, imageName string) {
	dockerfile := fmt.Sprintf(`FROM %s
		CMD ["echo", "finch-test-dummy-output"]
		`, localImages[defaultImage])
	buildContext := ffs.CreateBuildContext(dockerfile)
	ginkgo.DeferCleanup(os.RemoveAll, buildContext)
	command.Run(o, "build", "-q", "-t", imageName, buildContext)
}
