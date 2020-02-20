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

	"github.com/klmitch/overcover/common"
)

func TestRootCmdBase(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    0.0,
		"min_headroom": 0.0,
		"max_headroom": 0.0,
	}
	setConfigCalled := false
	writeConfigCalled := false
	loadCoverageCalled := false
	loadStatementsCalled := false
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			setConfigCalled = true
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			writeConfigCalled = true
			return nil
		}),
		patcher.SetVar(&loadCoverage, func(filename string) (common.DataSet, error) {
			assert.Equal(t, "coverage.out", filename)
			loadCoverageCalled = true
			return common.DataSet{
				common.FileData{
					Package: "some/package",
					Name:    "file1.go",
					Count:   10,
					Exec:    10,
				},
				common.FileData{
					Package: "some/package",
					Name:    "file2.go",
					Count:   5,
					Exec:    5,
				},
				common.FileData{
					Package: "other/package",
					Name:    "file3.go",
					Count:   4,
					Exec:    4,
				},
			}, nil
		}),
		patcher.SetVar(&loadStatements, func(ba, args []string) (common.DataSet, error) {
			assert.Equal(t, []string{}, ba)
			assert.Equal(t, []string{}, args)
			loadStatementsCalled = true
			return common.DataSet{}, nil
		}),
	).Install().Restore()

	rootCmd.Run(rootCmd, []string{})

	assert.Equal(t, "19 statements out of 19 covered; overall coverage: 100.0%\n", outStream.String())
	assert.Equal(t, "", errStream.String())
	assert.False(t, setConfigCalled)
	assert.False(t, writeConfigCalled)
	assert.True(t, loadCoverageCalled)
	assert.False(t, loadStatementsCalled)
}

func TestRootCmdStatementsOnly(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    0.0,
		"min_headroom": 0.0,
		"max_headroom": 0.0,
	}
	setConfigCalled := false
	writeConfigCalled := false
	loadCoverageCalled := false
	loadStatementsCalled := false
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
			setConfigCalled = true
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			writeConfigCalled = true
			return nil
		}),
		patcher.SetVar(&loadCoverage, func(filename string) (common.DataSet, error) {
			assert.Equal(t, "coverage.out", filename)
			loadCoverageCalled = true
			return common.DataSet{
				common.FileData{
					Package: "some/package",
					Name:    "file1.go",
					Count:   10,
					Exec:    10,
				},
				common.FileData{
					Package: "some/package",
					Name:    "file2.go",
					Count:   5,
					Exec:    5,
				},
				common.FileData{
					Package: "other/package",
					Name:    "file3.go",
					Count:   4,
					Exec:    4,
				},
			}, nil
		}),
		patcher.SetVar(&loadStatements, func(ba, args []string) (common.DataSet, error) {
			assert.Equal(t, []string{}, ba)
			assert.Equal(t, []string{"./..."}, args)
			loadStatementsCalled = true
			return common.DataSet{
				common.FileData{
					Package: "some/package",
					Name:    "file1.go",
					Count:   10,
				},
				common.FileData{
					Package: "some/package",
					Name:    "file2.go",
					Count:   5,
				},
				common.FileData{
					Package: "other/package",
					Name:    "file3.go",
					Count:   4,
				},
			}, nil
		}),
	).Install().Restore()

	rootCmd.Run(rootCmd, []string{"./..."})

	assert.Equal(t, "0 statements out of 19 covered; overall coverage: 0.0%\n", outStream.String())
	assert.Equal(t, "", errStream.String())
	assert.False(t, setConfigCalled)
	assert.False(t, writeConfigCalled)
	assert.False(t, loadCoverageCalled)
	assert.True(t, loadStatementsCalled)
}

