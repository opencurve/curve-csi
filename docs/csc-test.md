
# Test Using CSC Tool

- [Prerequisite](#prerequisite)
- [CurveBS](#curvebs)
- [CurveFS](#curvefs)

## Prerequisite

**Install csc tool:**

Get the csc code from `https://github.com/rexray/gocsi/tree/release/1.1/csc` and run `make`.

## CurveBS

#### Start curve csi driver

```bash
curve-csi --nodeid testnode \
    --type curvebs \
    --endpoint tcp://127.0.0.1:10000  \
    --snapshot-server http://127.0.0.1:5556 \
    -v 5
```

#### Get plugin info

```text
$ csc identity plugin-info --endpoint tcp://127.0.0.1:10000
"curve.csi.netease.com" "csi-v1.1.0-rc1"
```

#### Get supported capabilities

```text
$ csc identity plugin-capabilities --endpoint tcp://127.0.0.1:10000
CONTROLLER_SERVICE
ONLINE
```

#### Get controller implemented capabilities

```text
$ csc controller get-capabilities  --endpoint tcp://127.0.0.1:10000
&{type:CREATE_DELETE_VOLUME }
&{type:CLONE_VOLUME }
&{type:EXPAND_VOLUME }
```

#### Create a volume

```text
$ uuidgen
fa0c04c9-2e93-487e-8986-1e1625fd8c46

$ csc controller create --endpoint tcp://127.0.0.1:10000 \
    --req-bytes 10737418240 \
    --cap 1,mount,ext4 \
    --params user=k8s \
    volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
"0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46"  10737418240     "user"="k8s"
```

If the volume is block type, set: `--cap 5,1`

Check:

```text
$ sudo curve stat --user k8s --filename /k8s/csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
id: 39013
parentid: 39005
filetype: INODE_PAGEFILE
length(GB): 10
createtime: 2020-08-24 10:43:32
user: k8s
filename: csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
fileStatus: Created
root@pubt1-k8s-for-dev2:~#
```

#### NodeStage a volume

```text
$ sudo mkdir -p /mnt/test-csi/volume-globalmount
$ csc node stage --endpoint tcp://127.0.0.1:10000 \
   --cap 1,mount,ext4 \
   --staging-target-path /mnt/test-csi/volume-globalmount \
   0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
```

If the volume is block type, set: `--cap 5,1`

Check:

```text
$ sudo curve-nbd list-mapped
id    image                                                      device
97297 cbd:k8s//k8s/pvc-ce482926-91d8-11ea-bf6e-fa163e23ce53_k8s_ /dev/nbd0

$ sudo findmnt /mnt/test-csi/volume-globalmount/0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
TARGET                          SOURCE    FSTYPE OPTIONS
/mnt/test-csi/volume-globalmount/0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46 /dev/nbd0 ext4   rw,relatime,data=ordered
```

#### NodePublish a volume

```text
$ csc node publish --endpoint tcp://127.0.0.1:10000 \
    --target-path /mnt/test-csi/test-pod \
    --cap 1,mount,ext4 \
    --staging-target-path /mnt/test-csi/volume-globalmount \
    0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
```

Check:

```text
$ sudo findmnt /mnt/test-csi/test-pod
TARGET                 SOURCE    FSTYPE OPTIONS
/mnt/test-csi/test-pod /dev/nbd0 ext4   rw,relatime,data=ordered
```

#### NodeGetVolumeStats a volume

```text
$ csc node stats --endpoint tcp://127.0.0.1:10000 \
    0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46:/mnt/test-csi/test-pod
0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46    /mnt/test-csi/test-pod       10447220736     10501771264     37773312      BYTES
655349  655360  11      INODES
```

Check:

```
$ sudo df /mnt/test-csi/test-pod
Filesystem     1K-blocks  Used Available Use% Mounted on
/dev/nbd0      10255636 36888  10202364   1% /mnt/test-csi/test-pod

$ sudo df -i /mnt/test-csi/test-pod
Filesystem     Inodes IUsed  IFree IUse% Mounted on
/dev/nbd0      655360    11 655349    1% /mnt/test-csi/test-pod
```

#### Expand a Volume

```text
$ # controllerExpand:
$ csc controller expand --endpoint tcp://127.0.0.1:10000 \
    --req-bytes 21474836480 \
    0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
21474836480

$ # nodeExpand:
$ csc node expand --endpoint tcp://127.0.0.1:10000 \
    0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46 \
    /mnt/test-csi/test-pod
0
```

Check:

```text
$ sudo curve stat --user k8s --filename /k8s/csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46|grep length
length(GB): 20

$ sudo df -h /mnt/test-csi/test-pod
/dev/nbd0        20G   44M   20G   1% /mnt/test-csi/test-pod
```

#### NodeUnpublish a volume

```text
$ csc node unpublish --endpoint tcp://127.0.0.1:10000 \
    --target-path /mnt/test-csi/test-pod \
    0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
```

Check:

```
$ sudo file /mnt/test-csi/test-pod
/mnt/test-csi/test-pod: cannot open `/mnt/test-csi/test-pod' (No such file or directory)
```

#### NodeUnstage a volume

```text
$ csc node unstage --endpoint tcp://127.0.0.1:10000 \
    --staging-target-path /mnt/test-csi/volume-globalmount \
    0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
```

Check:

```
$ sudo file /mnt/test-csi/volume-globalmount
/mnt/test-csi/volume-globalmount: cannot open `/mnt/test-csi/volume-globalmount' (No such file or directory)

$ sudo curve-nbd list-mapped
```

#### Delete a volume

```text
$ csc controller delete --endpoint tcp://127.0.0.1:10000 \
  0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
```

Check:

```
$ sudo curve stat --user k8s --filename /k8s/csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
E 2020-08-24T13:57:55.946636+0800 58360 mds_client.cpp:395] GetFileInfo: filename = /k8s/volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46, owner = k8s, errocde = 6, error msg = kFileNotExists, log id = 1
stat fail, ret = -6
```

#### Snapshot

```
## create a snapshot
$ csc controller create-snapshot --endpoint tcp://127.0.0.1:10000 \
	--source-volume 0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46 \
	snapshot-215d24ff-c04c-4b08-a1fb-692c94627c63
"0024-9ea1a8fc-160d-47ef-b2ef-f0e09677b066-0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46"	23622320128	0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46	seconds:1639297566 nanos:780828000 	true

## delete a snapshot
$ csc controller delete-snapshot --endpoint tcp://127.0.0.1:10000 \
    0024-9ea1a8fc-160d-47ef-b2ef-f0e09677b066-0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
0024-9ea1a8fc-160d-47ef-b2ef-f0e09677b066-0003-k8s-csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46
```

## CurveBS

