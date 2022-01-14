[English version](../../curve-interface/curvebs-snapshot-clone-service.md)

# Curve块存储快照/克隆接口

- <a href="#createsnap">创建快照</a>
- <a href="#deletesnap">删除快照</a>
- <a href="#cancelsnap">取消快照</a>
- <a href="#getsnap">获取快照详情</a>
- <a href="#clone">克隆</a>
- [Flatten](#flatten)
- <a href="#gettask">查询指定任务的信息</a>
- <a href="#clean">清理克隆或恢复任务</a>
- <a href="#codecomp">编码对照表</a>

## <div id="createsnap">创建快照</div>

| 请求方法 | Url |
| --- | --- |
|GET | /SnapshotCloneService?Action=CreateSnapshot&Version=0.0.6&User=test&File=test&Name=test |

#### 请求数据

| 字段|类型|是否必须|描述|
| --- | --- | --- | --- |
|Action	|string | Y | CreateSnapshot |
|Version|string	| Y | API Version: 0.0.6|
|User	|string	| Y	| 卷所属用户 |
|File	|string	| Y | 快照的卷文件 |
|Name	|string	| Y | 快照名称 |

#### 响应数据

| 字段 | 类型 | 描述|
| --- | --- | --- |
|Code | string | 状态码 |
|Message | string| 额外信息 |
|RequestId | string	| 请求id|
|UUID |string | 快照uuid|

#### 举例

```
$ curl -XGET \
    'http://127.0.0.1:5555/SnapshotCloneService?Action=CreateSnapshot&Version=0.0.6&User=test&File=test&Name=test'
```

响应:

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

## <div id="deletesnap">删除快照</div>

| 请求方法 | Url |
| --- | --- |
|GET | /SnapshotCloneService?Action=DeleteSnapshot&Version=0.0.6&User=test&File=test&UUID=uuid1 |

#### 请求数据

| 字段|类型|是否必须|描述|
| --- | --- | --- | --- |
|Action	|string | Y | DeleteSnapshot |
|Version|string	| Y | API Version: 0.0.6|
|User	|string	| Y	| 卷所属用户 |
|File	|string	| Y | 快照的卷文件 |
|UUID	|string	| Y | 快照uuid |

#### 响应数据

| 字段 | 类型 | 描述|
| --- | --- | --- |
|Code | string | 状态码 |
|Message | string| 额外信息 |
|RequestId | string	| 请求id|


#### 举例

```
$ curl -XGET \
    'http://127.0.0.1:5555/SnapshotCloneService?Action=DeleteSnapshot&Version=0.0.6&User=test&File=test&UUID=uuid1'
```

响应:

```
HTTP/1.1 200 OK
Content-Length: xxx

{
    "Code" : "0",
    "Message" : "Exec success.",
    "RequestId" : "xxx"
}
```

## <div id="cancelsnap">取消快照</div>

| 请求方法 | Url |
| --- | --- |
|GET | /SnapshotCloneService?Action=CancelSnapshot&Version=0.0.6&User=test&File=test&UUID=uuid1 |

#### 请求数据

| 字段|类型|是否必须|描述|
| --- | --- | --- | --- |
|Action	|string | Y | CancelSnapshot |
|Version|string	| Y | API Version: 0.0.6|
|User	|string	| Y	| 卷所属用户 |
|File	|string	| Y | 快照的卷文件 |
|UUID	|string	| Y | 快照uuid |

#### 响应数据

| 字段 | 类型 | 描述|
| --- | --- | --- |
|Code | string | 状态码 |
|Message | string| 额外信息 |
|RequestId | string	| 请求id|

#### 举例

```
$ curl -XGET \
    'http://127.0.0.1:5555/SnapshotCloneService?Action=CancelSnapshot&Version=0.0.6&User=test&File=test&UUID=uuid1'
```

响应:

```
HTTP/1.1 200 OK
Content-Length: xxx

{
    "Code" : "0",
    "Message" : "Exec success.",
    "RequestId" : "xxx"
}
```

## <div id="getsnap">获取快照详情</div>

| 请求方法 | Url |
| --- | --- |
|GET | /SnapshotCloneService?Action=GetFileSnapshotInfo&Version=0.0.6&User=test&File=test&Limit=10&Offset=0 |
|GET | /SnapshotCloneService?Action=GetFileSnapshotInfo&Version=0.0.6&User=test&File=test&UUID=de06df66-b9e4-44df-ba3d-ac94ddee0b28 |

#### 请求数据

| 字段|类型|是否必须|描述|
| --- | --- | --- | --- |
|Action	|string | Y | GetFileSnapinfo |
|Version|string	| Y | API Version: 0.0.6|
|User	|string	| Y	| 卷所属用户 |
|File	|string	| Y | 快照的卷文件 |
|Limit  |int    | N | 快照数目，默认是10|
|Offset |int    | N | 查询offset, 默认为0|
|UUID	|string	| N | 快照uuid |

#### 响应数据

| 字段 | 类型 | 描述|
| --- | --- | --- |
|Code | string | 状态码 |
|Message | string| 额外信息 |
|RequestId | string	| 请求id|
|TotalCount | string | 快照数量|
|Snapshots | Snapshot | 快照详情列表 |

Snapshot结构:

| 字段 | 类型| 描述 |
| --- | --- | --- |
|UUID	|string	| 快照uuid |
|User	|string	| 卷所属用户 |
|File	|string	| 做快照的源卷文件
|SeqNum	|uint32	| 快照版本号 |
|Name	|string	| 快照名称 |
|Time	|uint64	| 打快照的时间戳 |
|FileLength	|uint32	| 快照大小，单位Byte |
|Status	|enum| 快照状态: <br/> (0:done, 1:pending,2:deleteing, 3:errorDeleting, 4:canceling, 5:error）|
|Progress|uint32 |	打快照的进度百分比 |

#### 举例

```
$ curl -XGET \
    'http://127.0.0.1:5555/SnapshotCloneService?Action=GetFileSnapshotInfo&Version=0.0.6&User=zjm&File=/zjm/test1&Limit=10'
```

响应:

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

## <div id="clone">克隆</div>

| 请求方法 | Url |
| --- | --- |
|GET | /SnapshotCloneService?Action=Clone&Version=0.0.6&User=zjm&Source=/zjm/test1&Destination=/zjm/clone1&Lazy=true |

#### 请求数据

| 字段|类型|是否必须|描述|
| --- | --- | --- | --- |
|Action	|string | Y | Clone |
|Version|string	| Y | API Version: 0.0.6|
|User	|string	| Y	| 卷所属用户 |
|Source	|string	| Y | 卷文件名称，或者快照uuid |
|Destination	|string	| Y |目标卷文件名称 |
|Lazy | bool | Y | lazy克隆 |

#### 响应数据

| 字段 | 类型 | 描述|
| --- | --- | --- |
|Code | string | 状态码 |
|Message | string| 额外信息 |
|RequestId | string	| 请求id|
|UUID | string | 克隆任务UUID |

#### 举例

```
$ curl -XGET \
    'http://127.0.0.1:5555/SnapshotCloneService?Action=Clone&Version=0.0.6&User=zjm&Source=/zjm/test1&Destination=/zjm/clone1&Lazy=true'
```

响应:

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

## <div id="recover">从快照恢复卷</div>

| 请求方法 | Url |
| --- | --- |
|GET | /SnapshotCloneService?Action=Recover&Version=0.0.6&User=zjm&Source=de06df66-b9e4-44df-ba3d-ac94ddee0b28&Destination=/zjm/recover1&Lazy=true |

#### 请求数据

| 字段|类型|是否必须|描述|
| --- | --- | --- | --- |
|Action	|string | Y | Recover |
|Version|string	| Y | API Version: 0.0.6|
|User	|string	| Y	| 卷所属用户 |
|Source	|string	| Y | 快照uuid |
|Destination	|string	| Y | 恢复的卷文件名称 |
|Lazy | bool | Y | lazy恢复 |

#### 响应数据

| 字段 | 类型 | 描述|
| --- | --- | --- |
|Code | string | 状态码 |
|Message | string| 额外信息 |
|RequestId | string	| 请求id|
|UUID | string | 恢复任务UUID |

#### 举例

```
$ curl -XGET \
    'http://127.0.0.1:5555/SnapshotCloneService?Action=Recover&Version=0.0.6&User=zjm&Source=de06df66-b9e4-44df-ba3d-ac94ddee0b28&Destination=/zjm/recover1&Lazy=true'
```

响应:

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

| 请求方法 | Url |
| --- | --- |
|GET | /SnapshotCloneService?Action=Flatten&Version=0.0.6&User=zjm&UUID=xxx |

####  请求数据

| 字段|类型|是否必须|描述|
| --- | --- | --- | --- |
|Action	|string | Y | Flatten |
|Version|string	| Y | API Version: 0.0.6|
|User	|string	| Y	| 卷所属用户 |
|UUID   |string | Y | 任务UUID |

#### 响应数据

| 字段 | 类型 | 描述|
| --- | --- | --- |
|Code | string | 状态码 |
|Message | string| 额外信息 |
|RequestId | string	| 请求id|

#### 举例

```
$ curl -XGET \
    'http://127.0.0.1:5555/SnapshotCloneService?Action=Flatten&Version=0.0.6&User=zjm&UUID=de06df66-b9e4-44df-ba3d-ac94ddee0b28'
```

响应:

```
HTTP/1.1 200 OK
Content-Length: xxx

{
    "Code" : "0",
    "Message" : "Exec success.",
    "RequestId" : "xxx"
}
```

## <div id="gettask">查询指定任务的信息</div>

获取指定用户的所有任务，可以限定：

- 指定任务UUID
- 指定卷文件

| 请求方法 | Url |
| --- | --- |
| GET |	/SnapshotCloneService?Action=GetCloneTasks&Version=0.0.6&User=zjm&Limit=10&Offset=0 |
| GET | /SnapshotCloneService?Action=GetCloneTasks&Version=0.0.6&User=zjm&UUID=xxx |
| GET | /SnapshotCloneService?Action=GetCloneTasks&Version=0.0.6&User=zjm&File=DestFileName |

#### 请求数据

| 字段|类型|是否必须|描述|
| --- | --- | --- | --- |
|Action	|string | Y | GetCloneTasks |
|Version|string	| Y | API Version: 0.0.6|
|User	|string	| Y	| 卷所属用户 |
|Limit	|int | N | 返回的最大任务数，默认是10 |
|Offset	|int | N | 查询offset, 默认为0|
|UUID	|string | N | 克隆/恢复任务的 UUID|
|File	|string | N | 克隆或恢复任务所属的卷文件名 |

#### 响应数据

| 字段 | 类型 | 描述|
| --- | --- | --- |
|Code | string | 状态码 |
|Message | string| 额外信息 |
|RequestId | string	| 请求id|
|TotalCount	| int| 返回的任务数目|
|TaskInfos	| TaskInfo | 任务详情列表|

TaskInfo结构:

| 字段 | 类型 | 描述 |
| --- | --- | --- |
|TaskType	|enum| 任务类型: <br/>（0:clone, 1:recover）|
|FileType   |enum| 文件类型: <br/> (0:SrcFile 1:SrcSnapshot) |
|IsLazy | bool |任务是否lazy|
|Progress|uint8| 任务进度百分比 |
|Src | string | 源文件，可以是卷文件或快照uuid |
|User	|string	| 卷所属用户 |
|File	|string| 克隆/恢复任务的目标文件名称 |
|Time	|uint64| 创建时间戳 |
|TaskStatus	|enum|任务状态: <br/>（0:done, 1:cloning, 2:recovering, 3:cleaning, 4:errorCleaning, 5:error，6:retrying, 7:metaInstalled）|

#### 举例

```
$ curl -XGET \
    'http://127.0.0.1:5555//SnapshotCloneService?Action=GetCloneTasks&Version=0.0.6&User=zjm&Limit=10'
```

响应:

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

## <div id="clean">清理克隆或恢复任务</div>

清理克隆/恢复任务。如果这个任务之前失败了，则删除克隆的临时文件，否则只删除任务。

| 请求方法 | Url |
| --- | --- |
| GET |	/SnapshotCloneService?Action=CleanCloneTask&Version=0.0.6&User=zjm&UUID=78e83875-2b50-438f-8f25-36715380f4f5 |

####  请求数据

| 字段|类型|是否必须|描述|
| --- | --- | --- | --- |
|Action	|string | Y | CleanCloneTask |
|Version|string	| Y | API Version: 0.0.6|
|User	|string	| Y	| 卷所属用户 |
|UUID	|string | Y | 克隆/恢复任务UUID|

#### 响应数据

| 字段 | 类型 | 描述|
| --- | --- | --- |
|Code | string | 状态码 |
|Message | string| 额外信息 |
|RequestId | string	| 请求id|

#### 举例

```
$ curl -XGET \
    'http://127.0.0.1:5555/SnapshotCloneService?Action=CleanCloneTask&Version=0.0.6&User=zjm&UUID=78e83875-2b50-438f-8f25-36715380f4f5'
```

响应:

```
HTTP/1.1 200 OK
Content-Length: xxx

{
    "Code" : "0",
    "Message" : "Exec success.",
    "RequestId" : "xxx"
}
```

## <div id="codecomp">编码对照表</div>

|Code |Message |	HTTP Status Code| Description |
| --- | --- | --- | --- |
|0	| Exec success.|	200| 执行成功|
|-1	|Internal error.|	500| 未知错误。此状态不应出现，否则应联系管理员处理|
|-2	|Server init fail.|	500	| 此状态不应出现。init阶段失败。|
|-3 |Server start fail.|500	| 此状态不应出现。start阶段失败。|
|-4	|Service is stop|	500	| 服务端在终止中。|
|-5	|BadRequest:"Invalid request."	|400| 请求确实必要的字段或字段非法。|
|-6	|Task already exist.|	500	|目前该状态不会出现|
|-7	|Invalid user.|	500	| 用户名和请求的文件/快照/任务等不匹配。|
|-8	|File not exist.|	500	| 快照时源文件不存在。 <br/>获取快照时快照不存在。<br/>克隆/恢复时快照或源文件不存在。<br/>获取或清理任务时任务不存在。|
|-9|File status invalid.|	500	| 对文件打快照时，文件处于克隆或恢复中。<br/>克隆时源文件或快照处于克隆或恢复中。|
|-10|Chunk size not aligned.	|500| 配置等chunk大小没校准，一般此情况不出现。|
|-11|FileName not match.	|500| 克隆或取消快照任务时，文件不匹配快照。|
|-12|Cannot delete unfinished.	|500| 不能删除未完成的快照。|
|-13|Cannot create when has error.|	500| 当以下情况时不能快照/克隆/恢复： <br/> - 源文件有错误 <br/> - 已存在的快照有错误。 <br/> - 存在的任务有错误。|
|-14|Cannot cancel finished.	|500| 取消快照时，快照还未完成或不存在。|
|-15|Invalid snapshot.	|500| 快照处于任务中或有错误时，不能克隆或恢复它。|
|-16|Cannot delete when using.	|500| 快照处于克隆或恢复中时，不能删除它。|
|-17|Cannot clean task unfinished.	|500	| 不能清除未完成的克隆或恢复任务。 |
|-18|Snapshot count reach the limit.|	500	| 快照数目触发到了最大值。|
|-19|File exist.	|500	| 克隆或快照时，目标文件已存在。|
|-20|Task is full.	|500	| 克隆/恢复任务触发到了最大值。|