func TestRootCmdStatements(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    0.0,
		"min_headroom": 0.0,
		"max_headroom": 0.0,
	}
	setConfigCalled := false
	writeConfigCalled := false
	loadCoverageCalled := false
	loadStatementsCalled := false
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			setConfigCalled = true
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			writeConfigCalled = true
			return nil
		}),
		patcher.SetVar(&loadCoverage, func(filename string) (common.DataSet, error) {
			assert.Equal(t, "coverage.out", filename)
			loadCoverageCalled = true
			return common.DataSet{
				common.FileData{
					Package: "some/package",
					Name:    "file1.go",
					Count:   10,
					Exec:    10,
				},
				common.FileData{
					Package: "some/package",
					Name:    "file2.go",
					Count:   5,
					Exec:    5,
				},
				common.FileData{
					Package: "other/package",
					Name:    "file3.go",
					Count:   4,
					Exec:    4,
				},
			}, nil
		}),
		patcher.SetVar(&loadStatements, func(ba, args []string) (common.DataSet, error) {
			assert.Equal(t, []string{}, ba)
			assert.Equal(t, []string{"./..."}, args)
			loadStatementsCalled = true
			return common.DataSet{
				common.FileData{
					Package: "some/package",
					Name:    "file1.go",
					Count:   10,
				},
				common.FileData{
					Package: "some/package",
					Name:    "file2.go",
					Count:   5,
				},
				common.FileData{
					Package: "other/package",
					Name:    "file3.go",
					Count:   4,
				},
			}, nil
		}),
	).Install().Restore()

	rootCmd.Run(rootCmd, []string{"./..."})

	assert.Equal(t, "19 statements out of 19 covered; overall coverage: 100.0%\n", outStream.String())
	assert.Equal(t, "", errStream.String())
	assert.False(t, setConfigCalled)
	assert.False(t, writeConfigCalled)
	assert.True(t, loadCoverageCalled)
	assert.True(t, loadStatementsCalled)
}

func TestRootCmdStatementsBuildArgs(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    0.0,
		"min_headroom": 0.0,
		"max_headroom": 0.0,
	}
	setConfigCalled := false
	writeConfigCalled := false
	loadCoverageCalled := false
	loadStatementsCalled := false
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			setConfigCalled = true
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			writeConfigCalled = true
			return nil
		}),
		patcher.SetVar(&loadCoverage, func(filename string) (common.DataSet, error) {
			assert.Equal(t, "coverage.out", filename)
			loadCoverageCalled = true
			return common.DataSet{
				common.FileData{
					Package: "some/package",
					Name:    "file1.go",
					Count:   10,
					Exec:    10,
				},
				common.FileData{
					Package: "some/package",
					Name:    "file2.go",
					Count:   5,
					Exec:    5,
				},
				common.FileData{
					Package: "other/package",
					Name:    "file3.go",
					Count:   4,
					Exec:    4,
				},
			}, nil
		}),
		patcher.SetVar(&loadStatements, func(ba, args []string) (common.DataSet, error) {
			assert.Equal(t, []string{"arg1", "arg2", "arg3"}, ba)
			assert.Equal(t, []string{"./..."}, args)
			loadStatementsCalled = true
			return common.DataSet{
				common.FileData{
					Package: "some/package",
					Name:    "file1.go",
					Count:   10,
				},
				common.FileData{
					Package: "some/package",
					Name:    "file2.go",
					Count:   5,
				},
				common.FileData{
					Package: "other/package",
					Name:    "file3.go",
					Count:   4,
				},
			}, nil
		}),
		patcher.SetVar(&buildArgs, []string{"arg1", "arg2", "arg3"}),
	).Install().Restore()

	rootCmd.Run(rootCmd, []string{"./..."})

	assert.Equal(t, "19 statements out of 19 covered; overall coverage: 100.0%\n", outStream.String())
	assert.Equal(t, "", errStream.String())
	assert.False(t, setConfigCalled)
	assert.False(t, writeConfigCalled)
	assert.True(t, loadCoverageCalled)
	assert.True(t, loadStatementsCalled)
}

func TestRootCmdStatementsExtra(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    0.0,
		"min_headroom": 0.0,
		"max_headroom": 0.0,
	}
	setConfigCalled := false
	writeConfigCalled := false
	loadCoverageCalled := false
	loadStatementsCalled := false
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			setConfigCalled = true
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			writeConfigCalled = true
			return nil
		}),
		patcher.SetVar(&loadCoverage, func(filename string) (common.DataSet, error) {
			assert.Equal(t, "coverage.out", filename)
			loadCoverageCalled = true
			return common.DataSet{
				common.FileData{
					Package: "some/package",
					Name:    "file1.go",
					Count:   10,
					Exec:    10,
				},
				common.FileData{
					Package: "some/package",
					Name:    "file2.go",
					Count:   5,
					Exec:    5,
				},
				common.FileData{
					Package: "other/package",
					Name:    "file3.go",
					Count:   4,
					Exec:    4,
				},
			}, nil
		}),
		patcher.SetVar(&loadStatements, func(ba, args []string) (common.DataSet, error) {
			assert.Equal(t, []string{}, ba)
			assert.Equal(t, []string{"./..."}, args)
			loadStatementsCalled = true
			return common.DataSet{
				common.FileData{
					Package: "some/package",
					Name:    "file4.go",
					Count:   10,
				},
				common.FileData{
					Package: "some/package",
					Name:    "file5.go",
					Count:   5,
				},
				common.FileData{
					Package: "other/package",
					Name:    "file6.go",
					Count:   4,
				},
			}, nil
		}),
	).Install().Restore()

	rootCmd.Run(rootCmd, []string{"./..."})

	assert.Equal(t, "19 statements out of 38 covered; overall coverage: 50.0%\n", outStream.String())
	assert.Equal(t, "", errStream.String())
	assert.False(t, setConfigCalled)
	assert.False(t, writeConfigCalled)
	assert.True(t, loadCoverageCalled)
	assert.True(t, loadStatementsCalled)
}

