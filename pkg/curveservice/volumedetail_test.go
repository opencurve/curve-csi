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

package curveservice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleParseVolumeDetail(t *testing.T) {
	type statCase struct {
		output    string
		volDetail CurveVolumeDetail
	}
	validStatus := []statCase{
		{
			output: `id: 39007
parentid: 39005
filetype: INODE_PAGEFILE
length(GB): 10
createtime: 2020-08-07 10:51:52
user: k8s
filename: pvc-ce482926-91d8-11ea-bf6e-fa163e23ce53
fileStatus: Created
`,
			volDetail: CurveVolumeDetail{
				Id:         "39007",
				ParentId:   "39005",
				FileType:   "INODE_PAGEFILE",
				LengthGiB:  10,
				CreateTime: "2020-08-07 10:51:52",
				User:       "k8s",
				FileName:   "pvc-ce482926-91d8-11ea-bf6e-fa163e23ce53",
				FileStatus: "Created",
			},
		}, {
			output: `
parentid: 39005
filetype: INODE_PAGEFILE
length(GB): 10
createtime: 2020-08-07 10:51:52
user: k8s
filename: pvc-ce482926-91d8-11ea-bf6e-fa163e23ce53
fileStatus: Created
paramA: valueA
`,
			volDetail: CurveVolumeDetail{
				Id:         "",
				ParentId:   "39005",
				FileType:   "INODE_PAGEFILE",
				LengthGiB:  10,
				CreateTime: "2020-08-07 10:51:52",
				User:       "k8s",
				FileName:   "pvc-ce482926-91d8-11ea-bf6e-fa163e23ce53",
				FileStatus: "Created",
			},
		},
	}

	for _, tcase := range validStatus {
		vDetail, err := simpleParseVolumeDetail([]byte(tcase.output))
		assert.NoError(t, err)
		assert.Equal(t, tcase.volDetail, vDetail)
	}

	invalidStatCase := statCase{
		output: `id: 39007
parentid: 39005
filetype: INODE_PAGEFILE
length(GB): a
createtime: 2020-08-07 10:51:52
user: k8s
filename: pvc-ce482926-91d8-11ea-bf6e-fa163e23ce53
fileStatus: Created`,
	}

	_, err := simpleParseVolumeDetail([]byte(invalidStatCase.output))
	assert.Error(t, err)
}
