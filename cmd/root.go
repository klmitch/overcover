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
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/klmitch/overcover/common"
	"github.com/klmitch/overcover/coverage"
	"github.com/klmitch/overcover/statements"
)

// Variables used to store the values of flags.
var (
	config       string
	readOnly     bool
	coverprofile string
	buildArgs    = []string{}
	detailed     bool
	summary      bool
)

// Variables used for mocking for the tests.
var (
	stdout         io.Writer                                        = os.Stdout
	stderr         io.Writer                                        = os.Stderr
	exit                                                            = os.Exit
	getFloat64     func(string) float64                             = viper.GetFloat64
	setConfigFile  func(string)                                     = viper.SetConfigFile
	readInConfig   func() error                                     = viper.ReadInConfig
	configFileUsed func() string                                    = viper.ConfigFileUsed
	setConfig      func(string, interface{})                        = viper.Set
	writeConfig    func(string) error                               = viper.WriteConfigAs
	loadCoverage   func(string) (common.DataSet, error)             = coverage.Load
	loadStatements func([]string, []string) (common.DataSet, error) = statements.Load
)

// rootCmd describes the overcover command to cobra.
var rootCmd = &cobra.Command{
	Use:   "overcover [flags] [PACKAGE ...]",
	Short: "Golang overall coverage tool with threshold enforcement",
	Long:  `A tool for reporting and testing the overall test suite coverage of a test suite written in go.  This parses the coverage profile output file (generated by passing a filename to the "-coverprofile" option of "go test") and reports the overall coverage of the test suite.  It can also test that the coverage meets a certain minimum threshold.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load the coverage; this reads the coverage profile
		// and sums the statement counts
		if coverprofile == "" && len(args) == 0 {
			fmt.Fprintf(stderr, "No coverage profile file specified!  Use -p or provide a configuration file.\n")
			_ = cmd.Usage()
			exit(2)
		}
		var ds common.DataSet
		if coverprofile != "" {
			var err error
			ds, err = loadCoverage(coverprofile)
			if err != nil {
				fmt.Fprintf(stderr, "Unable to read coverage profile file %q: %s\n", coverprofile, err)
				exit(3)
			}
		}

		// Next, read in the source if requested to
		if len(args) > 0 {
			direct, err := loadStatements(buildArgs, args)
			if err != nil {
				fmt.Fprintf(stderr, "Unable to read source: %s\n", err)
				exit(3)
			}

			// Merge the direct-read data
			var conflict common.DataSet
			ds, conflict = ds.Merge(direct)
			if len(conflict) > 0 {
				fmt.Fprintf(stderr, "WARNING: coverage profile %s may not match source; potentially altered files:\n", coverprofile)
				for _, fd := range conflict {
					fmt.Fprintf(stderr, "  %s\n", fd.Handle())
				}
			}
		}

		// Emit summary data, if requested
		if summary {
			summary := ds.Reduce()
			sort.Sort(summary)
			fmt.Fprintln(stdout, "Per package summary:")
			tab := tabwriter.NewWriter(stdout, 2, 8, 2, ' ', 0)
			fmt.Fprintf(tab, " Package\tExecuted\tTotal\tCoverage\n -------\t--------\t-----\t--------\n")
			for _, rec := range summary {
				fmt.Fprintf(tab, " %s\t%d\t%d\t%.1f%%\n", rec.Handle(), rec.Exec, rec.Count, rec.Coverage()*100.0)
			}
			tab.Flush()
			fmt.Fprintln(stdout, "")
		}

		// Emit detailed data, if requested
		if detailed {
			sort.Sort(ds)
			fmt.Fprintln(stdout, "Per file details:")
			tab := tabwriter.NewWriter(stdout, 2, 8, 2, ' ', 0)
			fmt.Fprintf(tab, " File\tExecuted\tTotal\tCoverage\n ----\t--------\t-----\t--------\n")
			for _, rec := range ds {
				fmt.Fprintf(tab, " %s\t%d\t%d\t%.1f%%\n", rec.Handle(), rec.Exec, rec.Count, rec.Coverage()*100.0)
			}
			tab.Flush()
			fmt.Fprintln(stdout, "")
		}

		// Compute the overall coverage record
		overall := ds.Sum()
		fmt.Fprintln(stdout, overall)

		// Verify that we met the threshold
		coverage := overall.Coverage() * 100.0
		threshold := getFloat64("threshold")
		if threshold > 0.0 && coverage < threshold {
			fmt.Fprintf(stderr, "\nFailed to meet coverage threshold of %.1f%%\n", threshold)
			exit(1)
		}

		// OK, now let's see if the threshold needs updating
		minHeadroom := getFloat64("min_headroom")
		maxHeadroom := getFloat64("max_headroom")
		if config != "" && minHeadroom >= 0.0 && maxHeadroom > minHeadroom && coverage > threshold+maxHeadroom {
			// Compute new threshold
			newThreshold := math.Round(coverage*10.0)/10.0 - minHeadroom

			// If we're read-only, generate an error
			if readOnly {
				fmt.Fprintf(stderr, "\nCoverage exceeds maximum headroom.  Update threshold to %.1f%%\n", newThreshold)
				exit(5)
			}

			// OK, update the configuration
			fmt.Fprintf(stdout, "Updating configuration file %s with new threshold value %.1f%%\n", config, newThreshold)
			setConfig("threshold", newThreshold)
			if err := writeConfig(config); err != nil {
				fmt.Fprintf(stderr, "\nFailed to write updated config with new threshold %.1f%% to %s: %s\n", newThreshold, config, err)
				exit(5)
			}
		}
	},
}

// Execute is the entrypoint for overcover.  This invokes the root
// command, which performs all the work.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(stderr, err)
		exit(4)
	}
}

// getBuildArgDefault is a helper that retrieves the default build
// arguments from the environment.
func getBuildArgDefault() []string {
	data, ok := os.LookupEnv("OVERCOVER_BUILD_ARG")
	if !ok {
		return []string{}
	}

	return strings.Split(data, " ")
}

// init initializes the flags for overcover.
func init() {
	// Initialize cobra and viper
	cobra.OnInitialize(readConfig)
	viper.SetEnvPrefix("overcover")

	// Set up the flags
	rootCmd.Flags().StringVarP(&config, "config", "c", os.Getenv("OVERCOVER_CONFIG"), "Configuration file to read.  All command line options may be set through the configuration file.")
	_, readOnlyDefault := os.LookupEnv("OVERCOVER_READONLY")
	rootCmd.Flags().BoolVarP(&readOnly, "readonly", "r", readOnlyDefault, "Used to indicate that the configuration file should only be read, not written.")
	rootCmd.Flags().Float64P("threshold", "t", 0, "Set the minimum threshold for coverage; coverage below this threshold will result in an error.")
	rootCmd.Flags().StringVarP(&coverprofile, "coverprofile", "p", os.Getenv("OVERCOVER_COVERPROFILE"), "Specify the coverage profile file to read.")
	rootCmd.Flags().Float64P("min-headroom", "m", 0, "Set the minimum headroom.  If the threshold is raised, it will be raised to the current coverage minus this value.")
	rootCmd.Flags().Float64P("max-headroom", "M", 0, "Set the maximum headroom.  If the coverage is more than the threshold plus this value, the threshold will be raised.")
	rootCmd.Flags().StringArrayVarP(&buildArgs, "build-arg", "b", getBuildArgDefault(), "Add a build argument.  Build arguments are used to select source files for later coverage checking.")
	_, detailedDefault := os.LookupEnv("OVERCOVER_DETAILED")
	rootCmd.Flags().BoolVarP(&detailed, "detailed", "d", detailedDefault, "Used to request per-file detailed coverage data be emitted.  May be used in conjunction with --summary.")
	_, summaryDefault := os.LookupEnv("OVERCOVER_SUMMARY")
	rootCmd.Flags().BoolVarP(&summary, "summary", "s", summaryDefault, "Used to request per-package summary coverage data be emitted.  May be used in conjunction with --detailed.")

	// Bind them to viper
	_ = viper.BindPFlag("threshold", rootCmd.Flags().Lookup("threshold"))
	_ = viper.BindEnv("threshold")
	_ = viper.BindPFlag("min_headroom", rootCmd.Flags().Lookup("min-headroom"))
	_ = viper.BindEnv("min_headroom")
	_ = viper.BindPFlag("max_headroom", rootCmd.Flags().Lookup("max-headroom"))
	_ = viper.BindEnv("max_headroom")
}

// readConfig reads the configuration file using Viper.
func readConfig() {
	// Is a configuration file set?
	if config == "" {
		return
	}

	// Select it
	setConfigFile(config)

	// Read the configuration
	if err := readInConfig(); err == nil {
		fmt.Fprintf(stdout, "Using configuration file %s\n", configFileUsed())
	}
}