func TestRootCmdStatementsConflict(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    0.0,
		"min_headroom": 0.0,
		"max_headroom": 0.0,
	}
	setConfigCalled := false
	writeConfigCalled := false
	loadCoverageCalled := false
	loadStatementsCalled := false
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			setConfigCalled = true
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			writeConfigCalled = true
			return nil
		}),
		patcher.SetVar(&loadCoverage, func(filename string) (common.DataSet, error) {
			assert.Equal(t, "coverage.out", filename)
			loadCoverageCalled = true
			return common.DataSet{
				common.FileData{
					Package: "some/package",
					Name:    "file1.go",
					Count:   10,
					Exec:    10,
				},
				common.FileData{
					Package: "some/package",
					Name:    "file2.go",
					Count:   5,
					Exec:    5,
				},
				common.FileData{
					Package: "other/package",
					Name:    "file3.go",
					Count:   4,
					Exec:    4,
				},
			}, nil
		}),
		patcher.SetVar(&loadStatements, func(ba, args []string) (common.DataSet, error) {
			assert.Equal(t, []string{}, ba)
			assert.Equal(t, []string{"./..."}, args)
			loadStatementsCalled = true
			return common.DataSet{
				common.FileData{
					Package: "some/package",
					Name:    "file1.go",
					Count:   10,
				},
				common.FileData{
					Package: "some/package",
					Name:    "file2.go",
					Count:   5,
				},
				common.FileData{
					Package: "other/package",
					Name:    "file3.go",
					Count:   5,
				},
			}, nil
		}),
	).Install().Restore()

	rootCmd.Run(rootCmd, []string{"./..."})

	assert.Equal(t, "19 statements out of 19 covered; overall coverage: 100.0%\n", outStream.String())
	assert.Equal(t, "WARNING: coverage profile coverage.out may not match source; potentially altered files:\n  other/package/file3.go\n", errStream.String())
	assert.False(t, setConfigCalled)
	assert.False(t, writeConfigCalled)
	assert.True(t, loadCoverageCalled)
	assert.True(t, loadStatementsCalled)
}

func TestRootCmdNoProfile(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    0.0,
		"min_headroom": 0.0,
		"max_headroom": 0.0,
	}
	setConfigCalled := false
	writeConfigCalled := false
	loadCoverageCalled := false
	loadStatementsCalled := false
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
			setConfigCalled = true
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			writeConfigCalled = true
			return nil
		}),
		patcher.SetVar(&loadCoverage, func(filename string) (common.DataSet, error) {
			assert.Equal(t, "coverage.out", filename)
			loadCoverageCalled = true
			return common.DataSet{
				common.FileData{
					Package: "some/package",
					Name:    "file1.go",
					Count:   10,
					Exec:    10,
				},
				common.FileData{
					Package: "some/package",
					Name:    "file2.go",
					Count:   5,
					Exec:    5,
				},
				common.FileData{
					Package: "other/package",
					Name:    "file3.go",
					Count:   4,
					Exec:    4,
				},
			}, nil
		}),
		patcher.SetVar(&loadStatements, func(ba, args []string) (common.DataSet, error) {
			assert.Equal(t, []string{}, ba)
			assert.Equal(t, []string{}, args)
			loadStatementsCalled = true
			return common.DataSet{}, nil
		}),
	).Install().Restore()

	assert.PanicsWithValue(t, "os.Exit(2)", func() { rootCmd.Run(rootCmd, []string{}) })
	assert.Equal(t, "", outStream.String())
	assert.Equal(t, "No coverage profile file specified!  Use -p or provide a configuration file.\n", errStream.String())
	assert.False(t, setConfigCalled)
	assert.False(t, writeConfigCalled)
	assert.False(t, loadCoverageCalled)
	assert.False(t, loadStatementsCalled)
}

