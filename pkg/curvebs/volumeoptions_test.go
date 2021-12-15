/*
Copyright 2020 The Netease Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package curvebs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundUpToGiBInt(t *testing.T) {
	var size1Mi int64 = 1 * 1024 * 1024
	size, err := roundUpToGiBInt(size1Mi)
	assert.NoError(t, err)
	assert.Equal(t, 10, size)

	var size20Gi int64 = 20 * 1024 * 1024 * 1024
	size, err = roundUpToGiBInt(size20Gi + size1Mi)
	assert.NoError(t, err)
	assert.Equal(t, 21, size)

	var size4Ti int64 = 4 * 1024 * 1024 * 1024 * 1024
	_, err = roundUpToGiBInt(size4Ti + size1Mi)
	assert.Error(t, err)
}
