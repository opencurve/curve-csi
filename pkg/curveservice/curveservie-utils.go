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

package curveservice

import (
	"context"
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"

	"github.com/opencurve/curve-csi/pkg/util"
	"github.com/opencurve/curve-csi/pkg/util/ctxlog"
)

const (
	nbdsMax = 128
	// The following three values are used for 30 seconds timeout
	// while waiting for curve Watcher to expire.
	curveFileWatcherInitDelay = 1 * time.Second
	curveFileWatcherFactor    = 1.4
	curveFileWatcherSteps     = 10

	curveNbdCmd = "curve-nbd"
)

func InitCurveNbd() {
	output, err := util.ExecCommandHost("modprobe", []string{"nbd", fmt.Sprintf("nbds_max=%d", nbdsMax)})
	if err != nil {
		klog.Errorf("curve-nbd: nbd modprobe failed with error %v, output: %v", err, string(output))
	}

	running, err := checkNebdDaemonRunning()
	if err == nil && !running {
		klog.Errorf("nebd-daemon not started, please run: nebd-daemon start")
	}
}

func checkNebdDaemonRunning() (bool, error) {
	output, err := util.ExecCommand("nebd-daemon", []string{"status"})
	if err != nil {
		klog.Warningf("failed to run nebd-daemon status, output: %v, err: %v", string(output), err)
		return false, err
	}
	if strings.Contains(string(output), "is running") {
		return true, nil
	}
	return false, nil
}

// Stat a path, if it doesn't exist, retry maxRetries times.
func waitForMapped(ctx context.Context, filePath, user string, maxRetries int) (string, bool) {
	for i := 0; i < maxRetries; i++ {
		if i != 0 {
			time.Sleep(time.Second)
		}

		if devicePath, err := getNbdDevFromFileName(ctx, filePath, user); devicePath != "" {
			return devicePath, true
		} else if err != nil {
			klog.Warning(err)
		}
	}
	return "", false
}

// cmd "curve-nbd list-mapped" return nbd device mapped locally.
// id      image                                                                device
// 1509297 cbd:k8s//k8s/csi-vol-pvc-647525be-c0d6-464b-b548-1fa26f6d183c_k8s_ /dev/nbd1
func getNbdDevFromFileName(ctx context.Context, filePath, user string) (string, error) {
	output, err := util.ExecCommand("curve-nbd", []string{"list-mapped"})
	if err != nil {
		return "", fmt.Errorf("can not run curve-nbd list-mapped, err: %v, output: %s", err, string(output))
	}
	for _, l := range strings.Split(string(output), "\n") {
		// 1509297 cbd:k8s//k8s/csi-vol-pvc-647525be-c0d6-464b-b548-1fa26f6d183c_k8s_ /dev/nbd1
		tLine := strings.TrimSpace(l)
		if tLine == "" || strings.HasPrefix(tLine, "id") {
			continue
		}
		lineSlice := strings.Fields(tLine)
		if len(lineSlice) < 3 {
			continue
		}
		// cbdMapPathSuffix: /k8s/csi-vol-pvc-647525be-c0d6-464b-b548-1fa26f6d183c_k8s_
		cbdMapPathSuffix := fmt.Sprintf("%s_%s_", filePath, user)
		if strings.HasSuffix(lineSlice[1], cbdMapPathSuffix) {
			ctxlog.Infof(ctx, "get device path: %s of filePath: %s", lineSlice[2], filePath)
			return lineSlice[2], nil
		}
	}

	ctxlog.Warningf(ctx, "can't find devicePath of filePath: %s", filePath)
	return "", nil
}

// Wait for the curve file ready and not mapped at other nodes
func waitForCurveFileReady(ctx context.Context, fileName, user string, disableInUseChecks bool) error {
	backoff := wait.Backoff{
		Duration: curveFileWatcherInitDelay,
		Factor:   curveFileWatcherFactor,
		Steps:    curveFileWatcherSteps,
	}

	err := wait.ExponentialBackoff(backoff, func() (bool, error) {
		used, output, err := curveStatus(ctx, fileName, user)
		if err != nil {
			return false, fmt.Errorf("fail to check curve file %s status with: (%v), output: (%s)", fileName, err, output)
		}
		if disableInUseChecks && used {
			ctxlog.Infof(ctx, "valid multi-node attach requested, ignoring watcher in-use result")
			return used, nil
		}
		return !used, nil
	})
	// return error if curve image has not become available for the specified timeout
	if err == wait.ErrWaitTimeout {
		return fmt.Errorf("curve file %s is still being used", fileName)
	}
	// return error if any other errors were encountered during waiting for the image to become available
	return err
}

// curveStatus checks if there is watcher on the file.
// It returns true if there is a watcher on the file, otherwise returns false.
func curveStatus(ctx context.Context, fileName, user string) (bool, string, error) {
	// Not implemented yet
	return false, "", nil
}