func TestRootCmdLoadCoverageFails(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    0.0,
		"min_headroom": 0.0,
		"max_headroom": 0.0,
	}
	setConfigCalled := false
	writeConfigCalled := false
	loadCoverageCalled := false
	loadStatementsCalled := false
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			setConfigCalled = true
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			writeConfigCalled = true
			return nil
		}),
		patcher.SetVar(&loadCoverage, func(filename string) (common.DataSet, error) {
			assert.Equal(t, "coverage.out", filename)
			loadCoverageCalled = true
			return nil, assert.AnError
		}),
		patcher.SetVar(&loadStatements, func(ba, args []string) (common.DataSet, error) {
			assert.Equal(t, []string{}, ba)
			assert.Equal(t, []string{}, args)
			loadStatementsCalled = true
			return common.DataSet{}, nil
		}),
	).Install().Restore()

	assert.PanicsWithValue(t, "os.Exit(3)", func() { rootCmd.Run(rootCmd, []string{}) })
	assert.Equal(t, "", outStream.String())
	assert.Equal(t, fmt.Sprintf("Unable to read coverage profile file \"coverage.out\": %s\n", assert.AnError), errStream.String())
	assert.False(t, setConfigCalled)
	assert.False(t, writeConfigCalled)
	assert.True(t, loadCoverageCalled)
	assert.False(t, loadStatementsCalled)
}

func TestRootCmdLoadStatementsFails(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    0.0,
		"min_headroom": 0.0,
		"max_headroom": 0.0,
	}
	setConfigCalled := false
	writeConfigCalled := false
	loadCoverageCalled := false
	loadStatementsCalled := false
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			setConfigCalled = true
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			writeConfigCalled = true
			return nil
		}),
		patcher.SetVar(&loadCoverage, func(filename string) (common.DataSet, error) {
			assert.Equal(t, "coverage.out", filename)
			loadCoverageCalled = true
			return common.DataSet{
				common.FileData{
					Package: "some/package",
					Name:    "file1.go",
					Count:   10,
					Exec:    10,
				},
				common.FileData{
					Package: "some/package",
					Name:    "file2.go",
					Count:   5,
					Exec:    5,
				},
				common.FileData{
					Package: "other/package",
					Name:    "file3.go",
					Count:   4,
					Exec:    4,
				},
			}, nil
		}),
		patcher.SetVar(&loadStatements, func(ba, args []string) (common.DataSet, error) {
			assert.Equal(t, []string{}, ba)
			assert.Equal(t, []string{"./..."}, args)
			loadStatementsCalled = true
			return nil, assert.AnError
		}),
	).Install().Restore()

	assert.PanicsWithValue(t, "os.Exit(3)", func() { rootCmd.Run(rootCmd, []string{"./..."}) })
	assert.Equal(t, "", outStream.String())
	assert.Equal(t, fmt.Sprintf("Unable to read source: %s\n", assert.AnError), errStream.String())
	assert.False(t, setConfigCalled)
	assert.False(t, writeConfigCalled)
	assert.True(t, loadCoverageCalled)
	assert.True(t, loadStatementsCalled)
}

func TestRootCmdDetailed(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    0.0,
		"min_headroom": 0.0,
		"max_headroom": 0.0,
	}
	setConfigCalled := false
	writeConfigCalled := false
	loadCoverageCalled := false
	loadStatementsCalled := false
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			setConfigCalled = true
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			writeConfigCalled = true
			return nil
		}),
		patcher.SetVar(&loadCoverage, func(filename string) (common.DataSet, error) {
			assert.Equal(t, "coverage.out", filename)
			loadCoverageCalled = true
			return common.DataSet{
				common.FileData{
					Package: "some/package",
					Name:    "file1.go",
					Count:   10,
					Exec:    10,
				},
				common.FileData{
					Package: "some/package",
					Name:    "file2.go",
					Count:   5,
					Exec:    5,
				},
				common.FileData{
					Package: "other/package",
					Name:    "file3.go",
					Count:   4,
					Exec:    4,
				},
			}, nil
		}),
		patcher.SetVar(&loadStatements, func(ba, args []string) (common.DataSet, error) {
			assert.Equal(t, []string{}, ba)
			assert.Equal(t, []string{}, args)
			loadStatementsCalled = true
			return common.DataSet{}, nil
		}),
		patcher.SetVar(&detailed, true),
	).Install().Restore()

	rootCmd.Run(rootCmd, []string{})

	assert.Equal(t, "Per file details:\n"+
		" File                    Executed  Total  Coverage\n"+
		" ----                    --------  -----  --------\n"+
		" other/package/file3.go  4         4      100.0%\n"+
		" some/package/file1.go   10        10     100.0%\n"+
		" some/package/file2.go   5         5      100.0%\n"+
		"\n"+
		"19 statements out of 19 covered; overall coverage: 100.0%\n",
		outStream.String(),
	)
	assert.Equal(t, "", errStream.String())
	assert.False(t, setConfigCalled)
	assert.False(t, writeConfigCalled)
	assert.True(t, loadCoverageCalled)
	assert.False(t, loadStatementsCalled)
}

