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

package curve

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComposeCSIID(t *testing.T) {
	user := "k8s"
	volName := csiVolNamingPrefix + "pvc-eeafeeb3-7a35-11ea-934a-fa163e28f309"

	id, err := composeCSIID(user, volName)
	assert.NoError(t, err)
	assert.Equal(t, "0003-k8s-csi-vol-pvc-eeafeeb3-7a35-11ea-934a-fa163e28f309", id)
}

func TestDecomposeCSIID(t *testing.T) {
	id := "0003-k8s-csi-vol-pvc-eeafeeb3-7a35-11ea-934a-fa163e28f309"
	user, volName, err := decomposeCSIID(id)
	assert.NoError(t, err)
	assert.Equal(t, "k8s", user)
	assert.Equal(t, "csi-vol-pvc-eeafeeb3-7a35-11ea-934a-fa163e28f309", volName)
}
