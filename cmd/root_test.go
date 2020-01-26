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
	values := map[string]float64{
		"threshold":    0.0,
		"min_headroom": 0.0,
		"max_headroom": 0.0,
	}
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "./testdata/full_coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			panic("should not be called")
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			panic("should not be called")
		}),
	).Install().Restore()

	rootCmd.Run(rootCmd, []string{})

	assert.Equal(t, "19 statements out of 19 covered; overall coverage: 100.0%\n", outStream.String())
	assert.Equal(t, "", errStream.String())
}

func TestRootCmdNoProfile(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    0.0,
		"min_headroom": 0.0,
		"max_headroom": 0.0,
	}
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return ""
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			panic("should not be called")
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			panic("should not be called")
		}),
	).Install().Restore()

	assert.PanicsWithValue(t, "os.Exit(2)", func() { rootCmd.Run(rootCmd, []string{}) })
	assert.Equal(t, "", outStream.String())
	assert.Equal(t, "No coverage profile file specified!  Use -p or provide a configuration file.\n", errStream.String())
}

func TestRootCmdOpenFails(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    0.0,
		"min_headroom": 0.0,
		"max_headroom": 0.0,
	}
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "./testdata/no_such_coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			panic("should not be called")
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			panic("should not be called")
		}),
	).Install().Restore()

	assert.PanicsWithValue(t, "os.Exit(2)", func() { rootCmd.Run(rootCmd, []string{}) })
	assert.Equal(t, "", outStream.String())
	assert.Contains(t, errStream.String(), "Unable to open coverage profile file: ")
}

func TestRootCmdLoadCoverageFails(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    0.0,
		"min_headroom": 0.0,
		"max_headroom": 0.0,
	}
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "./testdata/bad_coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			panic("should not be called")
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			panic("should not be called")
		}),
	).Install().Restore()

	assert.PanicsWithValue(t, "os.Exit(3)", func() { rootCmd.Run(rootCmd, []string{}) })
	assert.Equal(t, "", outStream.String())
	assert.Contains(t, errStream.String(), "Unable to read coverage profile file: ")
}

func TestRootCmdLowCoverageBase(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    0.0,
		"min_headroom": 0.0,
		"max_headroom": 0.0,
	}
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "./testdata/low_coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			panic("should not be called")
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			panic("should not be called")
		}),
	).Install().Restore()

	rootCmd.Run(rootCmd, []string{})

	assert.Equal(t, "5 statements out of 19 covered; overall coverage: 26.3%\n", outStream.String())
	assert.Equal(t, "", errStream.String())
}

func TestRootCmdLowCoverageLowThreshold(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    25.0,
		"min_headroom": 0.0,
		"max_headroom": 0.0,
	}
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "./testdata/low_coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			panic("should not be called")
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			panic("should not be called")
		}),
	).Install().Restore()

	rootCmd.Run(rootCmd, []string{})

	assert.Equal(t, "5 statements out of 19 covered; overall coverage: 26.3%\n", outStream.String())
	assert.Equal(t, "", errStream.String())
}

func TestRootCmdLowCoverageHighThreshold(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    50.0,
		"min_headroom": 0.0,
		"max_headroom": 0.0,
	}
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "./testdata/low_coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			panic("should not be called")
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			panic("should not be called")
		}),
	).Install().Restore()

	assert.PanicsWithValue(t, "os.Exit(1)", func() { rootCmd.Run(rootCmd, []string{}) })
	assert.Equal(t, "5 statements out of 19 covered; overall coverage: 26.3%\n", outStream.String())
	assert.Equal(t, "\nFailed to meet coverage threshold of 50.0%\n", errStream.String())
}

func TestRootCmdUpdateNeededNoConfig(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    75.0,
		"min_headroom": 1.0,
		"max_headroom": 2.0,
	}
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "./testdata/full_coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			panic("should not be called")
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			panic("should not be called")
		}),
	).Install().Restore()

	rootCmd.Run(rootCmd, []string{})

	assert.Equal(t, "19 statements out of 19 covered; overall coverage: 100.0%\n", outStream.String())
	assert.Equal(t, "", errStream.String())
}

func TestRootCmdUpdateNeededWithConfigReadOnly(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    75.0,
		"min_headroom": 1.0,
		"max_headroom": 2.0,
	}
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&config, "test.yaml"),
		patcher.SetVar(&readOnly, true),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "./testdata/full_coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			panic("should not be called")
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			panic("should not be called")
		}),
	).Install().Restore()

	assert.PanicsWithValue(t, "os.Exit(5)", func() { rootCmd.Run(rootCmd, []string{}) })
	assert.Equal(t, "19 statements out of 19 covered; overall coverage: 100.0%\n", outStream.String())
	assert.Equal(t, "\nCoverage exceeds maximum headroom.  Update threshold to 99.0%\n", errStream.String())
}