func TestRootCmdSummary(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    0.0,
		"min_headroom": 0.0,
		"max_headroom": 0.0,
	}
	setConfigCalled := false
	writeConfigCalled := false
	loadCoverageCalled := false
	loadStatementsCalled := false
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			setConfigCalled = true
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			writeConfigCalled = true
			return nil
		}),
		patcher.SetVar(&loadCoverage, func(filename string) (common.DataSet, error) {
			assert.Equal(t, "coverage.out", filename)
			loadCoverageCalled = true
			return common.DataSet{
				common.FileData{
					Package: "some/package",
					Name:    "file1.go",
					Count:   10,
					Exec:    10,
				},
				common.FileData{
					Package: "some/package",
					Name:    "file2.go",
					Count:   5,
					Exec:    5,
				},
				common.FileData{
					Package: "other/package",
					Name:    "file3.go",
					Count:   4,
					Exec:    4,
				},
			}, nil
		}),
		patcher.SetVar(&loadStatements, func(ba, args []string) (common.DataSet, error) {
			assert.Equal(t, []string{}, ba)
			assert.Equal(t, []string{}, args)
			loadStatementsCalled = true
			return common.DataSet{}, nil
		}),
		patcher.SetVar(&summary, true),
	).Install().Restore()

	rootCmd.Run(rootCmd, []string{})

	assert.Equal(t, "Per package summary:\n"+
		" Package         Executed  Total  Coverage\n"+
		" -------         --------  -----  --------\n"+
		" other/package/  4         4      100.0%\n"+
		" some/package/   15        15     100.0%\n"+
		"\n"+
		"19 statements out of 19 covered; overall coverage: 100.0%\n",
		outStream.String(),
	)
	assert.Equal(t, "", errStream.String())
	assert.False(t, setConfigCalled)
	assert.False(t, writeConfigCalled)
	assert.True(t, loadCoverageCalled)
	assert.False(t, loadStatementsCalled)
}

func TestRootCmdLowCoverageBase(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    0.0,
		"min_headroom": 0.0,
		"max_headroom": 0.0,
	}
	setConfigCalled := false
	writeConfigCalled := false
	loadCoverageCalled := false
	loadStatementsCalled := false
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			setConfigCalled = true
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			writeConfigCalled = true
			return nil
		}),
		patcher.SetVar(&loadCoverage, func(filename string) (common.DataSet, error) {
			assert.Equal(t, "coverage.out", filename)
			loadCoverageCalled = true
			return common.DataSet{
				common.FileData{
					Package: "some/package",
					Name:    "file1.go",
					Count:   10,
				},
				common.FileData{
					Package: "some/package",
					Name:    "file2.go",
					Count:   5,
					Exec:    5,
				},
				common.FileData{
					Package: "other/package",
					Name:    "file3.go",
					Count:   4,
				},
			}, nil
		}),
		patcher.SetVar(&loadStatements, func(ba, args []string) (common.DataSet, error) {
			assert.Equal(t, []string{}, ba)
			assert.Equal(t, []string{}, args)
			loadStatementsCalled = true
			return common.DataSet{}, nil
		}),
	).Install().Restore()

	rootCmd.Run(rootCmd, []string{})

	assert.Equal(t, "5 statements out of 19 covered; overall coverage: 26.3%\n", outStream.String())
	assert.Equal(t, "", errStream.String())
	assert.False(t, setConfigCalled)
	assert.False(t, writeConfigCalled)
	assert.True(t, loadCoverageCalled)
	assert.False(t, loadStatementsCalled)
}

