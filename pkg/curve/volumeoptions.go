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
	"fmt"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/apimachinery/pkg/api/resource"
	volumehelpers "k8s.io/cloud-provider/volume/helpers"
)

const (
	csiDefaultVolNamingPrefix = "csi-vol-"
)

type volumeOptions struct {
	reqName string
	volName string
	sizeGiB int
	user    string
}

func newVolumeOptions(req *csi.CreateVolumeRequest, curveVolumePrefix string) (*volumeOptions, error) {
	var (
		ok  bool
		err error
	)
	opts := &volumeOptions{
		reqName: req.GetName(),
	}

	if curveVolumePrefix != "" {
		opts.volName = curveVolumePrefix + opts.reqName
	} else {
		opts.volName = csiDefaultVolNamingPrefix + opts.reqName
	}

	volOptions := req.GetParameters()
	opts.user, ok = volOptions["user"]
	if !ok {
		return nil, fmt.Errorf("missing required field: user")
	}

	// volume size - default is 10GiB
	opts.sizeGiB = 10
	if req.GetCapacityRange() != nil {
		opts.sizeGiB, err = roundUpToGiBInt(req.GetCapacityRange().GetRequiredBytes())
		if err != nil {
			return nil, err
		}
	}

	return opts, nil
}

func newVolumeOptionsFromVolID(volumeId string, curveVolumePrefix string) (*volumeOptions, error) {
	var (
		volOptions volumeOptions
		err        error
	)

	volOptions.user, volOptions.volName, err = decomposeCSIID(volumeId)
	if err != nil {
		return nil, err
	}
	volNamingPrefix := curveVolumePrefix
	if volNamingPrefix == "" {
		volNamingPrefix = csiDefaultVolNamingPrefix
	}
	volOptions.reqName = strings.TrimPrefix(volOptions.volName, volNamingPrefix)

	return &volOptions, nil
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
