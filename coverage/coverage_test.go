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

package coverage

import (
	"testing"

	"github.com/klmitch/patcher"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/cover"

	"github.com/klmitch/overcover/common"
)

func TestLoadBase(t *testing.T) {
	profs := []*cover.Profile{
		{
			FileName: "example.com/some/package/file1.go",
			Blocks: []cover.ProfileBlock{
				{NumStmt: 5, Count: 0},
				{NumStmt: 10, Count: 1},
				{NumStmt: 4, Count: 0},
				{NumStmt: 1, Count: 1},
			},
		},
		{
			FileName: "example.com/some/package/file2.go",
			Blocks: []cover.ProfileBlock{
				{NumStmt: 10, Count: 0},
				{NumStmt: 5, Count: 1},
				{NumStmt: 1, Count: 0},
				{NumStmt: 4, Count: 1},
			},
		},
	}
	parseProfilesCalled := false
	defer patcher.SetVar(&parseProfiles, func(profile string) ([]*cover.Profile, error) {
		assert.Equal(t, "coverage.out", profile)
		parseProfilesCalled = true
		return profs, nil
	}).Install().Restore()

	result, err := Load("coverage.out")

	assert.NoError(t, err)
	assert.Equal(t, common.DataSet{
		common.FileData{
			Package: "example.com/some/package",
			Name:    "file1.go",
			Count:   20,
			Exec:    11,
		},
		common.FileData{
			Package: "example.com/some/package",
			Name:    "file2.go",
			Count:   20,
			Exec:    9,
		},
	}, result)
	assert.True(t, parseProfilesCalled)
}

func TestLoadError(t *testing.T) {
	parseProfilesCalled := false
	defer patcher.SetVar(&parseProfiles, func(profile string) ([]*cover.Profile, error) {
		assert.Equal(t, "coverage.out", profile)
		parseProfilesCalled = true
		return nil, assert.AnError
	}).Install().Restore()

	result, err := Load("coverage.out")

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	assert.True(t, parseProfilesCalled)
}