func TestRootCmdLowCoverageLowThreshold(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    25.0,
		"min_headroom": 0.0,
		"max_headroom": 0.0,
	}
	setConfigCalled := false
	writeConfigCalled := false
	loadCoverageCalled := false
	loadStatementsCalled := false
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			setConfigCalled = true
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			writeConfigCalled = true
			return nil
		}),
		patcher.SetVar(&loadCoverage, func(filename string) (common.DataSet, error) {
			assert.Equal(t, "coverage.out", filename)
			loadCoverageCalled = true
			return common.DataSet{
				common.FileData{
					Package: "some/package",
					Name:    "file1.go",
					Count:   10,
				},
				common.FileData{
					Package: "some/package",
					Name:    "file2.go",
					Count:   5,
					Exec:    5,
				},
				common.FileData{
					Package: "other/package",
					Name:    "file3.go",
					Count:   4,
				},
			}, nil
		}),
		patcher.SetVar(&loadStatements, func(ba, args []string) (common.DataSet, error) {
			assert.Equal(t, []string{}, ba)
			assert.Equal(t, []string{}, args)
			loadStatementsCalled = true
			return common.DataSet{}, nil
		}),
	).Install().Restore()

	rootCmd.Run(rootCmd, []string{})

	assert.Equal(t, "5 statements out of 19 covered; overall coverage: 26.3%\n", outStream.String())
	assert.Equal(t, "", errStream.String())
	assert.False(t, setConfigCalled)
	assert.False(t, writeConfigCalled)
	assert.True(t, loadCoverageCalled)
	assert.False(t, loadStatementsCalled)
}

func TestRootCmdLowCoverageHighThreshold(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    50.0,
		"min_headroom": 0.0,
		"max_headroom": 0.0,
	}
	setConfigCalled := false
	writeConfigCalled := false
	loadCoverageCalled := false
	loadStatementsCalled := false
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			setConfigCalled = true
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			writeConfigCalled = true
			return nil
		}),
		patcher.SetVar(&loadCoverage, func(filename string) (common.DataSet, error) {
			assert.Equal(t, "coverage.out", filename)
			loadCoverageCalled = true
			return common.DataSet{
				common.FileData{
					Package: "some/package",
					Name:    "file1.go",
					Count:   10,
				},
				common.FileData{
					Package: "some/package",
					Name:    "file2.go",
					Count:   5,
					Exec:    5,
				},
				common.FileData{
					Package: "other/package",
					Name:    "file3.go",
					Count:   4,
				},
			}, nil
		}),
		patcher.SetVar(&loadStatements, func(ba, args []string) (common.DataSet, error) {
			assert.Equal(t, []string{}, ba)
			assert.Equal(t, []string{}, args)
			loadStatementsCalled = true
			return common.DataSet{}, nil
		}),
	).Install().Restore()

	assert.PanicsWithValue(t, "os.Exit(1)", func() { rootCmd.Run(rootCmd, []string{}) })
	assert.Equal(t, "5 statements out of 19 covered; overall coverage: 26.3%\n", outStream.String())
	assert.Equal(t, "\nFailed to meet coverage threshold of 50.0%\n", errStream.String())
	assert.False(t, setConfigCalled)
	assert.False(t, writeConfigCalled)
	assert.True(t, loadCoverageCalled)
	assert.False(t, loadStatementsCalled)
}

func TestRootCmdUpdateNeededNoConfig(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    75.0,
		"min_headroom": 1.0,
		"max_headroom": 2.0,
	}
	setConfigCalled := false
	writeConfigCalled := false
	loadCoverageCalled := false
	loadStatementsCalled := false
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			setConfigCalled = true
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			writeConfigCalled = true
			return nil
		}),
		patcher.SetVar(&loadCoverage, func(filename string) (common.DataSet, error) {
			assert.Equal(t, "coverage.out", filename)
			loadCoverageCalled = true
			return common.DataSet{
				common.FileData{
					Package: "some/package",
					Name:    "file1.go",
					Count:   10,
					Exec:    10,
				},
				common.FileData{
					Package: "some/package",
					Name:    "file2.go",
					Count:   5,
					Exec:    5,
				},
				common.FileData{
					Package: "other/package",
					Name:    "file3.go",
					Count:   4,
					Exec:    4,
				},
			}, nil
		}),
		patcher.SetVar(&loadStatements, func(ba, args []string) (common.DataSet, error) {
			assert.Equal(t, []string{}, ba)
			assert.Equal(t, []string{}, args)
			loadStatementsCalled = true
			return common.DataSet{}, nil
		}),
	).Install().Restore()

	rootCmd.Run(rootCmd, []string{})

	assert.Equal(t, "19 statements out of 19 covered; overall coverage: 100.0%\n", outStream.String())
	assert.Equal(t, "", errStream.String())
	assert.False(t, setConfigCalled)
	assert.False(t, writeConfigCalled)
	assert.True(t, loadCoverageCalled)
	assert.False(t, loadStatementsCalled)
}

