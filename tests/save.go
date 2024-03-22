// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"archive/tar"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/ffs"
	"github.com/runfinch/common-tests/option"
)

type imageManifest struct {
	Layers   []string
	RepoTags []string
}

// Save tests saving an image to a tar archive.
func Save(o *option.Option) {
	ginkgo.Describe("save an image", func() {
		var tarFilePath string
		var tarFileContext string
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
			tarFilePath = ffs.CreateTarFilePath()
			tarFileContext = filepath.Join(tarFilePath, "../")
			ginkgo.DeferCleanup(os.RemoveAll, tarFileContext)
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})

		ginkgo.Context("when the images exist", func() {
			ginkgo.BeforeEach(func() {
				pullImage(o, localImages["defaultImage"])
			})

			ginkgo.It("should save an image to stdout", func() {
				stdout := command.New(o, "save", localImages["defaultImage"]).WithStdout(gbytes.NewBuffer()).Run().Out
				untar(stdout, tarFileContext)
				manifestContent := readManifestContent(tarFileContext)

				gomega.Expect(len(manifestContent)).Should(gomega.Equal(1))
				gomega.Expect(manifestContent[0].RepoTags[0]).Should(gomega.Equal(localImages["defaultImage"]))

				layersShouldExist(manifestContent[0].Layers, tarFileContext)
			})

			for _, outputOption := range []string{"-o", "--output"} {
				outputOption := outputOption
				ginkgo.It(fmt.Sprintf("should save an image with %s option", outputOption), func() {
					command.Run(o, "save", localImages["defaultImage"], outputOption, tarFilePath)

					untarFile(tarFilePath, tarFileContext)

					manifestContent := readManifestContent(tarFileContext)
					gomega.Expect(len(manifestContent)).Should(gomega.Equal(1))
					gomega.Expect(manifestContent[0].RepoTags[0]).Should(gomega.Equal(localImages["defaultImage"]))
					layersShouldExist(manifestContent[0].Layers, tarFileContext)
				})

				ginkgo.It(fmt.Sprintf("should save multiple images with %s option", outputOption), func() {
					pullImage(o, localImages["olderAlpineImage"])
					command.Run(o, "save", outputOption, tarFilePath, localImages["defaultImage"], localImages["olderAlpineImage"])

					untarFile(tarFilePath, tarFileContext)

					manifestContent := readManifestContent(tarFileContext)
					gomega.Expect(len(manifestContent)).Should(gomega.Equal(2))

					for i := range manifestContent {
						layersShouldExist(manifestContent[i].Layers, tarFileContext)
					}
				})
			}
		})

		ginkgo.It("should not be able to save an image if the image doesn't exist", func() {
			for _, outputOption := range []string{"-o", "--output"} {
				command.RunWithoutSuccessfulExit(o, "save", outputOption, tarFilePath, nonexistentImageName)
				command.RunWithoutSuccessfulExit(o, "save", outputOption, tarFilePath, nonexistentImageName)
			}
		})
	})
}

func layersShouldExist(layers []string, dir string) {
	for _, l := range layers {
		layerPath := filepath.Join(dir, l)
		_, err := os.Stat(layerPath)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	}
}

func untarFile(tarFilePath, tarFileContext string) {
	reader, err := os.Open(filepath.Clean(tarFilePath))
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	defer func() {
		gomega.Expect(reader.Close()).Should(gomega.Succeed())
	}()
	untar(reader, tarFileContext)
}

func untar(reader io.Reader, targetDir string) {
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		path := filepath.Join(targetDir, header.Name)                        //nolint:gosec // The following line resolves G305.
		gomega.Expect(path).To(gomega.HavePrefix(filepath.Clean(targetDir))) // https://security.snyk.io/research/zip-slip-vulnerability
		info := header.FileInfo()

		if info.IsDir() {
			gomega.Expect(os.MkdirAll(path, info.Mode())).Should(gomega.Succeed())
			continue
		}

		file, err := os.OpenFile(filepath.Clean(path), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		//nolint:gosec // Using io.CopyN to fix G110 seems to be an overkill considering the attack possibility.
		_, err = io.Copy(file, tarReader)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		gomega.Expect(file.Close()).Should(gomega.Succeed())
	}
}

func readManifestContent(tarFileContext string) []imageManifest {
	manifestFilePath := filepath.Join(tarFileContext, "manifest.json")
	manifestBytes, err := os.ReadFile(filepath.Clean(manifestFilePath))
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	var manifestContent []imageManifest
	err = json.Unmarshal(manifestBytes, &manifestContent)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	return manifestContent
}
