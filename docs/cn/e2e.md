[English version](../e2e.md)

# E2E测试

- <a href="#req">依赖</a>
- [CurveBS](#curvebs)
  - <a href="#bsstart">启动curvebs插件</a>
  - <a href="#bse2e">运行e2e测试</a>
- [CurveFS](#curvefs)


## <div id="req">依赖</div>

**安装csi-sanity**

启动csi-test v4.3.0版本，你可以通过命令`go get github.com/kubernetes-csi/csi-test/cmd/csi-sanity`获取csi-sanity，二进制在 `$GOPATH/bin/csi-sanity`。

`go get`命令总是获取最新版本，如果想使用特定版本，可以获取[源码](https://github.com/kubernetes-csi/csi-test/releases)，然后运行`make -C cmd/csi-sanity`。 二进制将产生在`cmd/csi-sanity/csi-sanity`。

## CurveBS

### <div id="bsstart">启动curvebs插件</div>

```
curve-csi --nodeid testnode \
    --type curvebs \
    --endpoint tcp://127.0.0.1:10000  \
    --snapshot-server http://127.0.0.1:5556 \
    -v 5
```

### <div id="bse2e">运行e2e测试</div>

跳过2个测试场景:

- 此插件的限制`len(user+volume)<=80`，所以不能满足csi卷名称的最大128长度限制。
- 插件没有实现：`should fail when requesting to create a snapshot with already existing name and different source volume ID`

```
cat > config.yaml <<EOF
user: k8s
cloneLazy: "true"
EOF

skip[0]="should not fail when creating volume with maximum-length name"
skip[1]="should fail when requesting to create a snapshot with already existing name and different source volume ID"

csi-sanity -csi.endpoint dns:///127.0.0.1:10000     \
-csi.testvolumeparameters ./config.yaml    \
 -ginkgo.skip "$(IFS=\| ; echo "${skip[*]}")"
```

结果如下:

```
Ran 54 of 78 Specs in 85.288 seconds
SUCCESS! -- 54 Passed | 0 Failed | 1 Pending | 23 Skipped
```

## CurveFS

