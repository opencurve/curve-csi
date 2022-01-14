[中文版](../cn/curve-interface/curvebs-cli.md)

# Curve Block Storage CLI

- [Create a directory](#create-a-directory)
- [Delete a directory](#delete-a-directory)
- [Create a volume](#create-a-volume)
- [Extend a volume](#extend-a-volume)
- [Get the volume information](#get-the-volume-information)
- [List the volumes in a directory](#list-the-volumes-in-a-directory)
- [Delete a volume](#delete-a-volume)
- [Code Comparison](#code-comparison)

### Create a directory

`curve mkdir [-h] --user USER --dirname DIRNAME`

Args:
- USER: the user of the current DIRNAME
- DIRNAME: the absolute path, length must less than 4096 bytes

Return Code:
- OK: create successfully
- AUTHFAIL: authentication failed
- EXISTS: the DIRNAME already exists
- NOTEXISTS: the parent path of DIRNAME not exists
- INTERNAL_ERROR: other internal error

e.g.
 
```bash
$ curve mkdir --user k8s --dirname /k8s
```

### Delete a directory

`curve rmdir [-h] --user USER --dirname DIRNAME`

Args:
- USER: the user of the current DIRNAME
- DIRNAME: the absolute path, length must less than 4096 bytes

Return Code:
- OK: delete sucessfully
- AUTHFAIL: authentication failed
- NOTEXISTS: the DIRNAME not exists
- NOT_EMPTY: the directory not empty
- INTERNAL_ERROR: other internal error

e.g.

```bash
$ curve rmdir --user k8s --dirname /k8s
```

### Create a volume

`curve create [-h] --filename FILENAME --length LENGTH --user USER`

Args: 
- FILENAME: the absolute path contains directory name and volume name
- LENGTH: the unit is GiB and size limits to 10GiB~4TiB
- USER: the user of the directory

Return Code:
- Ok: create successfully
- AUTHFAIL: authentication failed
- EXISTS: the volume already exists
- NOTEXISTS: the directory of the volume not exists
- FAILED: other internal error

e.g.

```bash
$ curve create --filename /k8s/myvol --length 10 --user k8s
```

### Extend a volume

`curve extend [-h] --user USER --filename FILENAME --length LENGTH`

Args:
- USER: the user of the directory
- FILENAME: the absolute path contains directory name and volume name
- LENGTH: new size of this volume, the unit is GiB and size limits to 10GiB~4TiB

Return Code:
- Ok: extend successfully
- AUTHFAIL: authentication failed
- NOTEXISTS: the volume not exists
- NOT_SUPPORT: does not support extending this volume
- NO_SHRINK_BIGGER_FILE: the specific new size `LENGTH` less than origin size
- INTERNAL_ERROR: other internal error

e.g.

```bash
$ curve extend --filename /k8s/myvol --length 20 --user k8s
```

### Get the volume information

`curve stat [-h] --user USER --filename FILENAME`

Args: 
- USER: the user of the directory
- FILENAME: the absolute path contains directory name and volume name

Return Code:
- Ok: get successfully
- AUTHFAIL: authentication failed
- NOTEXISTS: the volume not exists
- INTERNAL_ERROR: other internal error

FileStatus Code:
- Created
- Deleting
- Cloning
- CloneMetaInstalled
- Cloned
- BeingCloned

e.g.

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

### List the volumes in a directory

`curve list [-h] --user USER --dirname DIRNAME`

Args:
- USER: the user of the directory
- DIRNAME: the absolute path, length must less than 4096 bytes

Return Code:
- OK: list successfully
- AUTHFAIL: authentication failed
- NOTEXISTS: the DIRNAME not exists
- INTERNAL_ERROR: other internal error

e.g.

```bash
$ curve list --user k8s --dirname /k8s
myvol
```

### Delete a volume

`curve delete [-h] --user USER --filename FILENAME`

Args:
- FILENAME: the absolute path contains directory name and volume name
- USER: the user of the directory

Return Code:
- Ok: delete successfully
- AUTHFAIL: authentication failed
- NOTEXISTS: the volume not exists
- FILE_OCCUPIED: the volume occupied by other processes
- INTERNAL_ERROR: other internal error

e.g.

```bash
$ curve delete --user k8s --filename /k8s/myvol
```

### Code Comparison

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
