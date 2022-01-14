[中文版](cn/e2e.md)

# E2E

- [Prerequisite](#prerequisite)
- [CurveBS](#curvebs)
  - [Start the curvebs csi plugin](#start-the-curvebs-csi-plugin)
  - [Run e2e test](#run-e2e-test)
- [CurveFS](#curvefs)


## Prerequisite

**Install csi-sanity**

Starting with csi-test v4.3.0, you can build the csi-sanity command with `go get github.com/kubernetes-csi/csi-test/cmd/csi-sanity` and you'll find the compiled binary in `$GOPATH/bin/csi-sanity`.

go get always builds the latest revision from the master branch. To build a certain release, [get the source code](https://github.com/kubernetes-csi/csi-test/releases) and run `make -C cmd/csi-sanity`. This produces `cmd/csi-sanity/csi-sanity`.

## CurveBS

### Start the curvebs csi plugin

```
curve-csi --nodeid testnode \
    --type curvebs \
    --endpoint tcp://127.0.0.1:10000  \
    --snapshot-server http://127.0.0.1:5556 \
    -v 5
```

### Run e2e test

Skip 2 cases:

- As the plugin limit: `len(user+volume)<=80`, so can not support volume with maximum-length(128).
- The plugin does not implement: `should fail when requesting to create a snapshot with already existing name and different source volume ID`

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

The result is as follows:

```
Ran 54 of 78 Specs in 85.288 seconds
SUCCESS! -- 54 Passed | 0 Failed | 1 Pending | 23 Skipped
```

## CurveFS

