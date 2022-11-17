// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package fnet contains functions for network operations.
package fnet

import (
	"net"

	"github.com/onsi/gomega"
)

// GetFreePort returns a free port.
func GetFreePort() int {
	l, err := net.Listen("tcp", "localhost:0")
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	defer func() {
		gomega.Expect(l.Close()).To(gomega.Succeed())
	}()

	tcpAddr, ok := l.Addr().(*net.TCPAddr)
	gomega.Expect(ok).To(gomega.BeTrue())
	return tcpAddr.Port
}