func TestRootCmdUpdateNeededWithConfigReadOnly(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    75.0,
		"min_headroom": 1.0,
		"max_headroom": 2.0,
	}
	setConfigCalled := false
	writeConfigCalled := false
	loadCoverageCalled := false
	loadStatementsCalled := false
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
			return "coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			setConfigCalled = true
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			writeConfigCalled = true
			return nil
		}),
		patcher.SetVar(&loadCoverage, func(filename string) (common.DataSet, error) {
			assert.Equal(t, "coverage.out", filename)
			loadCoverageCalled = true
			return common.DataSet{
				common.FileData{
					Package: "some/package",
					Name:    "file1.go",
					Count:   10,
					Exec:    10,
				},
				common.FileData{
					Package: "some/package",
					Name:    "file2.go",
					Count:   5,
					Exec:    5,
				},
				common.FileData{
					Package: "other/package",
					Name:    "file3.go",
					Count:   4,
					Exec:    4,
				},
			}, nil
		}),
		patcher.SetVar(&loadStatements, func(ba, args []string) (common.DataSet, error) {
			assert.Equal(t, []string{}, ba)
			assert.Equal(t, []string{}, args)
			loadStatementsCalled = true
			return common.DataSet{}, nil
		}),
	).Install().Restore()

	assert.PanicsWithValue(t, "os.Exit(5)", func() { rootCmd.Run(rootCmd, []string{}) })
	assert.Equal(t, "19 statements out of 19 covered; overall coverage: 100.0%\n", outStream.String())
	assert.Equal(t, "\nCoverage exceeds maximum headroom.  Update threshold to 99.0%\n", errStream.String())
	assert.False(t, setConfigCalled)
	assert.False(t, writeConfigCalled)
	assert.True(t, loadCoverageCalled)
	assert.False(t, loadStatementsCalled)
}

func TestRootCmdUpdateUnneededWithConfigReadOnly(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    99.5,
		"min_headroom": 1.0,
		"max_headroom": 2.0,
	}
	setConfigCalled := false
	writeConfigCalled := false
	loadCoverageCalled := false
	loadStatementsCalled := false
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
			return "coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			setConfigCalled = true
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			writeConfigCalled = true
			return nil
		}),
		patcher.SetVar(&loadCoverage, func(filename string) (common.DataSet, error) {
			assert.Equal(t, "coverage.out", filename)
			loadCoverageCalled = true
			return common.DataSet{
				common.FileData{
					Package: "some/package",
					Name:    "file1.go",
					Count:   10,
					Exec:    10,
				},
				common.FileData{
					Package: "some/package",
					Name:    "file2.go",
					Count:   5,
					Exec:    5,
				},
				common.FileData{
					Package: "other/package",
					Name:    "file3.go",
					Count:   4,
					Exec:    4,
				},
			}, nil
		}),
		patcher.SetVar(&loadStatements, func(ba, args []string) (common.DataSet, error) {
			assert.Equal(t, []string{}, ba)
			assert.Equal(t, []string{}, args)
			loadStatementsCalled = true
			return common.DataSet{}, nil
		}),
	).Install().Restore()

	rootCmd.Run(rootCmd, []string{})

	assert.Equal(t, "19 statements out of 19 covered; overall coverage: 100.0%\n", outStream.String())
	assert.Equal(t, "", errStream.String())
	assert.False(t, setConfigCalled)
	assert.False(t, writeConfigCalled)
	assert.True(t, loadCoverageCalled)
	assert.False(t, loadStatementsCalled)
}

