#!/bin/bash
# modify MDS addr
sed -i "s/mds\.listen\.addr=.*/mds.listen.addr=${MDSADDR}/" /etc/curve/client.conf

# start nebd
nebd-daemon start

cd /bin && ./curve-csi "$@"