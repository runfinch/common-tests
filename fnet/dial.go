// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package fnet

import (
	"net/http"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// HTTPGetAndAssert sends an HTTP GET request to the specified URL, asserts the response status code against want, and closes the response body.
func HTTPGetAndAssert(url string, want int, maxRetry int, retryInterval time.Duration) {
	var (
		err  error
		resp *http.Response
	)
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	for i := 0; i < maxRetry; i++ {
		// #nosec G107 // it does not matter if url is not a constant for testing.
		resp, err = client.Get(url)
		if err != nil {
			time.Sleep(retryInterval)
			continue
		}
		defer func() { gomega.Expect(resp.Body.Close()).To(gomega.Succeed()) }()
		gomega.Expect(resp.StatusCode).To(gomega.Equal(want))
		return
	}
	ginkgo.Fail(err.Error())
}
