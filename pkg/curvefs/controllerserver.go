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

package curvefs

import (
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	csicommon "github.com/opencurve/curve-csi/pkg/csi-common"
	"github.com/opencurve/curve-csi/pkg/util"
	"github.com/opencurve/curve-csi/pkg/util/ctxlog"
)

type controllerServer struct {
	*csicommon.DefaultControllerServer

	// A map storing all volumes with ongoing operations so that additional operations
	// for that same volume (as defined by VolumeID/volume name) return an Aborted error
	volumeLocks *util.VolumeLocks
	// A map storing all snapshots with ongoing operations so that additional operations
	// for that same snapshot (as defined by SnapshotID/snapshot name) return an Aborted error
	snapshotLocks *util.VolumeLocks
}

// CreateVolume creates the volume in backend, if it is not already present
func (cs *controllerServer) CreateVolume(
	ctx context.Context,
	req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
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
	// TODO
	return nil, status.Error(codes.Unimplemented, "")
}

// DeleteVolume deletes the volume in backend and its reservation
func (cs *controllerServer) DeleteVolume(
	ctx context.Context,
	req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
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
	// TODO
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *controllerServer) ControllerExpandVolume(
	ctx context.Context,
	req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	if err := cs.validateExpandVolumeRequest(req); err != nil {
		ctxlog.ErrorS(ctx, err, "ExpandVolumeRequest validation failed")
		return nil, err
	}

	volumeId := req.GetVolumeId()
	// lock out parallel requests against the same volume ID
	if acquired := cs.volumeLocks.TryAcquire(volumeId); !acquired {
		ctxlog.Infof(ctx, util.VolumeOperationAlreadyExistsFmt, volumeId)
		return nil, status.Errorf(codes.Aborted, util.VolumeOperationAlreadyExistsFmt, volumeId)
	}
	defer cs.volumeLocks.Release(volumeId)

	// TODO
	return nil, status.Error(codes.Unimplemented, "")
}

// CreateSnapshot creates the snapshot in backend.
func (cs *controllerServer) CreateSnapshot(
	ctx context.Context,
	req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	if err := cs.validateSnapshotReq(req); err != nil {
		ctxlog.ErrorS(ctx, err, "CreateSnapshotRequest validation failed")
		return nil, err
	}

	snapshotName := req.GetName()
	// lock out parallel snapshot operations
	if acquired := cs.snapshotLocks.TryAcquire(snapshotName); !acquired {
		ctxlog.Infof(ctx, util.SnapshotOperationAlreadyExistsFmt, snapshotName)
		return nil, status.Errorf(codes.Aborted, util.SnapshotOperationAlreadyExistsFmt, snapshotName)
	}
	defer cs.snapshotLocks.Release(snapshotName)

	return nil, status.Error(codes.Unimplemented, "")
}

// DeleteSnapshot deletes thesnapshot in backend.
func (cs *controllerServer) DeleteSnapshot(
	ctx context.Context,
	req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	if err := cs.validateDeleteSnapshotReq(req); err != nil {
		ctxlog.ErrorS(ctx, err, "DeleteSnapshotRequest validation failed")
		return nil, err
	}

	snapshotId := req.GetSnapshotId()
	// lock out parallel snapshot
	if acquired := cs.snapshotLocks.TryAcquire(snapshotId); !acquired {
		ctxlog.Errorf(ctx, util.SnapshotOperationAlreadyExistsFmt, snapshotId)
		return nil, status.Errorf(codes.Aborted, util.SnapshotOperationAlreadyExistsFmt, snapshotId)
	}
	defer cs.snapshotLocks.Release(snapshotId)

	// TODO
	return nil, status.Error(codes.Unimplemented, "")
}

// ValidateVolumeCapabilities checks whether the volume capabilities requested are supported.
func (cs *controllerServer) ValidateVolumeCapabilities(
	ctx context.Context,
	req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	volumeId := req.GetVolumeId()
	if volumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "empty volume ID in request")
	}

	if len(req.VolumeCapabilities) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty volume capabilities in request")
	}

	// TODO
	return nil, status.Error(codes.Unimplemented, "")
}
