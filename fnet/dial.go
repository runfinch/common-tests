// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package fnet

import (
	"net"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// DialAndRead dials the network address, reads the data from the established connection, and asserts it against want.
func DialAndRead(network, address, want string, maxRetry int, retryInterval time.Duration) {
	var (
		conn net.Conn
		err  error
	)
	for i := 0; i < maxRetry; i++ {
		conn, err = net.Dial(network, address)
		if err != nil {
			time.Sleep(retryInterval)
			continue
		}

		b := make([]byte, len([]byte(want))) //nolint:makezero // The content of b does not matter,
		// but len(b) must be equal to len([]byte(want)) so that conn.Read can read the whole thing.
		gomega.Expect(conn.Read(b)).Error().ShouldNot(gomega.HaveOccurred())
		gomega.Expect(b).To(gomega.Equal([]byte(want)))
		gomega.Expect(conn.Close()).To(gomega.Succeed())
		return
	}
	ginkgo.Fail(err.Error())
}