func TestRootCmdUpdateNeededWithConfig(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    75.0,
		"min_headroom": 1.0,
		"max_headroom": 2.0,
	}
	setConfigCalled := false
	writeConfigCalled := false
	loadCoverageCalled := false
	loadStatementsCalled := false
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&config, "test.yaml"),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			assert.Equal(t, "threshold", name)
			assert.Equal(t, 99.0, value)
			setConfigCalled = true
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			assert.Equal(t, "test.yaml", fname)
			writeConfigCalled = true
			return nil
		}),
		patcher.SetVar(&loadCoverage, func(filename string) (common.DataSet, error) {
			assert.Equal(t, "coverage.out", filename)
			loadCoverageCalled = true
			return common.DataSet{
				common.FileData{
					Package: "some/package",
					Name:    "file1.go",
					Count:   10,
					Exec:    10,
				},
				common.FileData{
					Package: "some/package",
					Name:    "file2.go",
					Count:   5,
					Exec:    5,
				},
				common.FileData{
					Package: "other/package",
					Name:    "file3.go",
					Count:   4,
					Exec:    4,
				},
			}, nil
		}),
		patcher.SetVar(&loadStatements, func(ba, args []string) (common.DataSet, error) {
			assert.Equal(t, []string{}, ba)
			assert.Equal(t, []string{}, args)
			loadStatementsCalled = true
			return common.DataSet{}, nil
		}),
	).Install().Restore()

	rootCmd.Run(rootCmd, []string{})

	assert.Equal(t, "19 statements out of 19 covered; overall coverage: 100.0%\nUpdating configuration file test.yaml with new threshold value 99.0%\n", outStream.String())
	assert.Equal(t, "", errStream.String())
	assert.True(t, setConfigCalled)
	assert.True(t, writeConfigCalled)
	assert.True(t, loadCoverageCalled)
	assert.False(t, loadStatementsCalled)
}

func TestRootCmdUpdateNeededWithConfigFails(t *testing.T) {
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}
	values := map[string]float64{
		"threshold":    75.0,
		"min_headroom": 1.0,
		"max_headroom": 2.0,
	}
	setConfigCalled := false
	writeConfigCalled := false
	loadCoverageCalled := false
	loadStatementsCalled := false
	defer patcher.NewPatchMaster(
		patcher.SetVar(&stdout, outStream),
		patcher.SetVar(&stderr, errStream),
		patcher.SetVar(&config, "test.yaml"),
		patcher.SetVar(&exit, func(code int) {
			panic(fmt.Sprintf("os.Exit(%d)", code))
		}),
		patcher.SetVar(&getString, func(name string) string {
			assert.Equal(t, "coverprofile", name)
			return "coverage.out"
		}),
		patcher.SetVar(&getFloat64, func(name string) float64 {
			value, ok := values[name]
			assert.True(t, ok)
			return value
		}),
		patcher.SetVar(&setConfig, func(name string, value interface{}) {
			assert.Equal(t, "threshold", name)
			assert.Equal(t, 99.0, value)
			setConfigCalled = true
		}),
		patcher.SetVar(&writeConfig, func(fname string) error {
			assert.Equal(t, "test.yaml", fname)
			writeConfigCalled = true
			return assert.AnError
		}),
		patcher.SetVar(&loadCoverage, func(filename string) (common.DataSet, error) {
			assert.Equal(t, "coverage.out", filename)
			loadCoverageCalled = true
			return common.DataSet{
				common.FileData{
					Package: "some/package",
					Name:    "file1.go",
					Count:   10,
					Exec:    10,
				},
				common.FileData{
					Package: "some/package",
					Name:    "file2.go",
					Count:   5,
					Exec:    5,
				},
				common.FileData{
					Package: "other/package",
					Name:    "file3.go",
					Count:   4,
					Exec:    4,
				},
			}, nil
		}),
		patcher.SetVar(&loadStatements, func(ba, args []string) (common.DataSet, error) {
			assert.Equal(t, []string{}, ba)
			assert.Equal(t, []string{}, args)
			loadStatementsCalled = true
			return common.DataSet{}, nil
		}),
	).Install().Restore()

	assert.PanicsWithValue(t, "os.Exit(5)", func() { rootCmd.Run(rootCmd, []string{}) })
	assert.Equal(t, "19 statements out of 19 covered; overall coverage: 100.0%\nUpdating configuration file test.yaml with new threshold value 99.0%\n", outStream.String())
	assert.Contains(t, errStream.String(), "\nFailed to write updated config with new threshold 99.0% to test.yaml: ")
	assert.True(t, setConfigCalled)
	assert.True(t, writeConfigCalled)
	assert.True(t, loadCoverageCalled)
	assert.False(t, loadStatementsCalled)
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

func TestGetBuildArgDefaultUnset(t *testing.T) {
	defer patcher.UnsetEnv("OVERCOVER_BUILD_ARG").Install().Restore()

	result := getBuildArgDefault()

	assert.Equal(t, []string{}, result)
}

func TestGetBuildArgDefaultSet(t *testing.T) {
	defer patcher.SetEnv("OVERCOVER_BUILD_ARG", "this is a test").Install().Restore()

	result := getBuildArgDefault()

	assert.Equal(t, []string{"this", "is", "a", "test"}, result)
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
