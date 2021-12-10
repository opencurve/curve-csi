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
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	volumehelpers "k8s.io/cloud-provider/volume/helpers"

	csicommon "github.com/opencurve/curve-csi/pkg/csi-common"
	"github.com/opencurve/curve-csi/pkg/curveservice"
	"github.com/opencurve/curve-csi/pkg/util"
	"github.com/opencurve/curve-csi/pkg/util/ctxlog"
)

type controllerServer struct {
	*csicommon.DefaultControllerServer

	volumeLocks       *util.VolumeLocks
	curveVolumePrefix string
}

// CreateVolume creates the volume in backend, if it is not already present
func (cs *controllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	if err := cs.validateCreateVolumeRequest(req); err != nil {
		ctxlog.Errorf(ctx, err.Error())
		return nil, err
	}

	reqName := req.GetName()
	if acquired := cs.volumeLocks.TryAcquire(reqName); !acquired {
		ctxlog.Infof(ctx, util.VolumeOperationAlreadyExistsFmt, reqName)
		return nil, status.Errorf(codes.Aborted, util.VolumeOperationAlreadyExistsFmt, reqName)
	}
	defer cs.volumeLocks.Release(reqName)

	ctxlog.Infof(ctx, "starting creating volume requestNamed %s", reqName)
	// get volume options
	volOptions, err := newVolumeOptions(req, cs.curveVolumePrefix)
	if err != nil {
		ctxlog.ErrorS(ctx, err, "failed to new volume options")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// compose csi volume id
	csiVolumeId, err := composeCSIID(volOptions.user, volOptions.volName)
	if err != nil {
		ctxlog.ErrorS(ctx, err, "failed to composeCSIID")
		return nil, status.Error(codes.Internal, err.Error())
	}
	ctxlog.V(5).Infof(ctx, "build volumeOptions: %+v with csiVolumeId: %v", volOptions, csiVolumeId)

	// TODO: support volume clone and snapshot restore

	// verify the volume already exists
	curveVol := curveservice.NewCurveVolume(volOptions.user, volOptions.volName, volOptions.sizeGiB)
	volDetail, err := curveVol.Stat(ctx)
	if err != nil {
		ctxlog.ErrorS(ctx, err, "failed to get volDetail")
		return nil, status.Error(codes.Internal, err.Error())
	}
	if volDetail.FileStatus == curveservice.CurveVolumeStatusCreated {
		ctxlog.V(4).Infof(ctx, "the volume %s already created", reqName)
		return &csi.CreateVolumeResponse{
			Volume: &csi.Volume{
				VolumeId:      csiVolumeId,
				CapacityBytes: int64(volOptions.sizeGiB * volumehelpers.GiB),
				VolumeContext: req.GetParameters(),
			},
		}, nil
	}

	// create volume
	if err := curveVol.Create(ctx); err != nil {
		ctxlog.ErrorS(ctx, err, "failed to create volume")
		return nil, status.Error(codes.Internal, err.Error())
	}

	ctxlog.Infof(ctx, "successfully created volume named %s for request name %s", curveVol.FileName, reqName)
	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      csiVolumeId,
			CapacityBytes: int64(volOptions.sizeGiB * volumehelpers.GiB),
			VolumeContext: req.GetParameters(),
		},
	}, nil
}

// DeleteVolume deletes the volume in backend and its reservation
func (cs *controllerServer) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	if err := cs.validateDeleteVolumeRequest(req); err != nil {
		ctxlog.ErrorS(ctx, err, "DeleteVolumeRequest validation failed")
		return nil, err
	}

	volumeId := req.GetVolumeId()
	// lock out parallel delete operations
	if acquired := cs.volumeLocks.TryAcquire(volumeId); !acquired {
		ctxlog.Infof(ctx, util.VolumeOperationAlreadyExistsFmt, volumeId)
		return nil, status.Errorf(codes.Aborted, util.VolumeOperationAlreadyExistsFmt, volumeId)
	}
	defer cs.volumeLocks.Release(volumeId)

	ctxlog.Infof(ctx, "starting deleting volume id %s", volumeId)

	volOptions, err := newVolumeOptionsFromVolID(volumeId, cs.curveVolumePrefix)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// lock out parallel delete and create requests against the same volume name
	if acquired := cs.volumeLocks.TryAcquire(volOptions.reqName); !acquired {
		ctxlog.Infof(ctx, util.VolumeOperationAlreadyExistsFmt, volumeId)
		return nil, status.Errorf(codes.Aborted, util.VolumeOperationAlreadyExistsFmt, volOptions.reqName)
	}
	defer cs.volumeLocks.Release(volOptions.reqName)

	curveVol := curveservice.NewCurveVolume(volOptions.user, volOptions.volName, volOptions.sizeGiB)
	if err := curveVol.Delete(ctx); err != nil {
		ctxlog.ErrorS(ctx, err, "failed to delete volume", "volumeId", volumeId)
		return nil, status.Error(codes.Internal, err.Error())
	}

	ctxlog.Infof(ctx, "successfully deleted volume %s", volumeId)
	return &csi.DeleteVolumeResponse{}, nil
}

