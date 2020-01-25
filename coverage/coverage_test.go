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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadCoverageBase(t *testing.T) {
	r := bytes.NewBufferString(`mode: atomic
file1.go:1.1,3.3 3 5
file1.go:4.4,5.5 1 0
file1.go:6.6,7.7 1 1
file2.go:1.1,3.3 3 0
file2.go:4.4,5.5 1 15
`)

	result, err := LoadCoverage(r)

	assert.NoError(t, err)
	assert.Equal(t, int64(9), result.Total)
	assert.Equal(t, int64(5), result.Executed)
	assert.Equal(t, int64(4), result.Unexecuted)
}

type FailingReader struct{}

func (r *FailingReader) Read(p []byte) (int, error) {
	return 0, assert.AnError
}

func TestLoadCoverageReadAllFails(t *testing.T) {
	r := &FailingReader{}

	_, err := LoadCoverage(r)

	assert.Equal(t, assert.AnError, err)
}

func TestLoadCoverageStatementsFails(t *testing.T) {
	r := bytes.NewBufferString(`mode: atomic
file1.go:1.1,3.3 3 5
file1.go:4.4,5.5 1 0
file1.go:6.6,7.7 1 1
file2.go:1.1,3.3 bad 0
file2.go:4.4,5.5 1 15
`)

	_, err := LoadCoverage(r)

	assert.NotNil(t, err)
}

func TestLoadCoverageRunsFails(t *testing.T) {
	r := bytes.NewBufferString(`mode: atomic
file1.go:1.1,3.3 3 5
file1.go:4.4,5.5 1 0
file1.go:6.6,7.7 1 1
file2.go:1.1,3.3 3 bad
file2.go:4.4,5.5 1 15
`)

	_, err := LoadCoverage(r)

	assert.NotNil(t, err)
}