func TestRootCmdUpdateUnneededWithConfigReadOnly(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    99.5,
		"min_headroom": 1.0,
		"max_headroom": 2.0,
	}
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&config, "test.yaml"),
		patcher.SetVar(&readOnly, true),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "./testdata/full_coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			panic("should not be called")
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			panic("should not be called")
		}),
	).Install().Restore()

	rootCmd.Run(rootCmd, []string{})

	assert.Equal(t, "19 statements out of 19 covered; overall coverage: 100.0%\n", outStream.String())
	assert.Equal(t, "", errStream.String())
}

func TestRootCmdUpdateNeededWithConfig(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    75.0,
		"min_headroom": 1.0,
		"max_headroom": 2.0,
	}
	setCalled := false
	writeCalled := false
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&config, "test.yaml"),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "./testdata/full_coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			assert.Equal(t, "threshold", name)
			assert.Equal(t, 99.0, value)
			setCalled = true
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			assert.Equal(t, "test.yaml", fname)
			writeCalled = true
			return nil
		}),
	).Install().Restore()

	rootCmd.Run(rootCmd, []string{})

	assert.Equal(t, "19 statements out of 19 covered; overall coverage: 100.0%\nUpdating configuration file test.yaml with new threshold value 99.0%\n", outStream.String())
	assert.Equal(t, "", errStream.String())
	assert.True(t, setCalled)
	assert.True(t, writeCalled)
}

func TestRootCmdUpdateNeededWithConfigFails(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    75.0,
		"min_headroom": 1.0,
		"max_headroom": 2.0,
	}
	setCalled := false
	writeCalled := false
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&config, "test.yaml"),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "./testdata/full_coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			assert.Equal(t, "threshold", name)
			assert.Equal(t, 99.0, value)
			setCalled = true
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			assert.Equal(t, "test.yaml", fname)
			writeCalled = true
			return assert.AnError
		}),
	).Install().Restore()

	assert.PanicsWithValue(t, "os.Exit(5)", func() { rootCmd.Run(rootCmd, []string{}) })
	assert.Equal(t, "19 statements out of 19 covered; overall coverage: 100.0%\nUpdating configuration file test.yaml with new threshold value 99.0%\n", outStream.String())
	assert.Contains(t, errStream.String(), "\nFailed to write updated config with new threshold 99.0% to test.yaml: ")
	assert.True(t, setCalled)
	assert.True(t, writeCalled)
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

func TestReadConfigBase(t *testing.T) {
	outStream := &bytes.Buffer{}
	var setCalled, readCalled bool
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&config, "config"),
		patcher.SetVar(&setConfigFile, func(fname string) {
			assert.Equal(t, "config", fname)
			setCalled = true
		}),
		patcher.SetVar(&readInConfig, func() error {
			readCalled = true
			return nil
		}),
		patcher.SetVar(&configFileUsed, func() string {
			return "config.yaml"
		}),
	).Install().Restore()

	readConfig()

	assert.True(t, setCalled)
	assert.True(t, readCalled)
	assert.Equal(t, "Using configuration file config.yaml\n", outStream.String())
}

func TestReadConfigReadFails(t *testing.T) {
	outStream := &bytes.Buffer{}
	var setCalled, readCalled bool
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&config, "config"),
		patcher.SetVar(&setConfigFile, func(fname string) {
			assert.Equal(t, "config", fname)
			setCalled = true
		}),
		patcher.SetVar(&readInConfig, func() error {
			readCalled = true
			return assert.AnError
		}),
		patcher.SetVar(&configFileUsed, func() string {
			return "config.yaml"
		}),
	).Install().Restore()

	readConfig()

	assert.True(t, setCalled)
	assert.True(t, readCalled)
	assert.Equal(t, "", outStream.String())
}

func TestReadConfigNoConfig(t *testing.T) {
	outStream := &bytes.Buffer{}
	var setCalled, readCalled bool
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&setConfigFile, func(fname string) {
			assert.Equal(t, "config", fname)
			setCalled = true
		}),
		patcher.SetVar(&readInConfig, func() error {
			readCalled = true
			return nil
		}),
		patcher.SetVar(&configFileUsed, func() string {
			return "config.yaml"
		}),
	).Install().Restore()

	readConfig()

	assert.False(t, setCalled)
	assert.False(t, readCalled)
	assert.Equal(t, "", outStream.String())
}