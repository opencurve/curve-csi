ARG DEBIAN_ARCH=amd64
FROM us.gcr.io/k8s-artifacts-prod/build-image/debian-base-${DEBIAN_ARCH}:v2.1.0
ARG ARCH=amd64
LABEL name="curve-csi" \
      license="Apache Version 2.0" \
      maintainers="The Netease Kubernetes Authors" \
      description="Curve CSI Plugin" \
      architecture=$ARCH \
      distribution-scope="public" \
      summary="Curve CSI Plugin" \
      help="none"

ADD curve-csi /bin/curve-csi
RUN chmod +x /bin/curve-csi && \
    clean-install ca-certificates e2fsprogs mount xfsprogs udev

CMD ["/bin/curve-csi"]
