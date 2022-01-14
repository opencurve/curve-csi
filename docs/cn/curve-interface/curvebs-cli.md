[English version](../../curve-interface/curvebs-cli.md)

# Curve块存储命令行

- <a href="#createdir">创建目录</a>
- <a href="#rmdir">删除目录</a>
- <a href="#createvol">创建卷文件</a>
- <a href="#extendvol">卷扩容</a>
- <a href="#statvol">获取卷文件详情</a>
- <a href="#listdir">列出指定目录的所有卷文件</a>
- <a href="#deletevol">删除卷文件</a>
- <a href="#codecomp">编码对照表</a>


### <div id="createdir">创建目录</div>

`curve mkdir [-h] --user USER --dirname DIRNAME`

参数:
- USER: 目录的用户名，以此字段做租户隔离
- DIRNAME: 目录的绝对路径，且字段长度要小于4096字节

返回码:
- OK: 创建成功
- AUTHFAIL: 认证失败
- EXISTS: 目录已存在
- NOTEXISTS: 目录的父目录不存在
- INTERNAL_ERROR: 其它内部错误

举例：
 
```bash
$ curve mkdir --user k8s --dirname /k8s
```

### <div id="rmdir">删除目录</div>

`curve rmdir [-h] --user USER --dirname DIRNAME`

参数:
- USER: 目录的用户名，以此字段做租户隔离
- DIRNAME: 目录的绝对路径，且字段长度要小于4096字节

返回码:
- OK: 删除成功
- AUTHFAIL: 认证失败
- NOTEXISTS: 目录不存在
- NOT_EMPTY: 目录非空
- INTERNAL_ERROR: 其它内部错误

举例：

```bash
$ curve rmdir --user k8s --dirname /k8s
```

### <div id="createvol">创建卷文件</div>

`curve create [-h] --filename FILENAME --length LENGTH --user USER`

参数:
- FILENAME: 包含目录及文件名的绝对路径
- LENGTH: 卷大小，单位为GiB，且范围在10GiB~4TiB
- USER: 目录的所属用户

返回码:
- Ok: 创建成功
- AUTHFAIL: 认证失败
- EXISTS: 卷已存在
- NOTEXISTS: 目录不存在
- FAILED: 其它内部错误

举例：

```bash
$ curve create --filename /k8s/myvol --length 10 --user k8s
```

### <div id="extendvol">卷扩容</div>

`curve extend [-h] --user USER --filename FILENAME --length LENGTH`

参数:
- USER: 目录的所属用户
- FILENAME: 包含目录及文件名的绝对路径
- LENGTH: 新的大小, 单位为GiB，且范围在10GiB~4TiB

返回码:
- Ok: 扩容成功
- AUTHFAIL: 认证失败
- NOTEXISTS: 卷不存在
- NOT_SUPPORT: 此卷不支持扩容
- NO_SHRINK_BIGGER_FILE: 新指定的容量小于目前的容量
- INTERNAL_ERROR: 其它内部错误

举例：

```bash
$ curve extend --filename /k8s/myvol --length 20 --user k8s
```

### <div id="statvol">获取卷文件详情</div>

`curve stat [-h] --user USER --filename FILENAME`

参数:
- USER: 目录的所属用户
- FILENAME: 包含目录及文件名的绝对路径

返回码:
- Ok: 获取成功
- AUTHFAIL: 认证失败
- NOTEXISTS: 卷不存在
- INTERNAL_ERROR: 其它内部错误

卷文件状态码:
- Created
- Deleting
- Cloning
- CloneMetaInstalled
- Cloned
- BeingCloned

举例：

```bash
$ curve stat --user k8s --filename /k8s/myvol
id: 40004
parentid: 39005
filetype: INODE_PAGEFILE
length(GB): 10
createtime: 2020-08-24 19:06:35
user: k8s
filename: myvol
fileStatus: Created
```

### <div id="listdir">列出指定目录的所有卷文件</div>

`curve list [-h] --user USER --dirname DIRNAME`

参数:
- USER: 目录的所属用户
- DIRNAME: 目录的绝对路径，且字段长度要小于4096字节

返回码:
- OK: 列出成功
- AUTHFAIL: 认证失败
- NOTEXISTS: 目录不存在
- INTERNAL_ERROR: 其它内部错误

举例：

```bash
$ curve list --user k8s --dirname /k8s
myvol
```

### <div id="deletevol">删除卷文件</div>

`curve delete [-h] --user USER --filename FILENAME`

参数:
- FILENAME: 包含目录及文件名的绝对路径
- USER: 目录的所属用户

返回码:
- Ok: 删除成功
- AUTHFAIL: 认证失败
- NOTEXISTS: 卷不存在
- FILE_OCCUPIED: 卷被其它进程占用
- INTERNAL_ERROR: 其它内部错误

举例：

```bash
$ curve delete --user k8s --filename /k8s/myvol
```

### <div id="codecomp">编码对照表</div>

```text
LIBCURVE_ERROR {
    OK                          = 0,
    EXISTS                      = 1,
    FAILED                      = 2,
    DISABLEIO                   = 3,
    AUTHFAIL                    = 4,
    DELETING                    = 5,
    NOTEXIST                    = 6,
    UNDER_SNAPSHOT              = 7,
    NOT_UNDERSNAPSHOT           = 8,
    DELETE_ERROR                = 9,
    NOT_ALLOCATE                = 10,
    NOT_SUPPORT                 = 11,
    NOT_EMPTY                   = 12,
    NO_SHRINK_BIGGER_FILE       = 13,
    SESSION_NOTEXISTS           = 14,
    FILE_OCCUPIED               = 15,
    PARAM_ERROR                 = 16,
    INTERNAL_ERROR              = 17,
    CRC_ERROR                   = 18,
    INVALID_REQUEST             = 19,
    DISK_FAIL                   = 20,
    NO_SPACE                    = 21,
    NOT_ALIGNED                 = 22,
    BAD_FD                      = 23,
    LENGTH_NOT_SUPPORT          = 24,
    SESSION_NOT_EXIST           = 25,
    STATUS_NOT_MATCH            = 26,
    DELETE_BEING_CLONED         = 27,
    CLIENT_NOT_SUPPORT_SNAPSHOT = 28,
    SNAPSTHO_FROZEN             = 29,
    UNKNOWN                     = 100
};

FileStatus {
    Created            = 0,
    Deleting           = 1,
    Cloning            = 2,
    CloneMetaInstalled = 3,
    Cloned             = 4,
    BeingCloned        = 5
};
```
