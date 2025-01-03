// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package option

import "testing"

func TestSupportsEnvVarPassthrough(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mods   []Modifier
		assert func(*testing.T, *Option)
	}{
		{
			name: "IsEnvVarPassthroughByDefault",
			mods: []Modifier{},
			assert: func(t *testing.T, uut *Option) {
				if !uut.SupportsEnvVarPassthrough() {
					t.Fatal("expected SupportsEnvVarPassthrough to be true")
				}
			},
		},
		{
			name: "IsNotEnvVarPassthrough",
			mods: []Modifier{
				WithNoEnvironmentVariablePassthrough(),
			},
			assert: func(t *testing.T, uut *Option) {
				if uut.SupportsEnvVarPassthrough() {
					t.Fatal("expected SupportsEnvVarPassthrough to be false")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			uut, err := New([]string{"nerdctl"}, test.mods...)
			if err != nil {
				t.Fatal(err)
			}

			test.assert(t, uut)
		})
	}
}

func TestNerdctlVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mods   []Modifier
		assert func(*testing.T, *Option)
	}{
		{
			name: "IsNerdctlV2ByDefault",
			mods: []Modifier{},
			assert: func(t *testing.T, uut *Option) {
				if !uut.IsNerdctlV2() {
					t.Fatal("expected IsNerdctlV2 to be true")
				}
			},
		},
		{
			name: "IsNerdctlV1",
			mods: []Modifier{
				WithNerdctlVersion("1.7.7"),
			},
			assert: func(t *testing.T, uut *Option) {
				if !uut.IsNerdctlV1() {
					t.Fatal("expected IsNerdctlV1 to be true")
				}
			},
		},
		{
			name: "IsNerdctlV2",
			mods: []Modifier{
				WithNerdctlVersion("2.0.2"),
			},
			assert: func(t *testing.T, uut *Option) {
				if !uut.IsNerdctlV2() {
					t.Fatal("expected IsNerdctlV2 to be true")
				}
			},
		},
		{
			name: "IsPatchedNerdctlV2",
			mods: []Modifier{
				WithNerdctlVersion("2.0.2.m"),
			},
			assert: func(t *testing.T, uut *Option) {
				if !uut.IsNerdctlV2() {
					t.Fatal("expected IsNerdctlV2 to be true")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			uut, err := New([]string{"nerdctl"}, test.mods...)
			if err != nil {
				t.Fatal(err)
			}

			test.assert(t, uut)
		})
	}
}
