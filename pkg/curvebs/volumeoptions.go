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
	"fmt"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/apimachinery/pkg/api/resource"
	volumehelpers "k8s.io/cloud-provider/volume/helpers"
)

const (
	csiVolNamingPrefix = "csi-vol-"

	// max length of curve volume uesr
	curveUserMaxLen = 30
	// clone lazy
	curveCloneDefaultLazy = true
)

type volumeOptions struct {
	reqName   string
	volName   string
	volId     string
	sizeGiB   int
	user      string
	cloneLazy bool
}

func (vo *volumeOptions) genVolumePath() string {
	return "/" + vo.user + "/" + vo.volName
}

func newVolumeOptions(req *csi.CreateVolumeRequest) (*volumeOptions, error) {
	var (
		ok  bool
		err error
	)
	opts := &volumeOptions{
		reqName: req.GetName(),
		volName: csiVolNamingPrefix + req.GetName(),
	}

	parameters := req.GetParameters()
	opts.user, ok = parameters["user"]
	if !ok {
		return nil, fmt.Errorf("missing required field: user")
	}
	if len(opts.user) == 0 || len(opts.user) > curveUserMaxLen {
		return nil, fmt.Errorf("length of field user must be 1~%v", curveUserMaxLen)
	}

	cloneLazy, ok := parameters["cloneLazy"]
	if ok {
		opts.cloneLazy = cloneLazy == "true"
	} else {
		opts.cloneLazy = curveCloneDefaultLazy
	}

	// volume size - default is 10GiB
	opts.sizeGiB = 10
	if req.GetCapacityRange() != nil {
		opts.sizeGiB, err = roundUpToGiBInt(req.GetCapacityRange().GetRequiredBytes())
		if err != nil {
			return nil, err
		}
	}

	opts.volId, err = composeCSIID(opts.user, opts.volName)
	if err != nil {
		return nil, fmt.Errorf("failed to composeCSIID: %v", err)
	}

	return opts, nil
}

func newVolumeOptionsFromVolID(volumeId string) (*volumeOptions, error) {
	var err error
	volOptions := &volumeOptions{
		volId: volumeId,
	}
	volOptions.user, volOptions.volName, err = decomposeCSIID(volumeId)
	if err != nil {
		return nil, err
	}
	volOptions.reqName = strings.TrimPrefix(volOptions.volName, csiVolNamingPrefix)

	return volOptions, nil
}

func roundUpToGiBInt(sizeBytes int64) (int, error) {
	quantity := resource.NewQuantity(sizeBytes, resource.BinarySI)
	sizeGiB, err := volumehelpers.RoundUpToGiBInt(*quantity)
	if err != nil {
		return 0, err
	}
	if sizeGiB > 4*1024 {
		return 0, fmt.Errorf("the volume size must be less than 4TiB, got %dGiB", sizeGiB)
	}
	// must more than 10GiB
	if sizeGiB < 10 {
		sizeGiB = 10
	}
	return sizeGiB, nil
}

func parseSnapshotID(snapshotId string) (string, *volumeOptions, error) {
	snapCurveUUID, volId, err := decomposeSnapshotID(snapshotId)
	if err != nil {
		return "", nil, err
	}
	volOptions, err := newVolumeOptionsFromVolID(volId)

	return snapCurveUUID, volOptions, err
}
