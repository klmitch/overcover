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
	"path"

	"golang.org/x/tools/cover"

	"github.com/klmitch/overcover/common"
)

// Patch points for top-level functions called by functions in this
// file.
var (
	parseProfiles func(string) ([]*cover.Profile, error) = cover.ParseProfiles
)

// Load loads a coverage profile file and returns a list of FileData
// instances.
func Load(profile string) (common.DataSet, error) {
	// Begin by loading the profile file
	profs, err := parseProfiles(profile)
	if err != nil {
		return nil, err
	}

	// Next, process each profile
	var data []common.FileData
	for _, prof := range profs {
		fd := common.FileData{
			Package: path.Dir(prof.FileName),
			Name:    path.Base(prof.FileName),
		}

		// Process each block
		for _, blk := range prof.Blocks {
			fd.Count += int64(blk.NumStmt)

			// Has it been executed?
			if blk.Count > 0 {
				fd.Exec += int64(blk.NumStmt)
			}
		}

		// Append it to our results
		data = append(data, fd)
	}

	return data, nil
}
