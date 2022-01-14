[中文版](../cn/curve-interface/curvebs-snapshot-clone-service.md)

# Curve Block Storage Snapshot/Clone Interface

- [Create a snapshot](#create-a-snapshot)
- [Delete a snapshot](#delete-a-snapshot)
- [Cancel a snapshot](#cancel-a-snapshot)
- [Query the snapshot information](#query-the-snapshot-information)
- [Clone](#clone)
- [Volume recover from a snapshot](#volume-recover-from-a-snapshot)
- [Flatten](#flatten)
- [Query the information of the specific clone or recover task](#query-the-information-of-the-specific-clone-or-recover-task)
- [Clean a clone or recover task](#clean-a-clone-or-recover-task)
- [Code Comparison](#code-comparison)


## Create a snapshot

| Method | Url |
| --- | --- |
|GET | /SnapshotCloneService?Action=CreateSnapshot&Version=0.0.6&User=test&File=test&Name=test |

#### Query Data

| Name|Type|Require|Description|
| --- | --- | --- | --- |
|Action	|string | Y | CreateSnapshot |
|Version|string	| Y | API Version: 0.0.6|
|User	|string	| Y	| the user of the volume |
|File	|string	| Y | the volume which to snapshot |
|Name	|string	| Y | the name of snapshot |

#### Response Data

| Name | Type | Description|
| --- | --- | --- |
|Code | string | status code |
|Message | string| extra message |
|RequestId | string	| request id|
|UUID |string | the uuid of the snapshot|

#### e.g.

```
$ curl -XGET \
    'http://127.0.0.1:5555/SnapshotCloneService?Action=CreateSnapshot&Version=0.0.6&User=test&File=test&Name=test'
```

Response:

```
HTTP/1.1 200 OK
Content-Length: xxx

{
    "Code" : "0",
    "Message" : "Exec success.",
    "RequestId" : "xxx",
    "UUID" : "xxx"
}
```

## Delete a snapshot

| Method | Url |
| --- | --- |
|GET | /SnapshotCloneService?Action=DeleteSnapshot&Version=0.0.6&User=test&File=test&UUID=uuid1 |

#### Query Data

| Name|Type|Require|Description|
| --- | --- | --- | --- |
|Action	|string | Y | DeleteSnapshot |
|Version|string	| Y | API Version: 0.0.6|
|User	|string	| Y	| the user of the volume |
|File	|string	| Y | the volume which to snapshot |
|UUID	|string	| Y | the uuid of snapshot |

#### Response Data

| Name | Type | Description|
| --- | --- | --- |
|Code | string | status code |
|Message | string| extra message |
|RequestId | string	| request id|

#### e.g.

```
$ curl -XGET \
    'http://127.0.0.1:5555/SnapshotCloneService?Action=DeleteSnapshot&Version=0.0.6&User=test&File=test&UUID=uuid1'
```

Response:

```
HTTP/1.1 200 OK
Content-Length: xxx

{
    "Code" : "0",
    "Message" : "Exec success.",
    "RequestId" : "xxx"
}
```

## Cancel a snapshot

| Method | Url |
| --- | --- |
|GET | /SnapshotCloneService?Action=CancelSnapshot&Version=0.0.6&User=test&File=test&UUID=uuid1 |

#### Query Data

| Name|Type|Require|Description|
| --- | --- | --- | --- |
|Action	|string | Y | CancelSnapshot |
|Version|string	| Y | API Version: 0.0.6|
|User	|string	| Y	| the user of the volume |
|File	|string	| Y | the volume which to snapshot |
|UUID	|string	| Y | the uuid of snapshot |

#### Response Data

| Name | Type | Description|
| --- | --- | --- |
|Code | string | status code |
|Message | string| extra message |
|RequestId | string	| request id|

#### e.g.

```
$ curl -XGET \
    'http://127.0.0.1:5555/SnapshotCloneService?Action=CancelSnapshot&Version=0.0.6&User=test&File=test&UUID=uuid1'
```

Response:

```
HTTP/1.1 200 OK
Content-Length: xxx

{
    "Code" : "0",
    "Message" : "Exec success.",
    "RequestId" : "xxx"
}
```

## Query the snapshot information

| Method | Url |
| --- | --- |
|GET | /SnapshotCloneService?Action=GetFileSnapshotInfo&Version=0.0.6&User=test&File=test&Limit=10&Offset=0 |
|GET | /SnapshotCloneService?Action=GetFileSnapshotInfo&Version=0.0.6&User=test&File=test&UUID=de06df66-b9e4-44df-ba3d-ac94ddee0b28 |

#### Query Data

| Name|Type|Require|Description|
| --- | --- | --- | --- |
|Action	|string | Y | GetFileSnapinfo |
|Version|string	| Y | API Version: 0.0.6|
|User	|string	| Y	| the user of the volume |
|File	|string	| Y | the volume which to snapshot |
|Limit  |int    | N | the maximum snapshots, default to 10|
|Offset |int    | N | the query offset, default to 0|
|UUID	|string	| N | the uuid of snapshot |

#### Response Data

| Name | Type | Description |
| --- | --- | --- |
|Code | string | status code |
|Message | string| extra message |
|RequestId | string	| request id|
|TotalCount | string | the total count of the volume snapshots|
|Snapshots | Snapshot | the list of the snapshot information |

Snapshot struct:

| Name | Type| Description |
| --- | --- | --- |
|UUID	|string	| the uuid of snapshot |
|User	|string	| the user of the volume |
|File	|string	| the source volume which snapshotted |
|SeqNum	|uint32	| snapshot version number |
|Name	|string	| snapshot name |
|Time	|uint64	| created time stamp |
|FileLength	|uint32	|the size of the volume (unit Byte)|
|Status	|enum| snapshot status: <br/> (0:done, 1:pending,2:deleteing, 3:errorDeleting, 4:canceling, 5:error）|
|Progress|uint32 |	the percent of snapshot progress |

#### e.g.

```
$ curl -XGET \
    'http://127.0.0.1:5555/SnapshotCloneService?Action=GetFileSnapshotInfo&Version=0.0.6&User=zjm&File=/zjm/test1&Limit=10'
```

Response:

```
HTTP/1.1 200 OK
Content-Length: xxx

{
    "Code" : "0",
    "Message" : "Exec success.",
    "RequestId" : "xxx",
    "TotalCount": 1,
    "Snapshots":[{
        "File" : "/zjm/test1",
        "FileLength" : 10737418240,
        "Name" : "snap1",
        "Progress" : 30,
        "SeqNum" : 1,
        "Status" : 1,
        "Time" : 1564391913582677,
        "UUID" : "de06df66-b9e4-44df-ba3d-ac94ddee0b28",
        "User" : "zjm"
    }]
}
```

## Clone

| Method | Url |
| --- | --- |
|GET | /SnapshotCloneService?Action=Clone&Version=0.0.6&User=zjm&Source=/zjm/test1&Destination=/zjm/clone1&Lazy=true |

#### Query Data

| Name|Type|Require|Description|
| --- | --- | --- | --- |
|Action	|string | Y | Clone |
|Version|string	| Y | API Version: 0.0.6|
|User	|string	| Y	| the user of the volume |
|Source	|string	| Y | the volume name, or snapshot uuid |
|Destination	|string	| Y | cloned volume name |
|Lazy | bool | Y | lazy clone |

#### Response Data

| Name | Type | Description |
| --- | --- | --- |
|Code | string | status code |
|Message | string| extra message |
|RequestId | string	| request id|
|UUID | string | the clone task UUID |

#### e.g.

```
$ curl -XGET \
    'http://127.0.0.1:5555/SnapshotCloneService?Action=Clone&Version=0.0.6&User=zjm&Source=/zjm/test1&Destination=/zjm/clone1&Lazy=true'
```

Response:

```
HTTP/1.1 200 OK
Content-Length: xxx

{
    "Code" : "0",
    "Message" : "Exec success.",
    "RequestId" : "xxx",
    "UUID" : "xxx"
}
```

## Volume recover from a snapshot

| Method | Url |
| --- | --- |
|GET | /SnapshotCloneService?Action=Recover&Version=0.0.6&User=zjm&Source=de06df66-b9e4-44df-ba3d-ac94ddee0b28&Destination=/zjm/recover1&Lazy=true |

#### Query Data

| Name|Type|Require|Description|
| --- | --- | --- | --- |
|Action	|string | Y | Recover |
|Version|string	| Y | API Version: 0.0.6|
|User	|string	| Y	| the user of the volume |
|Source	|string	| Y | snapshot uuid |
|Destination	|string	| Y | recovered volume name |
|Lazy | bool | Y | lazy recover |

#### Response Data

| Name | Type | Description |
| --- | --- | --- |
|Code | string | status code |
|Message | string| extra message |
|RequestId | string	| request id|
|UUID | string | the recover task UUID |

#### e.g.

```
$ curl -XGET \
    'http://127.0.0.1:5555/SnapshotCloneService?Action=Recover&Version=0.0.6&User=zjm&Source=de06df66-b9e4-44df-ba3d-ac94ddee0b28&Destination=/zjm/recover1&Lazy=true'
```

Response:

```
HTTP/1.1 200 OK
Content-Length: xxx

{
    "Code" : "0",
    "Message" : "Exec success.",
    "RequestId" : "xxx",
    "UUID" : "xxx"
}
```

## Flatten

| Method | Url |
| --- | --- |
|GET | /SnapshotCloneService?Action=Flatten&Version=0.0.6&User=zjm&UUID=xxx |

#### Query Data

| Name|Type|Require|Description|
| --- | --- | --- | --- |
|Action	|string | Y | Flatten |
|Version|string	| Y | API Version: 0.0.6|
|User	|string	| Y	| the user of the volume |
|UUID   |string | Y | task UUID |

#### Response Data

| Name | Type | Description |
| --- | --- | --- |
|Code | string | status code |
|Message | string| extra message |
|RequestId | string	| request id|

#### e.g.

```
$ curl -XGET \
    'http://127.0.0.1:5555/SnapshotCloneService?Action=Flatten&Version=0.0.6&User=zjm&UUID=de06df66-b9e4-44df-ba3d-ac94ddee0b28'
```

Response:

```
HTTP/1.1 200 OK
Content-Length: xxx

{
    "Code" : "0",
    "Message" : "Exec success.",
    "RequestId" : "xxx"
}
```

## Query the information of the specific clone or recover task

Get all tasks of the specific user, you can limit with:

- specific UUID
- specific File

| Method | Url |
| --- | --- |
| GET |	/SnapshotCloneService?Action=GetCloneTasks&Version=0.0.6&User=zjm&Limit=10&Offset=0 |
| GET | /SnapshotCloneService?Action=GetCloneTasks&Version=0.0.6&User=zjm&UUID=xxx |
| GET | /SnapshotCloneService?Action=GetCloneTasks&Version=0.0.6&User=zjm&File=DestFileName |

#### Query Data

| Name|Type|Require|Description|
| --- | --- | --- | --- |
|Action	|string | Y | GetCloneTasks |
|Version|string	| Y | API Version: 0.0.6|
|User	|string	| Y	| the user of the volume |
|Limit	|int | N | the maximum tasks, default to 10|
|Offset	|int | N | query offset, default to 0|
|UUID	|string | N | UUID of the clone/recover task|
|File	|string | N | volume name which the clone/recover task belongs to|

#### Response Data

| Name | Type | Description |
| --- | --- | --- |
|Code | string | status code |
|Message | string| extra message |
|RequestId | string	| request id|
|TotalCount	| int|the total count of tasks|
|TaskInfos	| TaskInfo | the list of tasks information|

TaskInfo struct:

| Name | Type | Description |
| --- | --- | --- |
|TaskType	|enum| task type: <br/>（0:clone, 1:recover）|
|FileType   |enum| task file type: <br/> (0:SrcFile 1:SrcSnapshot) |
|IsLazy | bool |task is lazy|
|Progress|uint8| task progress percent |
|Src | string | task source snapshot or file |
|User	|string	|the user of the volume |
|File	|string|volume name which the clone/recover task belongs to|
|Time	|uint64| created time stamp |
|TaskStatus	|enum|	task status: <br/>（0:done, 1:cloning, 2:recovering, 3:cleaning, 4:errorCleaning, 5:error，6:retrying, 7:metaInstalled）|

#### e.g.

```
$ curl -XGET \
    'http://127.0.0.1:5555//SnapshotCloneService?Action=GetCloneTasks&Version=0.0.6&User=zjm&Limit=10'
```

Response:

```
HTTP/1.1 200 OK
Content-Length: xxx

{
    "Code": "0",
    "Message": "Exec success.",
    "RequestId": "xxx",
    "TotalCount": 1,
    "TaskInfos":[{
        "File" : "/zjm/clone1",
        "UUID" : "78e83875-2b50-438f-8f25-36715380f4f5",
        "TaskStatus" : 5,
        "TaskType" : 0,
        "Time" : 0,
        "User" : "zjm"
 }]
}
```

## Clean a clone or recover task

Clean the clone/recover task. Delete the temporary cloned file in curvefs server if the task has been failed, otherwise delete the task only.

| Method | Url |
| --- | --- |
| GET |	/SnapshotCloneService?Action=CleanCloneTask&Version=0.0.6&User=zjm&UUID=78e83875-2b50-438f-8f25-36715380f4f5 |

#### Query Data

| Name|Type|Require|Description|
| --- | --- | --- | --- |
|Action	|string | Y | CleanCloneTask |
|Version|string	| Y | API Version: 0.0.6|
|User	|string	| Y	| the user of the volume |
|UUID	|string | Y | UUID of the clone/recover task|

#### Response Data

| Name | Type | Description |
| --- | --- | --- |
|Code | string | status code |
|Message | string| extra message |
|RequestId | string	| request id|

#### e.g.

```
$ curl -XGET \
    'http://127.0.0.1:5555/SnapshotCloneService?Action=CleanCloneTask&Version=0.0.6&User=zjm&UUID=78e83875-2b50-438f-8f25-36715380f4f5'
```

Response:

```
HTTP/1.1 200 OK
Content-Length: xxx

{
    "Code" : "0",
    "Message" : "Exec success.",
    "RequestId" : "xxx"
}
```

## Code Comparison

|Code |Message |	HTTP Status Code| Description |
| --- | --- | --- | --- |
|0	| Exec success.|	200| Execute successfully.|
|-1	|Internal error.|	500| Unknown internal error. This case shoud not happen. Connect with the administrator if this error is returned|
|-2	|Server init fail.|	500	| This case shoud not happen. Failed at the init phase.|
|-3 |Server start fail.|500	| This case shoud not happen. Failed at the start phase.|
|-4	|Service is stop|	500	| Get this code when the server is terminating.|
|-5	|BadRequest:"Invalid request."	|400| Request lack of necessary fields or with invalid value.|
|-6	|Task already exist.|	500	| Now this case will never happed.|
|-7	|Invalid user.|	500	| Mismatch between the User field in the request and the File/Image/Snapshot to be operated.|
|-8	|File not exist.|	500	| Source volume not exists when snapshot. <br/>Snapshot not exists when get the information. <br/>Snapshot/Image not exists when clone/recover. <br/>Task not exists when clean/get the task.|
|-9|File status invalid.|	500	| The file is cloning/recovering instead of Normal status when task snapshots. <br/>The source image is cloning/recovering instead of Normal status when clone it.|
|-10|Chunk size not aligned.	|500| The configured chunk size not aligned. In general, it should not happen.|
|-11|FileName not match.	|500| The file is not matched with the snapshot, when clean/cancel the snapshot task|
|-12|Cannot delete unfinished.	|500|	Can not delete unfinished snapshot.|
|-13|Cannot create when has error.|	500| Can not snapshot/clone/recover when: <br/> - the file with error. <br/> - the existing snapshot with error. <br/> - the existing task with error.|
|-14|Cannot cancel finished.	|500| The snapshot has been finished or not exists when cancel it.|
|-15|Invalid snapshot.	|500| Can not clone/recover a snapshot which is processing or has error.|
|-16|Cannot delete when using.	|500| Can not delete the snapshot which is cloning/recovering.|
|-17|Cannot clean task unfinished.	|500	| Can not clean unfinished clone/recover tasks. |
|-18|Snapshot count reach the limit.|	500	| The count of snapshots reaches the max limit.|
|-19|File exist.	|500	| Target file exists when snapshot/clone.|
|-20|Task is full.	|500	| The count of clone/recover tasks reaches the max limit.|
