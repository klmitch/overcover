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

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileDataCoverageBase(t *testing.T) {
	fd := FileData{
		Count: 100,
		Exec:  50,
	}

	result := fd.Coverage()

	assert.Equal(t, .5, result)
}

func TestFileDataCoverageEmpty(t *testing.T) {
	fd := FileData{
		Count: 0,
		Exec:  50,
	}

	result := fd.Coverage()

	assert.Equal(t, 1.0, result)
}

func TestFileDataHandleBase(t *testing.T) {
	fd := FileData{
		Package: "example.com/some/package",
		Name:    "file.go",
	}

	result := fd.Handle()

	assert.Equal(t, "example.com/some/package/file.go", result)
}

func TestFileDataHandlePackage(t *testing.T) {
	fd := FileData{
		Package: "example.com/some/package",
	}

	result := fd.Handle()

	assert.Equal(t, "example.com/some/package/", result)
}

func TestFileDataHandleOverall(t *testing.T) {
	fd := FileData{}

	result := fd.Handle()

	assert.Equal(t, "", result)
}

func TestFileDataStringBase(t *testing.T) {
	fd := FileData{
		Package: "example.com/some/package",
		Name:    "file.go",
		Count:   100,
		Exec:    50,
	}

	result := fd.String()

	assert.Equal(t, "example.com/some/package/file.go: 50 statements out of 100 covered; coverage: 50.0%", result)
}

func TestFileDataStringOverall(t *testing.T) {
	fd := FileData{
		Count: 100,
		Exec:  50,
	}

	result := fd.String()

	assert.Equal(t, "50 statements out of 100 covered; overall coverage: 50.0%", result)
}

func TestDataSetMergeBase(t *testing.T) {
	ds := DataSet{
		FileData{
			Package: "example.com/some/package",
			Name:    "file1.go",
			Count:   30,
			Exec:    20,
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file2.go",
			Count:   30,
		},
	}
	other := DataSet{
		FileData{
			Package: "example.com/other/package",
			Name:    "file3.go",
			Count:   40,
			Exec:    10,
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file2.go",
			Count:   30,
			Exec:    20,
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file1.go",
			Count:   30,
		},
	}

	result, conflict := ds.Merge(other)

	assert.Equal(t, DataSet{
		FileData{
			Package: "example.com/some/package",
			Name:    "file1.go",
			Count:   30,
			Exec:    20,
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file2.go",
			Count:   30,
			Exec:    20,
		},
		FileData{
			Package: "example.com/other/package",
			Name:    "file3.go",
			Count:   40,
			Exec:    10,
		},
	}, result)
	assert.Nil(t, conflict)
}

func TestDataSetMergeConflict(t *testing.T) {
	ds := DataSet{
		FileData{
			Package: "example.com/some/package",
			Name:    "file1.go",
			Count:   30,
			Exec:    20,
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file2.go",
			Count:   30,
		},
	}
	other := DataSet{
		FileData{
			Package: "example.com/other/package",
			Name:    "file3.go",
			Count:   40,
			Exec:    10,
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file2.go",
			Count:   20,
			Exec:    20,
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file1.go",
			Count:   30,
		},
	}

	result, conflict := ds.Merge(other)

	assert.Equal(t, DataSet{
		FileData{
			Package: "example.com/some/package",
			Name:    "file1.go",
			Count:   30,
			Exec:    20,
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file2.go",
			Count:   30,
		},
		FileData{
			Package: "example.com/other/package",
			Name:    "file3.go",
			Count:   40,
			Exec:    10,
		},
	}, result)
	assert.Equal(t, DataSet{
		FileData{
			Package: "example.com/some/package",
			Name:    "file2.go",
			Count:   20,
			Exec:    20,
		},
	}, conflict)
}

