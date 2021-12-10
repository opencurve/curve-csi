FROM golang:1.16.11

WORKDIR /go/src/github.com/opencurve/curve-csi
COPY . .
ARG VERSION
RUN VERSION=${VERSION} make build

FROM us.gcr.io/k8s-artifacts-prod/build-image/debian-base-amd64:v2.1.0
LABEL name="curve-csi" \
      license="Apache Version 2.0" \
      maintainers="The Netease Kubernetes Authors" \
      description="Curve CSI Plugin" \
      architecture="amd64" \
      distribution-scope="public" \
      summary="Curve CSI Plugin" \
      help="none"
RUN clean-install ca-certificates e2fsprogs mount xfsprogs udev
COPY --from=0 /go/src/github.com/opencurve/curve-csi/_output/curve-csi /bin/curve-csi
RUN chmod +x /bin/curve-csi
    
ENTRYPOINT ["/bin/curve-csi"]
