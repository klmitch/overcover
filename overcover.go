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
	"flag"
	"fmt"
	"os"

	"github.com/klmitch/overcover/coverage"
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
	file, err := os.Open(coverprofile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open coverage profile file: %s\n", err)
		os.Exit(2)
	}
	defer file.Close()
	cov, err := coverage.LoadCoverage(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read coverage profile file: %s\n", err)
		os.Exit(2)
	}

	// Compute the overall coverage
	coverage := float64(cov.Executed) / float64(cov.Total) * 100.0
	fmt.Printf("%d statements out of %d covered; overall coverage: %.1f%%\n", cov.Executed, cov.Total, coverage)

	// Verify that we met the threshold
	if threshold > 0.0 && coverage < threshold {
		fmt.Fprintf(os.Stderr, "\nFailed to meet coverage threshold of %.1f%%\n", threshold)
		os.Exit(1)
	}
}
