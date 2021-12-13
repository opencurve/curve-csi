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

type SnapshotStatus uint8

type RespCode string

const (
	ExecSuccess   RespCode = "0"
	FileNotExists RespCode = "-8"
)

// SnapshotStatus:
//  （0:done, 1:pending, 2:deleteing, 3:errorDeleting, 4:canceling, 5:error）
const (
	SnapshotStatusDone SnapshotStatus = iota
	SnapshotStatusPending
	SnapshotStatusDeleteing
	SnapshotStatusErrorDeleting
	SnapshotStatusCanceling
	SnapshotStatusError
)

type TaskType uint8

// TaskType: (0:clone, 1:recover）
const (
	TaskTypeClone TaskType = iota
	TaskTypeRecover
)

type TaskFileType uint8

// TaskFileType: (0:SrcFile 1:SrcSnapshot)
const (
	TaskFileTypeSrcFile TaskFileType = iota
	TaskFileTypeSrcSnapshot
)

type TaskStatus uint8

// TaskStatus
//  （0:done, 1:cloning, 2:recovering, 3:cleaning, 4:errorCleaning, 5:error，6:retrying, 7:metaInstalled）
const (
	TaskStatusDone TaskStatus = iota
	TaskStatusCloning
	TaskStatusRecovering
	TaskStatusCleaning
	TaskStatusErrorCleaning
	TaskStatusError
	TaskStatusRetrying
	TaskStatusMetaInstalled
)

type SnapshotCommonResp struct {
	Code      RespCode `json:"Code"`
	Message   string   `json:"Message"`
	RequestId string   `json:"RequestId"`
}

type CreateSnapshotResp struct {
	SnapshotCommonResp
	UUID string `json:"UUID,omitempty"`
}

type DeleteSnapshotResp SnapshotCommonResp

type CancelSnapshotResp SnapshotCommonResp

type GetSnapshotResp struct {
	SnapshotCommonResp
	TotalCount int        `json:"TotalCount,omitempty"`
	Snapshots  []Snapshot `json:"Snapshots,omitempty"`
}

type Snapshot struct {
	UUID       string         `json:"UUID"`
	User       string         `json:"User"`
	File       string         `json:"File"`
	SeqNum     uint32         `json:"SeqNum"`
	Name       string         `json:"Name"`
	Time       int64          `json:"Time"`
	FileLength uint64         `json:"FileLength"` //unit Byte
	Status     SnapshotStatus `json:"Status"`
	Progress   uint8          `json:"Progress"`
}

type CloneResp struct {
	SnapshotCommonResp
	UUID string `json:"UUID,omitempty"`
}

type RecoverResp struct {
	SnapshotCommonResp
	UUID string `json:"UUID,omitempty"`
}

type FlattenResp SnapshotCommonResp

// also use as recover task
type GetCloneTaskResp struct {
	SnapshotCommonResp
	TotalCount int        `json:"TotalCount,omitempty"`
	TaskInfos  []TaskInfo `json:"TaskInfos,omitempty"`
}

type TaskInfo struct {
	File         string       `json:"File"`
	TaskFileType TaskFileType `json:"FileType,omitempty"`
	IsLazy       bool         `json:"IsLazy,omitempty"`
	Progress     uint8        `json:"Progress,omitempty"`
	Src          string       `json:"Src,omitempty"`
	TaskStatus   TaskStatus   `json:"TaskStatus"`
	TaskType     TaskType     `json:"TaskType"`
	Time         int64        `json:"Time"`
	UUID         string       `json:"UUID"`
	User         string       `json:"User"`
}

type CleanCloneTaskResp SnapshotCommonResp
