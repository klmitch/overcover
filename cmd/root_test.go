// Copyright (c) 2020 Kevin L. Mitchell
//
// Licensed under the Apache License, Version 2.0 (the "License"); you
// may not use this file except in compliance with the License.  You
// may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.  See the License for the specific language governing
// permissions and limitations under the License.

package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/klmitch/patcher"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestRootCmdBase(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&coverprofile, "./testdata/full_coverage.out"),
	).Install().Restore()

	rootCmd.Run(rootCmd, []string{})

	assert.Equal(t, "19 statements out of 19 covered; overall coverage: 100.0%\n", outStream.String())
	assert.Equal(t, "", errStream.String())
}

func TestRootCmdOpenFails(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&coverprofile, "./testdata/no_such_coverage.out"),
	).Install().Restore()

	assert.PanicsWithValue(t, "os.Exit(2)", func() { rootCmd.Run(rootCmd, []string{}) })
	assert.Equal(t, "", outStream.String())
	assert.Contains(t, errStream.String(), "Unable to open coverage profile file: ")
}

func TestRootCmdLoadCoverageFails(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&coverprofile, "./testdata/bad_coverage.out"),
	).Install().Restore()

	assert.PanicsWithValue(t, "os.Exit(3)", func() { rootCmd.Run(rootCmd, []string{}) })
	assert.Equal(t, "", outStream.String())
	assert.Contains(t, errStream.String(), "Unable to read coverage profile file: ")
}

func TestRootCmdLowCoverageBase(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&coverprofile, "./testdata/low_coverage.out"),
	).Install().Restore()

	rootCmd.Run(rootCmd, []string{})

	assert.Equal(t, "5 statements out of 19 covered; overall coverage: 26.3%\n", outStream.String())
	assert.Equal(t, "", errStream.String())
}

func TestRootCmdLowCoverageLowThreshold(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&threshold, 25.0),
		patcher.SetVar(&coverprofile, "./testdata/low_coverage.out"),
	).Install().Restore()

	rootCmd.Run(rootCmd, []string{})

	assert.Equal(t, "5 statements out of 19 covered; overall coverage: 26.3%\n", outStream.String())
	assert.Equal(t, "", errStream.String())
}

func TestRootCmdLowCoverageHighThreshold(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&threshold, 50.0),
		patcher.SetVar(&coverprofile, "./testdata/low_coverage.out"),
	).Install().Restore()

	assert.PanicsWithValue(t, "os.Exit(1)", func() { rootCmd.Run(rootCmd, []string{}) })
	assert.Equal(t, "5 statements out of 19 covered; overall coverage: 26.3%\n", outStream.String())
	assert.Equal(t, "\nFailed to meet coverage threshold of 50.0%\n", errStream.String())
}

func TestExecuteSuccess(t *testing.T) {
	errStream := &bytes.Buffer{}
	defer patcher.NewPatchMaster(
		patcher.SetVar(&rootCmd.Run, func(cmd *cobra.Command, args []string) {}),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
	).Install().Restore()

	Execute()

	assert.Equal(t, "", errStream.String())
}

func TestExecuteFailure(t *testing.T) {
	errStream := &bytes.Buffer{}
	defer patcher.NewPatchMaster(
		patcher.SetVar(&rootCmd.RunE, func(cmd *cobra.Command, args []string) error {
			return assert.AnError
		}),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
	).Install().Restore()

	assert.PanicsWithValue(t, "os.Exit(4)", Execute)
	assert.Equal(t, fmt.Sprintf("%s\n", assert.AnError), errStream.String())
}
