[English version](../../curve-interface/curve-nbd.md)

# Curve nbd tool 

挂载/卸载cure设备。

- <a href="#usage">用法</a>
- <a href="#map">挂载</a>
- <a href="#list">列出本地已挂载设备</a>
- <a href="#unmap">卸载</a>

### <div id="usage">用法</div>

```
Usage: curve-nbd [options] map <image>  Map an image to nbd device
                 unmap <device|image>   Unmap nbd device
                 [options] list-mapped  List mapped nbd devices
Map options:
  --device <device path>  Specify nbd device path (/dev/nbd{num})
  --read-only             Map read-only
  --nbds_max <limit>      Override for module param nbds_max
  --max_part <limit>      Override for module param max_part
  --timeout <seconds>     Set nbd request timeout
  --try-netlink           Use the nbd netlink interface
```

### <div id="map">挂载</div>

设备名称格式为：`cbd:<user>/<filename_full_path>_<user>_`。

比如：

```
$ curve-nbd map cbd:k8s//k8s/csi-vol-volume-fa0c04c9-2e93-487e-8986-1e1625fd8c46_k8s_
```

### <div id="list">列出本地已挂载设备</div>

```
$ curve-nbd list-mapped
id      image                                                                device options
1509297 cbd:k8s//k8s/csi-vol-pvc-647525be-c0d6-464b-b548-1fa26f6d183c_k8s_ /dev/nbd1 timeout=86400
```

### <div id="unmap">卸载</div>

```
$ curve-nbd unmap /dev/nbd1
```
