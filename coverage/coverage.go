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
	"bytes"
	"io"
	"io/ioutil"
	"strconv"
)

// Coverage contains the overall coverage as determined by reading a
// coverage profile file.
type Coverage struct {
	Total      int64
	Executed   int64
	Unexecuted int64
}

// LoadCoverage reads a stream containing coverage profile data and
// constructs a Coverage from it.
func LoadCoverage(r io.Reader) (cov Coverage, err error) {
	// Begin by reading the data
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	// Interpret the data
	for lno, line := range bytes.Split(data, []byte{'\n'}) {
		// Skip empty lines and the mode directive
		if len(line) == 0 || (lno == 0 && bytes.HasPrefix(line, []byte("mode:"))) {
			continue
		}

		// We're only really interested in the last two
		// fields: the number of statements and the execution
		// count (where we're only wanting to know if it's 0)
		var stmts, runs int
		fields := bytes.Split(line, []byte{' '})
		stmts, err = strconv.Atoi(string(fields[len(fields)-2]))
		if err != nil {
			return
		}
		runs, err = strconv.Atoi(string(fields[len(fields)-1]))
		if err != nil {
			return
		}

		// OK, categorize it as run or unrun and update the
		// coverage data
		cov.Total += int64(stmts)
		if runs == 0 {
			cov.Unexecuted += int64(stmts)
		} else {
			cov.Executed += int64(stmts)
		}
	}

	return
}
