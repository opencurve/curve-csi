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

package curve

import (
	"net"
	"net/http"
	"strconv"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/klog/v2"
	"k8s.io/utils/mount"

	"github.com/opencurve/curve-csi/cmd/options"
	csicommon "github.com/opencurve/curve-csi/pkg/csi-common"
	"github.com/opencurve/curve-csi/pkg/curveservice"
	"github.com/opencurve/curve-csi/pkg/logs"
	"github.com/opencurve/curve-csi/pkg/util"
)

type curveDriver struct {
	driver *csicommon.CSIDriver

	ids *identityServer
	cs  *controllerServer
	ns  *nodeServer
}

func NewCurveDriver() *curveDriver {
	return &curveDriver{}
}

func NewIdentityServer(d *csicommon.CSIDriver) *identityServer {
	return &identityServer{
		DefaultIdentityServer: csicommon.NewDefaultIdentityServer(d),
	}
}

func NewControllerServer(d *csicommon.CSIDriver, curveVolumePrefix string) *controllerServer {
	return &controllerServer{
		DefaultControllerServer: csicommon.NewDefaultControllerServer(d),
		volumeLocks:             util.NewVolumeLocks(),
		curveVolumePrefix:       curveVolumePrefix,
	}
}

func NewNodeServer(d *csicommon.CSIDriver, curveVolumePrefix string) *nodeServer {
	curveservice.InitCurveNbd()
	mounter := mount.New("")
	return &nodeServer{
		DefaultNodeServer: csicommon.NewDefaultNodeServer(d),
		mounter:           mounter,
		volumeLocks:       util.NewVolumeLocks(),
		curveVolumePrefix: curveVolumePrefix,
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

func (c *curveDriver) Run(curveConf options.CurveConf) {
	// Initialize default library driver
	c.driver = csicommon.NewCSIDriver(curveConf.DriverName, util.Version, curveConf.NodeID)
	if c.driver == nil {
		klog.Fatalln("Failed to initialize CSI Driver")
	}
	if curveConf.IsControllerServer || !curveConf.IsNodeServer {
		c.driver.AddControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{
			csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
			csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT,
			csi.ControllerServiceCapability_RPC_CLONE_VOLUME,
			csi.ControllerServiceCapability_RPC_EXPAND_VOLUME,
		})
		c.driver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{
			csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
			csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
		})
	}

	c.ids = NewIdentityServer(c.driver)
	if curveConf.IsControllerServer {
		c.cs = NewControllerServer(c.driver, curveConf.CurveVolumePrefix)
	}
	if curveConf.IsNodeServer {
		c.ns = NewNodeServer(c.driver, curveConf.CurveVolumePrefix)
	}

	if !curveConf.IsControllerServer && !curveConf.IsNodeServer {
		c.cs = NewControllerServer(c.driver, curveConf.CurveVolumePrefix)
		c.ns = NewNodeServer(c.driver, curveConf.CurveVolumePrefix)
	}

	s := csicommon.NewNonBlockingGRPCServer()
	s.Start(curveConf.Endpoint, c.ids, c.cs, c.ns)

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
