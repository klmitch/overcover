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

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

var (
	threshold    float64
	coverprofile string
)

func init() {
	flag.Float64Var(&threshold, "threshold", 0, "Set the minimum threshold for coverage; coverage below this threshold will result in an error.")
	flag.StringVar(&coverprofile, "coverprofile", "coverage.out", "Specify the coverage profile file to read.")
}

func main() {
	flag.Parse()

	// Load the coverage; this reads the coverage profile and sums the
	// statement counts
	exec, unexec, err := loadCoverage(coverprofile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read coverage profile file: %s\n", err)
		os.Exit(2)
	}

	// Compute the overall coverage
	coverage := float64(exec) / float64(exec+unexec) * 100.0
	fmt.Printf("%d statements out of %d covered; overall coverage: %.1f%%\n", exec, exec+unexec, coverage)

	// Verify that we met the threshold
	if threshold > 0.0 && coverage < threshold {
		fmt.Fprintf(os.Stderr, "\nFailed to meet coverage threshold of %.1f%%\n", threshold)
		os.Exit(1)
	}
}

func loadCoverage(filename string) (exec int64, unexec int64, err error) {
	// Read the coverage file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	// Split it up
	exec = 0
	unexec = 0
	for lno, line := range bytes.Split(data, []byte{'\n'}) {
		if len(line) == 0 || (lno == 0 && bytes.HasPrefix(line, []byte("mode:"))) {
			continue
		}

		// Parse out the last two fields
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

		if runs == 0 {
			unexec += int64(stmts)
		} else {
			exec += int64(stmts)
		}
	}

	return
}