func TestDataSetReduce(t *testing.T) {
	ds := DataSet{
		FileData{
			Package: "example.com/some/package",
			Name:    "file1.go",
			Count:   30,
			Exec:    20,
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file2.go",
			Count:   30,
			Exec:    20,
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file3.go",
			Count:   40,
			Exec:    10,
		},
		FileData{
			Package: "example.com/other/package",
			Name:    "file1.go",
			Count:   30,
			Exec:    20,
		},
		FileData{
			Package: "example.com/other/package",
			Name:    "file2.go",
			Count:   30,
			Exec:    20,
		},
		FileData{
			Package: "example.com/other/package",
			Name:    "file3.go",
			Count:   40,
			Exec:    10,
		},
	}

	result := ds.Reduce()

	assert.Equal(t, DataSet{
		FileData{
			Package: "example.com/some/package",
			Count:   100,
			Exec:    50,
		},
		FileData{
			Package: "example.com/other/package",
			Count:   100,
			Exec:    50,
		},
	}, result)
}

func TestDataSetSum(t *testing.T) {
	ds := DataSet{
		FileData{
			Package: "example.com/some/package",
			Name:    "file1.go",
			Count:   30,
			Exec:    20,
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file2.go",
			Count:   30,
			Exec:    20,
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file3.go",
			Count:   40,
			Exec:    10,
		},
	}

	result := ds.Sum()

	assert.Equal(t, FileData{
		Count: 100,
		Exec:  50,
	}, result)
}

func TestDataSetLen(t *testing.T) {
	ds := DataSet{
		FileData{},
		FileData{},
		FileData{},
	}

	result := ds.Len()

	assert.Equal(t, 3, result)
}

func TestDataSetLessCoverageEqualTrue(t *testing.T) {
	ds := DataSet{
		FileData{
			Package: "example.com/some/package",
			Name:    "file1.go",
			Count:   20,
			Exec:    10,
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file2.go",
			Count:   20,
			Exec:    10,
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file3.go",
			Count:   10,
			Exec:    5,
		},
	}

	result := ds.Less(1, 2)

	assert.True(t, result)
}

func TestDataSetLessCoverageLess(t *testing.T) {
	ds := DataSet{
		FileData{
			Package: "example.com/some/package",
			Name:    "file1.go",
			Count:   20,
			Exec:    10,
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file2.go",
			Count:   20,
			Exec:    9,
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file3.go",
			Count:   10,
			Exec:    5,
		},
	}

	result := ds.Less(1, 2)

	assert.True(t, result)
}

func TestDataSetLessCoverageGreater(t *testing.T) {
	ds := DataSet{
		FileData{
			Package: "example.com/some/package",
			Name:    "file1.go",
			Count:   20,
			Exec:    10,
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file2.go",
			Count:   20,
			Exec:    11,
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file3.go",
			Count:   10,
			Exec:    5,
		},
	}

	result := ds.Less(1, 2)

	assert.False(t, result)
}

func TestDataSetLessCoverageEqualFalse(t *testing.T) {
	ds := DataSet{
		FileData{
			Package: "example.com/some/package",
			Name:    "file1.go",
			Count:   20,
			Exec:    10,
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file2.go",
			Count:   20,
			Exec:    10,
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file3.go",
			Count:   10,
			Exec:    5,
		},
	}

	result := ds.Less(2, 1)

	assert.False(t, result)
}

func TestDataSeSwap(t *testing.T) {
	ds := DataSet{
		FileData{
			Package: "example.com/some/package",
			Name:    "file1.go",
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file2.go",
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file3.go",
		},
	}

	ds.Swap(1, 2)

	assert.Equal(t, DataSet{
		FileData{
			Package: "example.com/some/package",
			Name:    "file1.go",
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file3.go",
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file2.go",
		},
	}, ds)
}

func TestDataSetSorts(t *testing.T) {
	ds := DataSet{
		FileData{
			Package: "example.com/some/package",
			Name:    "file3.go",
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file2.go",
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file1.go",
		},
	}

	sort.Sort(ds)

	assert.Equal(t, DataSet{
		FileData{
			Package: "example.com/some/package",
			Name:    "file1.go",
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file2.go",
		},
		FileData{
			Package: "example.com/some/package",
			Name:    "file3.go",
		},
	}, ds)
}
