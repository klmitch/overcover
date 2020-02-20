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

package common

import "fmt"

// FileData contains the summarized data about the file, including its
// package, its base name, the total number of statements, and the
// number of executed statements.
type FileData struct {
	Package string // Name of the package
	Name    string // Name of the file (basename)
	Count   int64  // Number of statements in the file
	Exec    int64  // Number of statements in the file that were executed
}

// Coverage reports the coverage of the file as a float.
func (fd FileData) Coverage() float64 {
	// Avoid divide-by-zero
	if fd.Count <= 0 {
		return 1.0
	}

	return float64(fd.Exec) / float64(fd.Count)
}

// Handle reports a handle for the FileData record.  This is typically
// the full file name, including package, but it could be just the
// package name, or the empty string.
func (fd FileData) Handle() string {
	if fd.Package == "" && fd.Name == "" {
		return ""
	}

	return fmt.Sprintf("%s/%s", fd.Package, fd.Name)
}

// String reports the coverage of the entity described by the FileData
// object.
func (fd FileData) String() string {
	handle := fd.Handle()
	if handle == "" {
		return fmt.Sprintf("%d statements out of %d covered; overall coverage: %.1f%%", fd.Exec, fd.Count, fd.Coverage()*100.0)
	}

	return fmt.Sprintf("%s: %d statements out of %d covered; coverage: %.1f%%", handle, fd.Exec, fd.Count, fd.Coverage()*100.0)
}

// DataSet is a list of FileData instances.  It implements
// sort.Interface, and thus is capable of being sorted by handle.
type DataSet []FileData

// Merge is a utility function that merges a list of FileData
// instances with another FileData list.  It ensures that the Count
// fields are the same and picks an Exec field.  It returns the
// resulting list, along with a list of FileData where Count did not
// match; the order will match the elements of a first, followed by
// any elements of b which did not appear in a.
func (ds DataSet) Merge(other DataSet) (DataSet, DataSet) {
	// Set up our index and result sets
	idx := map[string]map[string]int{}
	var result DataSet
	var conflict DataSet

	// Copy ds into the result set
	i := 0
	for _, fd := range ds {
		// Mark it seen
		if _, ok := idx[fd.Package]; !ok {
			idx[fd.Package] = map[string]int{}
		}
		idx[fd.Package][fd.Name] = i

		// Save it to the result set
		result = append(result, fd)
		i++
	}

	// Now walk through other
	for _, fd := range other {
		// Have we seen it?
		if _, ok := idx[fd.Package]; !ok {
			idx[fd.Package] = map[string]int{}
		}
		if j, ok := idx[fd.Package][fd.Name]; ok {
			// Are the counts consistent?
			if result[j].Count != fd.Count {
				conflict = append(conflict, fd)
			} else if result[j].Exec == 0 {
				result[j].Exec = fd.Exec
			}
			continue
		}

		// Add the element to the result set
		idx[fd.Package][fd.Name] = i
		result = append(result, fd)
		i++
	}

	return result, conflict
}

// Reduce is a utility function that reduces a list of FileData
// instances into a list of FileData instances that only contain
// counts for packages.
func (ds DataSet) Reduce() DataSet {
	// Set up our index and result set
	idx := map[string]int{}
	var result DataSet

	// Construct it
	i := 0
	for _, fd := range ds {
		// Has the package been seen yet?
		if j, ok := idx[fd.Package]; ok {
			result[j].Count += fd.Count
			result[j].Exec += fd.Exec
			continue
		}

		// Add a new entry
		idx[fd.Package] = i
		result = append(result, FileData{
			Package: fd.Package,
			Count:   fd.Count,
			Exec:    fd.Exec,
		})
		i++
	}

	return result
}

// Sum is a utility function similar to Reduce, but it reduces a list
// of FileData instances down to a single summarizing FileData.
func (ds DataSet) Sum() FileData {
	result := FileData{}
	for _, fd := range ds {
		result.Count += fd.Count
		result.Exec += fd.Exec
	}

	return result
}

// Len returns the number of elements in the data set.
func (ds DataSet) Len() int {
	return len(ds)
}

// Less reports whether the element with index i should sort before
// the element with index j.
func (ds DataSet) Less(i, j int) bool {
	iCov := ds[i].Coverage()
	jCov := ds[j].Coverage()
	switch {
	case iCov < jCov:
		return true
	case jCov < iCov:
		return false
	default:
		return ds[i].Handle() < ds[j].Handle()
	}
}

// Swap swaps the two elements with the specified indexes.
func (ds DataSet) Swap(i, j int) {
	ds[i], ds[j] = ds[j], ds[i]
}