func (cs *controllerServer) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	if err := cs.validateExpandVolumeRequest(req); err != nil {
		ctxlog.ErrorS(ctx, err, "ExpandVolumeRequest validation failed")
		return nil, err
	}
	reqSizeGiB, err := roundUpToGiBInt(req.GetCapacityRange().GetRequiredBytes())
	if err != nil {
		return nil, err
	}

	volumeId := req.GetVolumeId()

	// lock out parallel requests against the same volume ID
	if acquired := cs.volumeLocks.TryAcquire(volumeId); !acquired {
		ctxlog.Infof(ctx, util.VolumeOperationAlreadyExistsFmt, volumeId)
		return nil, status.Errorf(codes.Aborted, util.VolumeOperationAlreadyExistsFmt, volumeId)
	}
	defer cs.volumeLocks.Release(volumeId)

	volOptions, err := newVolumeOptionsFromVolID(volumeId, cs.curveVolumePrefix)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// lock out parallel delete/create/expand requests against the same volume name
	if acquired := cs.volumeLocks.TryAcquire(volOptions.reqName); !acquired {
		ctxlog.Infof(ctx, util.VolumeOperationAlreadyExistsFmt, volumeId)
		return nil, status.Errorf(codes.Aborted, util.VolumeOperationAlreadyExistsFmt, volOptions.reqName)
	}
	defer cs.volumeLocks.Release(volOptions.reqName)

	curveVol := curveservice.NewCurveVolume(volOptions.user, volOptions.volName, reqSizeGiB)
	// get volume information
	volDetail, err := curveVol.Stat(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if volDetail.FileStatus == curveservice.CurveVolumeStatusNotExist {
		return nil, status.Errorf(codes.Internal, "the curve volume %s not exists", volOptions.volName)
	}

	ctxlog.Infof(ctx, "volume %s(status %s) size is %dGiB, reqSize is round up to %dGiB",
		volDetail.FileName, volDetail.FileStatus, volDetail.LengthGiB, reqSizeGiB)
	if reqSizeGiB <= volDetail.LengthGiB {
		return &csi.ControllerExpandVolumeResponse{
			CapacityBytes:         int64(volDetail.LengthGiB * volumehelpers.GiB),
			NodeExpansionRequired: false,
		}, nil
	}

	if err := curveVol.Extend(ctx, reqSizeGiB); err != nil {
		ctxlog.ErrorS(ctx, err, "failed to delete volume", "volumeId", volumeId)
		return nil, status.Error(codes.Internal, err.Error())
	}

	ctxlog.Infof(ctx, "successfully extend volume %s size to %dGiB", volDetail.FileName, reqSizeGiB)
	return &csi.ControllerExpandVolumeResponse{
		CapacityBytes:         int64(reqSizeGiB * volumehelpers.GiB),
		NodeExpansionRequired: true,
	}, nil
}

// ValidateVolumeCapabilities checks whether the volume capabilities requested are supported.
func (cs *controllerServer) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	if req.GetVolumeId() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty volume ID in request")
	}

	if len(req.VolumeCapabilities) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty volume capabilities in request")
	}

	return &csi.ValidateVolumeCapabilitiesResponse{
		Confirmed: &csi.ValidateVolumeCapabilitiesResponse_Confirmed{
			VolumeCapabilities: req.VolumeCapabilities,
		},
	}, nil
}
