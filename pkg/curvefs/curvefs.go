/*
Copyright 2020 The Netease Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package curvefs

import (
	"net"
	"net/http"
	"strconv"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/klog/v2"
	"k8s.io/utils/mount"

	"github.com/opencurve/curve-csi/cmd/options"
	csicommon "github.com/opencurve/curve-csi/pkg/csi-common"
	"github.com/opencurve/curve-csi/pkg/logs"
	"github.com/opencurve/curve-csi/pkg/util"
)

type fsDriver struct {
	driver *csicommon.CSIDriver

	ids *identityServer
	cs  *controllerServer
	ns  *nodeServer
}

func NewCurveFSDriver() *fsDriver {
	return &fsDriver{}
}

func NewIdentityServer(d *csicommon.CSIDriver) *identityServer {
	return &identityServer{
		DefaultIdentityServer: csicommon.NewDefaultIdentityServer(d),
	}
}

func NewControllerServer(d *csicommon.CSIDriver) *controllerServer {
	return &controllerServer{
		DefaultControllerServer: csicommon.NewDefaultControllerServer(d),
		volumeLocks:             util.NewVolumeLocks(),
		snapshotLocks:           util.NewVolumeLocks(),
	}
}

func NewNodeServer(d *csicommon.CSIDriver) *nodeServer {
	mounter := mount.New("")
	return &nodeServer{
		DefaultNodeServer: csicommon.NewDefaultNodeServer(d),
		mounter:           mounter,
		volumeLocks:       util.NewVolumeLocks(),
	}
}

func listenAndServeDebugger(port int) {
	address := "127.0.0.1"
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/flags/v", util.StringFlagPutHandler(logs.GlogSetter))

	klog.Infof("starting debug http server to listen on %s:%d", address, port)
	err := http.ListenAndServe(net.JoinHostPort(address, strconv.Itoa(port)), mux)
	if err != nil {
		klog.Errorf("can not start debug http server, err: %v", err)
	}
}

func (f *fsDriver) Run(curveConf options.CurveConf) {
	// Initialize default library driver
	f.driver = csicommon.NewCSIDriver(curveConf.DriverName, util.Version, curveConf.NodeID)
	if f.driver == nil {
		klog.Fatalln("Failed to initialize CSI Driver")
	}
	if curveConf.IsControllerServer || !curveConf.IsNodeServer {
		f.driver.AddControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{
			csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
			csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT,
			csi.ControllerServiceCapability_RPC_CLONE_VOLUME,
			csi.ControllerServiceCapability_RPC_EXPAND_VOLUME,
		})
		f.driver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{
			csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
			csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
		})
	}

	f.ids = NewIdentityServer(f.driver)
	if curveConf.IsControllerServer {
		f.cs = NewControllerServer(f.driver)
	}
	if curveConf.IsNodeServer {
		f.ns = NewNodeServer(f.driver)
	}

	if !curveConf.IsControllerServer && !curveConf.IsNodeServer {
		f.cs = NewControllerServer(f.driver)
		f.ns = NewNodeServer(f.driver)
	}

	s := csicommon.NewNonBlockingGRPCServer()
	s.Start(curveConf.Endpoint, f.ids, f.cs, f.ns)

	// start debug server
	if curveConf.DebugPort > 0 {
		go listenAndServeDebugger(curveConf.DebugPort)
	}
	if curveConf.EnableProfiling {
		klog.Infof("Registering profiling handler")
		go util.EnableProfiling()
	}

	s.Wait()
}
