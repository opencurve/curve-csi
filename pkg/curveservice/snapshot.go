/*
Copyright 2021 The Netease Kubernetes Authors.

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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/opencurve/curve-csi/pkg/util"
	"github.com/opencurve/curve-csi/pkg/util/ctxlog"
)

const (
	// The following three values are used for 30 seconds timeout
	// while waiting for snapshot Watcher to expire.
	snapshotWatcherInitDelay = 1 * time.Second
	snapshotWatcherFactor    = 1.4
	snapshotWatcherSteps     = 10
)

type SnapshotServer struct {
	URL      string `json:"server"`
	User     string `json:"user"`
	FilePath string `json:"filepath"`
}

func NewSnapshotServer(server, user, volName string) *SnapshotServer {
	return &SnapshotServer{
		URL:      server + "/SnapshotCloneService",
		User:     user,
		FilePath: "/" + user + "/" + volName,
	}
}

// GetSnapshotByName gets the snapshot with specific name
func (cs *SnapshotServer) GetFileSnapshotOfName(ctx context.Context, snapName string) (Snapshot, error) {
	var snap Snapshot
	limit, offset, total := 20, 0, 1
	for offset < total {
		snapshotResp, err := cs.getFileSnapshots(ctx, "", limit, offset)
		if err != nil {
			return snap, err
		}
		for _, oneSnap := range snapshotResp.Snapshots {
			if oneSnap.Name == snapName {
				ctxlog.V(4).Infof(ctx, "[curve snapshot] get snapshot (%+v) by name: %v", oneSnap, snapName)
				return oneSnap, nil
			}
		}

		total = snapshotResp.TotalCount
		offset += limit
	}

	return snap, util.NewNotFoundErr()
}

// GetSnapshotById gets the snapshot with specific uuid
func (cs *SnapshotServer) GetFileSnapshotOfId(ctx context.Context, uuid string) (Snapshot, error) {
	var snap Snapshot
	snapshotResp, err := cs.getFileSnapshots(ctx, uuid, 0, 0)
	if err != nil {
		ctxlog.ErrorS(ctx, err, "failed to get snapshot by id", "uuid", uuid)
		return snap, err
	}

	// check snapshot field
	if len(snapshotResp.Snapshots) > 1 {
		return snap, fmt.Errorf("found multi snapshots with uuid %v", uuid)
	}
	if snapshotResp.Snapshots[0].UUID != uuid {
		return snap, fmt.Errorf("the snapshot(%+v) in response not matched with uuid: %v", snapshotResp.Snapshots[0], uuid)
	}

	ctxlog.V(4).Infof(ctx, "[curve snapshot] get snapshot %v successfully.", snapshotResp.Snapshots[0])
	return snapshotResp.Snapshots[0], nil
}

// getFileSnapshots get snapshots list
func (cs *SnapshotServer) getFileSnapshots(ctx context.Context, uuid string, limit, offset int) (GetSnapshotResp, error) {
	var resp GetSnapshotResp
	queryMap := map[string]string{
		"Action":  "GetFileSnapshotInfo",
		"Version": "0.0.6",
		"User":    cs.User,
		"File":    cs.FilePath,
	}
	if limit > 0 {
		queryMap["Limit"] = strconv.Itoa(limit)
	}
	if offset > 0 {
		queryMap["Offset"] = strconv.Itoa(offset)
	}
	if uuid != "" {
		queryMap["UUID"] = uuid
	}

	ctxlog.V(4).Infof(ctx, "starting to get snapshots: %v", queryMap)
	statusCode, data, err := util.HttpGet(cs.URL, queryMap)
	if err != nil {
		return resp, fmt.Errorf("failed to get snapshot, err: %v", err)
	}

	if err = json.Unmarshal(data, &resp); err != nil {
		return resp, fmt.Errorf("unmarshal failed when get snapshot. statusCode: %v, data: %v, err: %v", statusCode, string(data), err)
	}

	if resp.Code == FileNotExists || len(resp.Snapshots) == 0 {
		ctxlog.V(4).Infof(ctx, "not found, resp: %+v", resp)
		if uuid != "" {
			return resp, util.NewNotFoundErr(uuid)
		}
		return resp, util.NewNotFoundErr()
	}

	if resp.Code != ExecSuccess {
		return resp, fmt.Errorf("faied to get snapshot, resp: %+v", resp)
	}

	ctxlog.V(5).Infof(ctx, "[curve snapshot] get snapshots: %+v", resp)
	return resp, nil
}

// CreateSnapshot creates a snapshot and returns the uuid
func (cs *SnapshotServer) CreateSnapshot(ctx context.Context, snapName string) (string, error) {
	queryMap := map[string]string{
		"Action":  "CreateSnapshot",
		"Version": "0.0.6",
		"User":    cs.User,
		"File":    cs.FilePath,
		"Name":    snapName,
	}

	ctxlog.V(4).Infof(ctx, "starting to create snapshot: %v", queryMap)
	statusCode, data, err := util.HttpGet(cs.URL, queryMap)
	if err != nil {
		return "", fmt.Errorf("failed to create snapshot, err: %v", err)
	}
	if statusCode != http.StatusOK {
		return "", fmt.Errorf("failed to create snapshot, statusCode: %v, response data: %v", statusCode, string(data))
	}

	var resp CreateSnapshotResp
	if err = json.Unmarshal(data, &resp); err != nil {
		return "", fmt.Errorf("failed to unmarshal data: %v, err: %v", string(data), err)
	}

	if resp.Code != ExecSuccess {
		return "", fmt.Errorf("faied to create snapshot, resp: %+v", resp)
	}

	ctxlog.V(4).Infof(ctx, "[curve snapshot] create snapshot successfully with uuid: %v", resp.UUID)
	return resp.UUID, nil
}

// DeleteSnapshot detetes a snapshot
func (cs *SnapshotServer) DeleteSnapshot(ctx context.Context, uuid string) error {
	queryMap := map[string]string{
		"Action":  "DeleteSnapshot",
		"Version": "0.0.6",
		"User":    cs.User,
		"File":    cs.FilePath,
		"UUID":    uuid,
	}

	ctxlog.V(4).Infof(ctx, "starting to delete snapshot: %v", queryMap)
	statusCode, data, err := util.HttpGet(cs.URL, queryMap)
	if err != nil {
		return fmt.Errorf("failed to delete snapshot, err: %v", err)
	}
	if statusCode != http.StatusOK {
		return fmt.Errorf("failed to delete snapshot, statusCode: %v, response data: %v", statusCode, string(data))
	}

	var resp DeleteSnapshotResp
	if err = json.Unmarshal(data, &resp); err != nil {
		return fmt.Errorf("failed to unmarshal data: %v, err: %v", string(data), err)
	}

	if resp.Code != ExecSuccess {
		return fmt.Errorf("faied to delete snapshot, resp: %+v", resp)
	}

	ctxlog.V(4).Infof(ctx, "[curve snapshot] delete snapshot successfully with uuid: %v", uuid)
	return nil
}

// CancelSnapshot cancels a snapshot
func (cs *SnapshotServer) CancelSnapshot(ctx context.Context, uuid string) error {
	queryMap := map[string]string{
		"Action":  "CancelSnapshot",
		"Version": "0.0.6",
		"User":    cs.User,
		"File":    cs.FilePath,
		"UUID":    uuid,
	}

	ctxlog.V(4).Infof(ctx, "starting to cancel snapshot: %v", queryMap)
	statusCode, data, err := util.HttpGet(cs.URL, queryMap)
	if err != nil {
		return fmt.Errorf("failed to cancel snapshot, err: %v", err)
	}

	var resp CancelSnapshotResp
	if err = json.Unmarshal(data, &resp); err != nil {
		return fmt.Errorf("unmarshal failed when cancel snapshot. statusCode: %v, data: %v, err: %v", statusCode, string(data), err)
	}

	if resp.Code == FileNotExists {
		return util.NewNotFoundErr(uuid)
	}
	if resp.Code != ExecSuccess {
		return fmt.Errorf("faied to cancel snapshot, resp: %+v", resp)
	}

	ctxlog.V(4).Infof(ctx, "[curve snapshot] cancel snapshot successfully with uuid: %v", uuid)
	return nil
}

// Wait for the snapshot ready
func (cs *SnapshotServer) WaitForSnapshotDone(ctx context.Context, uuid string) (Snapshot, error) {
	var (
		snap Snapshot
		err  error
	)
	backoff := wait.Backoff{
		Duration: snapshotWatcherInitDelay,
		Factor:   snapshotWatcherFactor,
		Steps:    snapshotWatcherSteps,
	}

	waitErr := wait.ExponentialBackoff(backoff, func() (bool, error) {
		snap, err = cs.GetFileSnapshotOfId(ctx, uuid)
		if err != nil {
			return false, fmt.Errorf("failed to get snapshort for uuid %v, err: %v", uuid, err)
		}
		ctxlog.V(4).Infof(ctx, "the snapshot (name: %v uuid: %v) process %v%%", snap.Name, uuid, snap.Progress)
		return snap.Status == SnapshotStatusDone, nil
	})
	// return error if err has not become available for the specified timeout
	if waitErr == wait.ErrWaitTimeout {
		return snap, fmt.Errorf("timeout to wait, snapshot (uuid %v) is still not done", uuid)
	}
	// return error if any other errors were encountered during waiting for the snapshot to become done
	return snap, waitErr
}

// Get task with specific destination
func (cs *SnapshotServer) GetCloneTaskOfDestination(ctx context.Context, destination string) (TaskInfo, error) {
	var taskInfo TaskInfo
	resp, err := cs.getCloneTask(ctx, "", destination, 0, 0)
	if err != nil {
		return taskInfo, err
	}
	if len(resp.TaskInfos) > 1 {
		return taskInfo, fmt.Errorf("found multi task of destination: %v", destination)
	}
	ctxlog.V(4).Infof(ctx, "[curve snapshot] get clone task successfully: %+v", resp.TaskInfos[0])
	return resp.TaskInfos[0], nil
}

// Get task with specific uuid
func (cs *SnapshotServer) GetCloneTaskOfId(ctx context.Context, uuid string) (TaskInfo, error) {
	var taskInfo TaskInfo
	resp, err := cs.getCloneTask(ctx, uuid, "", 0, 0)
	if err != nil {
		return taskInfo, err
	}
	if len(resp.TaskInfos) > 1 {
		return taskInfo, fmt.Errorf("found multi task of uuid: %v", uuid)
	}
	ctxlog.V(4).Infof(ctx, "[curve snapshot] get clone task successfully: %+v", resp.TaskInfos[0])
	return resp.TaskInfos[0], nil
}

// Get tasks
func (cs *SnapshotServer) getCloneTask(ctx context.Context, uuid, destination string, limit, offset int) (GetCloneTaskResp, error) {
	var resp GetCloneTaskResp
	queryMap := map[string]string{
		"Action":  "GetCloneTasks",
		"Version": "0.0.6",
		"User":    cs.User,
	}
	if uuid != "" {
		queryMap["UUID"] = uuid
	}
	if destination != "" {
		queryMap["File"] = destination
	}
	if limit != 0 {
		queryMap["Limit"] = strconv.Itoa(limit)
	}
	if offset != 0 {
		queryMap["Offset"] = strconv.Itoa(offset)
	}

	ctxlog.V(4).Infof(ctx, "starting to get clone task: %v", queryMap)
	statusCode, data, err := util.HttpGet(cs.URL, queryMap)
	if err != nil {
		return resp, fmt.Errorf("failed to get clone task, err: %v", err)
	}

	if err = json.Unmarshal(data, &resp); err != nil {
		return resp, fmt.Errorf("unmarshal failed when get task. statusCode: %v, data: %v, err: %v", statusCode, string(data), err)
	}

	if resp.Code == FileNotExists || len(resp.TaskInfos) == 0 {
		return resp, util.NewNotFoundErr()
	}
	if resp.Code != ExecSuccess {
		return resp, fmt.Errorf("faied to get task, resp: %+v", resp)
	}
	return resp, nil
}

// Clone a volume from source to destination
func (cs *SnapshotServer) Clone(ctx context.Context, source, destination string, lazy bool) (string, error) {
	queryMap := map[string]string{
		"Action":      "Clone",
		"Version":     "0.0.6",
		"User":        cs.User,
		"Source":      source,
		"Destination": destination,
		"Lazy":        strconv.FormatBool(lazy),
	}

	ctxlog.V(4).Infof(ctx, "starting to clone snapshot: %v", queryMap)
	statusCode, data, err := util.HttpGet(cs.URL, queryMap)
	if err != nil {
		return "", fmt.Errorf("failed to clone snapshot, err: %v", err)
	}

	var resp CloneResp
	if err = json.Unmarshal(data, &resp); err != nil {
		return "", fmt.Errorf("unmarshal failed when clone. statusCode: %v, data: %v, err: %v", statusCode, string(data), err)
	}

	if resp.Code == FileNotExists {
		return "", util.NewNotFoundErr()
	}
	if resp.Code != ExecSuccess {
		return "", fmt.Errorf("faied to clone snapshot, resp: %+v", resp)
	}

	ctxlog.V(4).Infof(ctx, "[curve snapshot] clone %v to %v successfully with task id: %v", source, destination, resp.UUID)
	return resp.UUID, nil
}

// Clean a clone task, flatten if it is unfinished.
func (cs *SnapshotServer) CleanCloneTask(ctx context.Context, uuid string) error {
	ctxlog.V(4).Infof(ctx, "get task status %v before clean it", uuid)
	taskInfo, err := cs.GetCloneTaskOfId(ctx, uuid)
	if err != nil {
		if util.IsNotFoundErr(err) {
			return nil
		}
		return err
	}
	if taskInfo.TaskStatus == TaskStatusMetaInstalled {
		if err = cs.Flatten(ctx, uuid); err != nil {
			return err
		}
	}
	if err = cs.waitForCloneTaskStatus(ctx, taskInfo.File, TaskStatusDone, TaskStatusError); err != nil {
		if util.IsNotFoundErr(err) {
			return nil
		}
		return err
	}

	return cs.cleanCloneTask(ctx, uuid)
}

func (cs *SnapshotServer) cleanCloneTask(ctx context.Context, uuid string) error {
	queryMap := map[string]string{
		"Action":  "CleanCloneTask",
		"Version": "0.0.6",
		"User":    cs.User,
		"UUID":    uuid,
	}

	ctxlog.V(4).Infof(ctx, "starting to clean cloneTask: %v", queryMap)
	statusCode, data, err := util.HttpGet(cs.URL, queryMap)
	if err != nil {
		return fmt.Errorf("failed to clean cloneTask, err: %v", err)
	}

	var resp CleanCloneTaskResp
	if err = json.Unmarshal(data, &resp); err != nil {
		return fmt.Errorf("unmarshal failed when clean cloneTask. statusCode: %v, data: %v, err: %v", statusCode, string(data), err)
	}

	if resp.Code == FileNotExists {
		return nil
	}
	if resp.Code != ExecSuccess {
		return fmt.Errorf("faied to clean cloneTask, resp: %+v", resp)
	}

	ctxlog.V(4).Infof(ctx, "[curve snapshot] successfully clean cloneTask: %v", uuid)
	return nil
}

func (cs *SnapshotServer) Flatten(ctx context.Context, uuid string) error {
	queryMap := map[string]string{
		"Action":  "Flatten",
		"Version": "0.0.6",
		"User":    cs.User,
		"UUID":    uuid,
	}

	ctxlog.V(4).Infof(ctx, "starting to flatten task: %v", queryMap)
	statusCode, data, err := util.HttpGet(cs.URL, queryMap)
	if err != nil {
		return fmt.Errorf("failed to flatten task, err: %v", err)
	}

	var resp FlattenResp
	if err = json.Unmarshal(data, &resp); err != nil {
		return fmt.Errorf("unmarshal failed when flatten task. statusCode: %v, data: %v, err: %v", statusCode, string(data), err)
	}

	if resp.Code == FileNotExists {
		return util.NewNotFoundErr()
	}
	if resp.Code != ExecSuccess {
		return fmt.Errorf("faied to flatten task, resp: %+v", resp)
	}

	ctxlog.V(4).Infof(ctx, "[curve snapshot] successfully flatten task: %v", uuid)
	return nil
}

func (cs *SnapshotServer) waitForCloneTaskStatus(ctx context.Context, destination string, taskStatus ...TaskStatus) error {
	if len(taskStatus) == 0 {
		return nil
	}

	backoff := wait.Backoff{
		Duration: snapshotWatcherInitDelay,
		Factor:   snapshotWatcherFactor,
		Steps:    snapshotWatcherSteps,
	}

	waitErr := wait.ExponentialBackoff(backoff, func() (bool, error) {
		taskInfo, err := cs.GetCloneTaskOfDestination(ctx, destination)
		if err != nil {
			ctxlog.ErrorS(ctx, err, "failed to get clone task", "destination", destination)
			return false, err
		}
		for _, status := range taskStatus {
			if taskInfo.TaskStatus == status {
				return true, nil
			}
		}
		return false, nil
	})
	// return error if err has not become available for the specified timeout
	if waitErr == wait.ErrWaitTimeout {
		return fmt.Errorf("timeout to wait, task is still not %v", taskStatus)
	}
	// return error if any other errors were encountered during waiting for the task to become done
	return waitErr
}

// Wait for the task ready: Done or MetaInstalled
func (cs *SnapshotServer) WaitForCloneTaskReady(ctx context.Context, destination string) error {
	ctxlog.V(4).Infof(ctx, "wait for task of destination %q status ready", destination)
	return cs.waitForCloneTaskStatus(ctx, destination, TaskStatusDone, TaskStatusMetaInstalled)
}

func (cs *SnapshotServer) WaitForCloneTaskDone(ctx context.Context, destination string) error {
	ctxlog.V(4).Infof(ctx, "wait for task of destination %q status done", destination)
	return cs.waitForCloneTaskStatus(ctx, destination, TaskStatusDone)
}

func (cs *SnapshotServer) EnsureTaskFromSourceDone(ctx context.Context, source string) error {
	ctxlog.V(4).Infof(ctx, "ensure task created from %v status done", source)

	needFlatten := make([]TaskInfo, 0)
	limit, offset, total := 20, 0, 1
	for offset < total {
		taskResp, err := cs.getCloneTask(ctx, "", "", limit, offset)
		if err != nil {
			if util.IsNotFoundErr(err) {
				break
			}
			return err
		}
		for _, oneTask := range taskResp.TaskInfos {
			if oneTask.Src == source && oneTask.TaskStatus == TaskStatusMetaInstalled {
				needFlatten = append(needFlatten, oneTask)
			}
		}
		total = taskResp.TotalCount
		offset += limit
	}

	ctxlog.V(4).Infof(ctx, "need flatten tasks: %v", needFlatten)
	for _, t := range needFlatten {
		if err := cs.Flatten(ctx, t.UUID); err != nil {
			return err
		}
	}
	for _, t := range needFlatten {
		taskInfo, err := cs.GetCloneTaskOfId(ctx, t.UUID)
		if err != nil {
			if util.IsNotFoundErr(err) {
				continue
			}
			return err
		}
		if taskInfo.TaskStatus == TaskStatusError {
			ctxlog.Warningf(ctx, "%v status err, just clean it", t)
			if err = cs.cleanCloneTask(ctx, taskInfo.UUID); err != nil {
				return err
			}
		}
		if err := cs.waitForCloneTaskStatus(ctx, t.File, TaskStatusDone); err != nil {
			return err
		}
	}
	ctxlog.V(4).Infof(ctx, "[curve snapshot] all tasks from %v done", source)
	return nil
}
