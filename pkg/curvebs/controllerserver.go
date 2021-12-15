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
	"context"
	"fmt"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	volumehelpers "k8s.io/cloud-provider/volume/helpers"

	csicommon "github.com/opencurve/curve-csi/pkg/csi-common"
	"github.com/opencurve/curve-csi/pkg/curveservice"
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

	snapshotServer string
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
	// get volume options
	volOptions, err := newVolumeOptions(req)
	if err != nil {
		ctxlog.ErrorS(ctx, err, "failed to new volume options")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	ctxlog.V(5).Infof(ctx, "build volumeOptions: %+v", volOptions)

	// verify the volume already exists
	curveVol := curveservice.NewCurveVolume(volOptions.user, volOptions.volName, volOptions.sizeGiB)
	volDetail, err := curveVol.Stat(ctx)
	if err == nil {
		ctxlog.V(4).Infof(ctx, "the volume %v already created, status: %v", volOptions.volName, volDetail.FileStatus)
		if volDetail.LengthGiB != volOptions.sizeGiB {
			return nil, status.Errorf(codes.AlreadyExists, "request size %vGiB not equal with existing %vGiB", volOptions.sizeGiB, volDetail.LengthGiB)
		}
		return &csi.CreateVolumeResponse{
			Volume: &csi.Volume{
				VolumeId:      volOptions.volId,
				CapacityBytes: int64(volOptions.sizeGiB * volumehelpers.GiB),
				VolumeContext: req.GetParameters(),
			},
		}, nil
	}
	if !util.IsNotFoundErr(err) {
		ctxlog.ErrorS(ctx, err, "failed to get volDetail")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// create volume from contentSource: snapshot or clone from an existing volume
	volSource, err := cs.createVolFromContentSource(ctx, req, volOptions, curveVol)
	if err != nil {
		return nil, err
	}
	if len(volSource) > 0 {
		volContext := req.GetParameters()
		volContext["volSource"] = volSource
		return &csi.CreateVolumeResponse{
			Volume: &csi.Volume{
				VolumeId:      volOptions.volId,
				CapacityBytes: int64(volOptions.sizeGiB * volumehelpers.GiB),
				ContentSource: req.GetVolumeContentSource(),
				VolumeContext: volContext,
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
			VolumeId:      volOptions.volId,
			CapacityBytes: int64(volOptions.sizeGiB * volumehelpers.GiB),
			VolumeContext: req.GetParameters(),
		},
	}, nil
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

	volOptions, err := newVolumeOptionsFromVolID(volumeId)
	if err != nil {
		ctxlog.Warningf(ctx, "failed to new volOptions from volume id %v", volumeId)
		return &csi.DeleteVolumeResponse{}, nil
	}
	// lock out parallel delete and create requests against the same volume name
	if acquired := cs.volumeLocks.TryAcquire(volOptions.reqName); !acquired {
		ctxlog.Infof(ctx, util.VolumeOperationAlreadyExistsFmt, volumeId)
		return nil, status.Errorf(codes.Aborted, util.VolumeOperationAlreadyExistsFmt, volOptions.reqName)
	}
	defer cs.volumeLocks.Release(volOptions.reqName)

	if cs.snapshotServer == "" {
		// delete volume
		curveVol := curveservice.NewCurveVolume(volOptions.user, volOptions.volName, volOptions.sizeGiB)
		if err := curveVol.Delete(ctx); err != nil {
			ctxlog.ErrorS(ctx, err, "failed to delete volume", "volumeId", volumeId)
			return nil, status.Error(codes.Internal, err.Error())
		}
		ctxlog.Infof(ctx, "successfully deleted volume %s", volumeId)
		return &csi.DeleteVolumeResponse{}, nil
	}

	// ensure all the tasks created from this volume status done.
	snapServer := curveservice.NewSnapshotServer(cs.snapshotServer, volOptions.user, volOptions.volName)
	if err = snapServer.EnsureTaskFromSourceDone(ctx, volOptions.genVolumePath()); err != nil {
		ctxlog.Errorf(ctx, "failed to ensure tasks from %v status done: %v", volumeId, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// detete volume
	curveVol := curveservice.NewCurveVolume(volOptions.user, volOptions.volName, volOptions.sizeGiB)
	if err := curveVol.Delete(ctx); err != nil {
		ctxlog.ErrorS(ctx, err, "failed to delete volume", "volumeId", volumeId)
		return nil, status.Error(codes.Internal, err.Error())
	}
	ctxlog.Infof(ctx, "successfully deleted volume %s", volumeId)

	// clean cloneTask if the volume is cloned
	taskInfo, err := snapServer.GetCloneTaskOfDestination(ctx, volOptions.genVolumePath())
	if err != nil {
		if util.IsNotFoundErr(err) {
			ctxlog.Infof(ctx, "the volume is not cloned, need not clean tasks")
		} else {
			ctxlog.Warningf(ctx, "can not get taskInfo of path %v", volOptions.genVolumePath())
		}
		return &csi.DeleteVolumeResponse{}, nil
	}
	if err = snapServer.CleanCloneTask(ctx, taskInfo.UUID); err != nil {
		ctxlog.Warningf(ctx, "can not clean task %v", taskInfo.UUID)
	}
	return &csi.DeleteVolumeResponse{}, nil
}

func (cs *controllerServer) ControllerExpandVolume(
	ctx context.Context,
	req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
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

	volOptions, err := newVolumeOptionsFromVolID(volumeId)
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
	sizeGiB, resizeRequired, err := expandVolume(ctx, curveVol, reqSizeGiB)
	if err != nil {
		ctxlog.ErrorS(ctx, err, "failed to expandVolume")
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &csi.ControllerExpandVolumeResponse{
		CapacityBytes:         int64(sizeGiB * volumehelpers.GiB),
		NodeExpansionRequired: resizeRequired,
	}, nil
}

// CreateSnapshot creates the snapshot in backend.
func (cs *controllerServer) CreateSnapshot(
	ctx context.Context,
	req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	if cs.snapshotServer == "" {
		return nil, status.Error(codes.Unimplemented, "")
	}
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

	// build source volume options from volume id
	sourceVolId := req.GetSourceVolumeId()
	volOptions, err := newVolumeOptionsFromVolID(sourceVolId)
	if err != nil {
		ctxlog.ErrorS(ctx, err, "failed to new volume options from id", "volumeId", sourceVolId)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	ctxlog.V(5).Infof(ctx, "build volOptions: %+v", volOptions)

	// lock out parallel delete/create/snapshot requests against the same volume
	if acquired := cs.volumeLocks.TryAcquire(volOptions.reqName); !acquired {
		ctxlog.Infof(ctx, util.VolumeOperationAlreadyExistsFmt, volOptions.reqName)
		return nil, status.Errorf(codes.Aborted, util.VolumeOperationAlreadyExistsFmt, volOptions.reqName)
	}
	defer cs.volumeLocks.Release(volOptions.reqName)

	snapServer := curveservice.NewSnapshotServer(cs.snapshotServer, volOptions.user, volOptions.volName)
	// verify the snapshot already exists
	curveSnapshot, err := snapServer.GetFileSnapshotOfName(ctx, snapshotName)
	if err == nil {
		ctxlog.V(4).Infof(ctx, "snapshot (name %v) already exists, check status...", snapshotName)
		return waitSnapshotDone(ctx, snapServer, curveSnapshot, sourceVolId)
	}
	if !util.IsNotFoundErr(err) {
		ctxlog.ErrorS(ctx, err, "failed to get snapshot by name", "snapshotName", snapshotName)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// check source volume status
	curveVol := curveservice.NewCurveVolume(volOptions.user, volOptions.volName, volOptions.sizeGiB)
	volDetail, err := curveVol.Stat(ctx)
	if err != nil {
		ctxlog.ErrorS(ctx, err, "failed to stat source volume", "volumeId", sourceVolId)
		return nil, status.Error(codes.Internal, err.Error())
	}
	if volDetail.FileStatus == curveservice.CurveVolumeStatusBeingCloned {
		ctxlog.Warningf(ctx, "the source volume %v status is BeingCloned, flatten it", sourceVolId)
		if err = snapServer.EnsureTaskFromSourceDone(ctx, volOptions.genVolumePath()); err != nil {
			ctxlog.ErrorS(ctx, err, "failed to flatten all tasks sourced volume", "volumeId", sourceVolId)
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	// do snapshot
	snapCurveUUID, err := snapServer.CreateSnapshot(ctx, snapshotName)
	if err != nil {
		ctxlog.ErrorS(ctx, err, "failed to create snapshot of name", "snapshotName", snapshotName)
		return nil, status.Error(codes.Internal, err.Error())
	}
	curveSnapshot, err = snapServer.GetFileSnapshotOfId(ctx, snapCurveUUID)
	if err != nil {
		ctxlog.ErrorS(ctx, err, "failed to get snapshot by id", "snapCurveUUID", snapCurveUUID)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return waitSnapshotDone(ctx, snapServer, curveSnapshot, sourceVolId)
}

// DeleteSnapshot deletes thesnapshot in backend.
func (cs *controllerServer) DeleteSnapshot(
	ctx context.Context,
	req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	if cs.snapshotServer == "" {
		return nil, status.Error(codes.Unimplemented, "")
	}
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

	snapCurveUUID, volOptions, err := parseSnapshotID(snapshotId)
	if err != nil {
		ctxlog.Warningf(ctx, "failed to parse snapshot id: %v", snapshotId)
		return &csi.DeleteSnapshotResponse{}, nil
	}

	snapServer := curveservice.NewSnapshotServer(cs.snapshotServer, volOptions.user, volOptions.volName)
	// get snapshot
	curveSnapshot, err := snapServer.GetFileSnapshotOfId(ctx, snapCurveUUID)
	if err != nil {
		if util.IsNotFoundErr(err, snapCurveUUID) {
			ctxlog.Infof(ctx, "snapshot %v not found, maybe already deleted.", snapshotId)
			return &csi.DeleteSnapshotResponse{}, nil
		}
		ctxlog.ErrorS(ctx, err, "failed to get snapshot", "snapCurveUUID", snapCurveUUID)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// lock out parallel snapshot
	if acquired := cs.snapshotLocks.TryAcquire(curveSnapshot.Name); !acquired {
		ctxlog.Errorf(ctx, util.SnapshotOperationAlreadyExistsFmt, curveSnapshot.Name)
		return nil, status.Errorf(codes.Aborted, util.SnapshotOperationAlreadyExistsFmt, curveSnapshot.Name)
	}
	defer cs.snapshotLocks.Release(curveSnapshot.Name)

	// ensure all the tasks created from this snapshot status done.
	if err = snapServer.EnsureTaskFromSourceDone(ctx, snapCurveUUID); err != nil {
		ctxlog.Errorf(ctx, "failed to ensure tasks from %v status done: %v", snapCurveUUID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// do delete
	if err = snapServer.DeleteSnapshot(ctx, snapCurveUUID); err != nil {
		ctxlog.ErrorS(ctx, err, "failed to delete snapshot", "snapCurveUUID", snapCurveUUID)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &csi.DeleteSnapshotResponse{}, nil
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

	volOptions, err := newVolumeOptionsFromVolID(volumeId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	curveVol := curveservice.NewCurveVolume(volOptions.user, volOptions.volName, volOptions.sizeGiB)
	volDetail, err := curveVol.Stat(ctx)
	if err != nil || volDetail.FileStatus == curveservice.CurveVolumeStatusNotExist {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &csi.ValidateVolumeCapabilitiesResponse{
		Confirmed: &csi.ValidateVolumeCapabilitiesResponse_Confirmed{
			VolumeCapabilities: req.VolumeCapabilities,
		},
	}, nil
}

// createVolFromContentSource creates a volume from the request contentSource
// return non-empty volSource if clone successfully.
func (cs *controllerServer) createVolFromContentSource(
	ctx context.Context,
	req *csi.CreateVolumeRequest,
	destVolOptions *volumeOptions,
	curveVol *curveservice.CurveVolume) (volSource string, err error) {
	if cs.snapshotServer == "" {
		return "", status.Error(codes.Unimplemented, "")
	}

	if req.VolumeContentSource == nil {
		return "", nil
	}

	volDestination := destVolOptions.genVolumePath()
	// check contentSource
	switch req.VolumeContentSource.Type.(type) {
	case *csi.VolumeContentSource_Snapshot:
		snapshotId := req.VolumeContentSource.GetSnapshot().GetSnapshotId()
		// lock out parallel snapshot
		if acquired := cs.snapshotLocks.TryAcquire(snapshotId); !acquired {
			ctxlog.Errorf(ctx, util.SnapshotOperationAlreadyExistsFmt, snapshotId)
			return "", status.Errorf(codes.Aborted, util.SnapshotOperationAlreadyExistsFmt, snapshotId)
		}
		defer cs.snapshotLocks.Release(snapshotId)
		// ensure the source snapshot exists,
		// and get the snapshot UUID as the source to create a new volume
		volSource, err = ensureSnapshotExists(ctx, cs.snapshotServer, snapshotId)
	case *csi.VolumeContentSource_Volume:
		volumeId := req.VolumeContentSource.GetVolume().GetVolumeId()
		// lock out parallel source volume
		if acquired := cs.volumeLocks.TryAcquire(volumeId); !acquired {
			ctxlog.Errorf(ctx, util.VolumeOperationAlreadyExistsFmt, volumeId)
			return "", status.Errorf(codes.Aborted, util.VolumeOperationAlreadyExistsFmt, volumeId)
		}
		defer cs.volumeLocks.Release(volumeId)
		// ensurce the source volume exists,
		// and get the volume path as the source to create a new volume
		volSource, err = ensureVolumeExists(ctx, cs.snapshotServer, volumeId)
	default:
		err = status.Errorf(codes.InvalidArgument, "not a proper volume source %v", req.VolumeContentSource)
	}
	if err != nil {
		return "", err
	}

	ctxlog.V(4).Infof(ctx, "clone/snapshot volume from %v to %v", volSource, volDestination)
	snapServer := curveservice.NewSnapshotServer(cs.snapshotServer, destVolOptions.user, destVolOptions.volName)
	var taskUUID string
	taskUUID, err = cloneVolume(ctx, snapServer, volSource, volDestination, destVolOptions.cloneLazy)
	if err != nil {
		ctxlog.ErrorS(ctx, err, "failed to clone volume")
		return "", status.Error(codes.Internal, err.Error())
	}
	ctxlog.V(4).Infof(ctx, "clone %v status done", taskUUID)

	// fix size if the cloned volume size less than requested size.
	_, _, err = expandVolume(ctx, curveVol, destVolOptions.sizeGiB)
	if err != nil {
		ctxlog.ErrorS(ctx, err, "failed to expand volume")
		return "", status.Error(codes.Internal, err.Error())
	}

	return volSource, nil
}

// Expand volume if the existing size is less than reqSizeGiB
func expandVolume(
	ctx context.Context,
	curveVol *curveservice.CurveVolume,
	reqSizeGiB int) (sizeGiB int, resizeRequired bool, err error) {
	volDetail, err := curveVol.Stat(ctx)
	if err != nil {
		return 0, false, err
	}
	if volDetail.FileStatus == curveservice.CurveVolumeStatusNotExist {
		return 0, false, fmt.Errorf("the curve volume not exists")
	}
	ctxlog.Infof(ctx, "volume %s(status %s) size is %dGiB, reqSize is round up to %dGiB",
		volDetail.FileName, volDetail.FileStatus, volDetail.LengthGiB, reqSizeGiB)
	if reqSizeGiB <= volDetail.LengthGiB {
		return volDetail.LengthGiB, false, nil
	}

	if err := curveVol.Extend(ctx, reqSizeGiB); err != nil {
		return 0, false, fmt.Errorf("failed to extend volume")
	}
	ctxlog.Infof(ctx, "successfully extend volume %s size to %dGiB", volDetail.FileName, reqSizeGiB)
	return reqSizeGiB, true, nil
}

// Waits the snapshot status Done.
// generate a snapshotId from the UUID in curve and the source volume id, then return the response.
func waitSnapshotDone(
	ctx context.Context,
	snapServer *curveservice.SnapshotServer,
	curveSnapshot curveservice.Snapshot,
	sourceVolId string) (*csi.CreateSnapshotResponse, error) {
	// wait snapshot status done
	if curveSnapshot.Status != curveservice.SnapshotStatusDone {
		var err error
		curveSnapshot, err = snapServer.WaitForSnapshotDone(ctx, curveSnapshot.UUID)
		if err != nil {
			ctxlog.ErrorS(ctx, err, "failed to wait snapshot status Done", "snapName", curveSnapshot.Name, "UUID", curveSnapshot.UUID)
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	snapshotId, err := composeSnapshotID(curveSnapshot.UUID, sourceVolId)
	if err != nil {
		ctxlog.ErrorS(ctx, err, "failed to compose snapshot id", "snapId", curveSnapshot.UUID, "sourceVolId", sourceVolId)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	createTime := time.Unix(0, curveSnapshot.Time*1000)
	ctxlog.Infof(ctx, "Snapshot(name %v csiId %v) status Done", curveSnapshot.Name, snapshotId)
	return &csi.CreateSnapshotResponse{
		Snapshot: &csi.Snapshot{
			SizeBytes:      int64(curveSnapshot.FileLength),
			SnapshotId:     snapshotId,
			SourceVolumeId: sourceVolId,
			CreationTime:   timestamppb.New(createTime),
			ReadyToUse:     true,
		},
	}, nil
}

// Clone volume from volSource to volDestination and wait for the clone task ready to use.
func cloneVolume(
	ctx context.Context,
	snapServer *curveservice.SnapshotServer,
	volSource, volDestination string,
	cloneLazy bool) (string, error) {
	taskInfo, err := snapServer.GetCloneTaskOfDestination(ctx, volDestination)
	if err == nil {
		ctxlog.V(4).Infof(ctx, "get existing task when clone: %v", taskInfo)
		if taskInfo.TaskStatus != curveservice.TaskStatusDone {
			return taskInfo.UUID, snapServer.WaitForCloneTaskReady(ctx, volDestination)
		}
		return taskInfo.UUID, nil
	}
	if !util.IsNotFoundErr(err) {
		return "", err
	}

	taskUUID, err := snapServer.Clone(ctx, volSource, volDestination, cloneLazy)
	if err != nil {
		return taskUUID, err
	}
	return taskUUID, snapServer.WaitForCloneTaskReady(ctx, volDestination)
}

// Ensure the snapshot exists.
func ensureSnapshotExists(ctx context.Context, snapshotServer, snapshotId string) (string, error) {
	snapCurveUUID, volOptions, err := parseSnapshotID(snapshotId)
	if err != nil {
		return "", status.Errorf(codes.NotFound, "snapshot id %v not found", snapshotId)
	}
	snapServer := curveservice.NewSnapshotServer(snapshotServer, volOptions.user, volOptions.volName)
	if _, err := snapServer.GetFileSnapshotOfId(ctx, snapCurveUUID); err != nil {
		if util.IsNotFoundErr(err, snapCurveUUID) {
			return "", status.Errorf(codes.NotFound, "the source snapshot(UUID %v) not found", snapCurveUUID)
		}
		return "", status.Error(codes.Internal, err.Error())
	}
	return snapCurveUUID, nil
}

// Ensure the volume exists.
// If the volume was cloned, ensure the clone task done.
func ensureVolumeExists(ctx context.Context, snapshotServer, volumeId string) (string, error) {
	volOptions, err := newVolumeOptionsFromVolID(volumeId)
	if err != nil {
		return "", status.Errorf(codes.NotFound, "volume id %v not found", volumeId)
	}
	curveVol := curveservice.NewCurveVolume(volOptions.user, volOptions.volName, volOptions.sizeGiB)
	if _, err := curveVol.Stat(ctx); err != nil {
		if util.IsNotFoundErr(err) {
			return "", status.Errorf(codes.NotFound, "the source volume (%v) not found", volOptions)
		}
		return "", status.Error(codes.Internal, err.Error())
	}
	// flatten the volume if it was cloned by other lazy
	snapServer := curveservice.NewSnapshotServer(snapshotServer, volOptions.user, volOptions.volName)
	volPath := volOptions.genVolumePath()
	taskInfo, err := snapServer.GetCloneTaskOfDestination(ctx, volPath)
	if err != nil {
		if util.IsNotFoundErr(err) {
			return volPath, nil
		}
		return "", status.Error(codes.Internal, err.Error())
	}
	if taskInfo.TaskStatus == curveservice.TaskStatusDone {
		return volPath, nil
	}
	if taskInfo.TaskStatus == curveservice.TaskStatusMetaInstalled {
		if err = snapServer.Flatten(ctx, taskInfo.UUID); err != nil {
			return "", status.Error(codes.Internal, err.Error())
		}
	}

	// wait done
	if err = snapServer.WaitForCloneTaskDone(ctx, volPath); err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}
	ctxlog.V(4).Infof(ctx, "the clone task of destination %v done", volPath)
	return volPath, nil
}
