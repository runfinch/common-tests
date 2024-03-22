// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/onsi/ginkgo/v2"

	"github.com/runfinch/common-tests/command"
	"github.com/runfinch/common-tests/ffs"
	"github.com/runfinch/common-tests/option"
)

// Load tests loading images from tar file or stdin.
func Load(o *option.Option) {
	ginkgo.Describe("load an image", func() {
		var tarFilePath string
		ginkgo.BeforeEach(func() {
			command.RemoveAll(o)
			pullImage(o, localImages["defaultImage"])
			tarFilePath = ffs.CreateTarFilePath()
			ginkgo.DeferCleanup(os.RemoveAll, filepath.Join(tarFilePath, "../"))
		})
		ginkgo.AfterEach(func() {
			command.RemoveAll(o)
		})
		// TODO: add test for input redirection sign
		// REF issue: https://github.com/lima-vm/lima/issues/1078
		for _, inputOption := range []string{"-i", "--input"} {
			inputOption := inputOption
			ginkgo.It(fmt.Sprintf("should load an image with %s option", inputOption), func() {
				command.Run(o, "save", "-o", tarFilePath, localImages["defaultImage"])

				command.Run(o, "rmi", localImages["defaultImage"])
				imageShouldNotExist(o, localImages["defaultImage"])

				command.Run(o, "load", inputOption, tarFilePath)
				imageShouldExist(o, localImages["defaultImage"])
			})

			ginkgo.It(fmt.Sprintf("should load multiple images with %s option", inputOption), func() {
				pullImage(o, localImages["olderAlpineImage"])
				command.Run(o, "save", "-o", tarFilePath, localImages["defaultImage"], localImages["olderAlpineImage"])

				command.Run(o, "rmi", localImages["defaultImage"], localImages["olderAlpineImage"])
				imageShouldNotExist(o, localImages["defaultImage"])
				imageShouldNotExist(o, localImages["olderAlpineImage"])

				command.Run(o, "load", inputOption, tarFilePath)
				imageShouldExist(o, localImages["defaultImage"])
				imageShouldExist(o, localImages["olderAlpineImage"])
			})
		}
	})
}
