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

package main

import (
	"flag"
	"fmt"
	"os"

	"k8s.io/klog/v2"

	"github.com/opencurve/curve-csi/cmd/options"
	"github.com/opencurve/curve-csi/pkg/curve"
	"github.com/opencurve/curve-csi/pkg/logs"
	"github.com/opencurve/curve-csi/pkg/util"
)

const (
	driverDefaultName = "curve.csi.netease.com"
)

var (
	curveConf   options.CurveConf
	showVersion = flag.Bool("version", false, "Print version")
)

func init() {
	// common flags
	flag.StringVar(&curveConf.Endpoint, "endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	flag.StringVar(&curveConf.DriverName, "drivername", "", "name of the driver")
	flag.StringVar(&curveConf.NodeID, "nodeid", "", "node id")

	// CSI spec flags
	flag.BoolVar(&curveConf.IsNodeServer, "node-server", false, "start curve-csi node server")
	flag.BoolVar(&curveConf.IsControllerServer, "controller-server", false, "start curve-csi controller server")

	// curve volume name prefix
	flag.StringVar(&curveConf.CurveVolumePrefix, "curve-volume-prefix", "csi-vol-", "curve volume name prefix")

	// debug
	flag.IntVar(&curveConf.DebugPort, "debug-port", 0, "debug port, set 0 to disable")
	flag.BoolVar(&curveConf.EnableProfiling, "enableprofiling", false, "enable go profiling")
}

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Println(util.GetVersion())
		os.Exit(0)
	}
	logs.InitLogs()
	defer logs.FlushLogs()

	if curveConf.DriverName == "" {
		curveConf.DriverName = driverDefaultName
	}
	if err := util.ValidateDriverName(curveConf.DriverName); err != nil {
		klog.Fatalln(err)
	}

	klog.Infof("Starting the driver %s with version: %v", curveConf.DriverName, util.GetVersion())
	curveDriver := curve.NewCurveDriver()
	curveDriver.Run(curveConf)
	os.Exit(0)
}
