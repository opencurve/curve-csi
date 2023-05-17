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
	"strconv"
	"strings"

	"github.com/opencurve/curve-csi/pkg/util"
	"github.com/opencurve/curve-csi/pkg/util/ctxlog"
)

const (
	// curve cli return code
	retExist      = 1
	retFailed     = 2
	retAuthFailed = 4
	retNotExist   = 6

	retFailFormat = "fail, ret = -%d"

	// curve file status
	CurveVolumeStatusNotExist      CurveVolumeStatus = "kFileNotExists"
	CurveVolumeStatusExist         CurveVolumeStatus = "kFileExists"
	CurveVolumeStatusCreated       CurveVolumeStatus = "Created"
	CurveVolumeStatusOwnerAuthFail CurveVolumeStatus = "kOwnerAuthFail"
	CurveVolumeStatusClonedLazy    CurveVolumeStatus = "CloneMetaInstalled"
	CurveVolumeStatusBeingCloned   CurveVolumeStatus = "BeingCloned"
	CurveVolumeStatusUnknown       CurveVolumeStatus = "Unknown"
)

type CurveVolume struct {
	FileName string `json:"filename"`
	FilePath string `json:"filepath"`
	DirPath  string `json:"dirpath"`
	User     string `json:"user"`
	SizeGiB  int    `json:"size"`
}

func NewCurveVolume(user, volName string, sizeGiB int) *CurveVolume {
	return &CurveVolume{
		FileName: volName,
		FilePath: "/" + user + "/" + volName,
		DirPath:  "/" + user,
		User:     user,
		SizeGiB:  sizeGiB,
	}
}

// curve stat [-h] --user USER --filename FILENAME
func (cv *CurveVolume) Stat(ctx context.Context) (*CurveVolumeDetail, error) {
	args := []string{"stat", "--user", cv.User, "--filename", cv.FilePath}
	ctxlog.V(4).Infof(ctx, "starting exec: curve %v", args)
	output, err := util.ExecCommand("curve", args)
	outputStr := string(output)
	if err == nil {
		ctxlog.V(5).Infof(ctx, "[curve] successfully stat the volume, output: %v", outputStr)
		return simpleParseVolumeDetail(output)
	}

	ctxlog.Warningf(ctx, "[curve] failed to stat the file %s, err: %v, output: %v", cv.FilePath, err, outputStr)
	if strings.Contains(outputStr, fmt.Sprintf(retFailFormat, retNotExist)) {
		return nil, util.NewNotFoundErr()
	}

	return nil, fmt.Errorf("can not run curve %v, err: %v, output: %v", args, err, outputStr)
}

// curve create file, mkdir the dir if not exists
func (cv *CurveVolume) Create(ctx context.Context) error {
	output, err := cv.create(ctx)
	if err == nil {
		ctxlog.V(4).Infof(ctx, "[curve] successfully create %v", cv.FilePath)
		return nil
	}

	if strings.Contains(string(output), fmt.Sprintf(retFailFormat, retNotExist)) {
		ctxlog.V(4).Infof(ctx, "[curve] try to mkdir %s before creating volume %s", cv.DirPath, cv.FilePath)
		if err := cv.mkdir(ctx); err != nil {
			return fmt.Errorf("failed to mkdir %v, err: %v", cv.DirPath, err)
		}
		// recreate
		output, err = cv.create(ctx)
		if err == nil {
			ctxlog.V(4).Infof(ctx, "[curve] successfully create %v", cv.FilePath)
			return nil
		}
	}

	return fmt.Errorf("failed to create %s, err: %v, output: %v", cv.FilePath, err, string(output))
}

// curve delete [-h] --user USER --filename FILENAME
func (cv *CurveVolume) Delete(ctx context.Context) error {
	args := []string{"delete", "--user", cv.User, "--filename", cv.FilePath}
	ctxlog.V(4).Infof(ctx, "starting exec: curve %v", args)
	output, err := util.ExecCommand("curve", args)
	if err != nil {
		if strings.Contains(string(output), fmt.Sprintf(retFailFormat, retNotExist)) {
			ctxlog.Warningf(ctx, "[curve] the file %s already deleted, ignore deleting it", cv.FilePath)
			return nil
		}
		return fmt.Errorf("failed to delete %s, err: %v, output: %v", cv.FilePath, err, string(output))
	}

	ctxlog.V(4).Infof(ctx, "[curve] successfully delete %v", cv.FilePath)
	return nil
}

// curve extend [-h] --user USER --filename FILENAME --length LENGTH
func (cv *CurveVolume) Extend(ctx context.Context, newSizeGiB int) error {
	volLength := strconv.Itoa(newSizeGiB)
	args := []string{"extend", "--user", cv.User, "--filename", cv.FilePath, "--length", volLength}
	ctxlog.V(4).Infof(ctx, "starting exec: curve %v", args)
	output, err := util.ExecCommand("curve", args)
	if err != nil {
		return fmt.Errorf("failed to extend %s, err: %v, output: %v", cv.FilePath, err, string(output))
	}
	ctxlog.V(4).Infof(ctx, "[curve] successfully extend %v", cv.FilePath)
	return nil
}

// curve mkdir [-h] --user USER --dirname DIRNAME
func (cv *CurveVolume) mkdir(ctx context.Context) error {
	args := []string{"mkdir", "--user", cv.User, "--dirname", cv.DirPath}
	ctxlog.V(4).Infof(ctx, "starting exec: curve %v", args)
	output, err := util.ExecCommand("curve", args)
	if err != nil {
		if strings.Contains(string(output), fmt.Sprintf(retFailFormat, retExist)) {
			ctxlog.V(4).Infof(ctx, "[curve] the dir %s of user %s already exists, ignore to mkdir", cv.DirPath, cv.User)
			return nil
		}
		return fmt.Errorf("failed to run curve %v, err: %v, output: %v", args, err, string(output))
	}
	ctxlog.V(4).Infof(ctx, "[curve] successfully mkdir %s of user %s", cv.DirPath, cv.User)
	return nil
}

// curve create [-h] --filename FILENAME --length LENGTH --user USER
func (cv *CurveVolume) create(ctx context.Context) (output []byte, err error) {
	volLength := strconv.Itoa(cv.SizeGiB)
	args := []string{"create", "--filename", cv.FilePath, "--length", volLength, "--user", cv.User}
	ctxlog.V(4).Infof(ctx, "starting exec: curve %v", args)
	output, err = util.ExecCommand("curve", args)
	if err != nil {
		if strings.Contains(string(output), fmt.Sprintf(retFailFormat, retExist)) {
			ctxlog.Warningf(ctx, "[curve] the file %s already exists, ignore recreating it", cv.FilePath)
			return output, nil
		}
	}
	ctxlog.V(5).Infof(ctx, "[curve] create result: %v, err: %v", string(output), err)
	return output, err
}

// curve list [-h] --user USER --dirname DIRNAME
func (cv *CurveVolume) list(ctx context.Context) ([]string, error) {
	args := []string{"list", "--user", cv.User, "--dirname", cv.DirPath}
	ctxlog.V(4).Infof(ctx, "starting exec: curve %v", args)
	output, err := util.ExecCommand("curve", args)
	outputStr := string(output)
	if err != nil {
		if strings.Contains(outputStr, fmt.Sprintf(retFailFormat, retNotExist)) {
			ctxlog.Warningf(ctx, "[curve] the %s not exist, output: %v", cv.DirPath, outputStr)
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to run curve %v, err: %v, output: %v", args, err, outputStr)
	}

	ctxlog.V(4).Infof(ctx, "[curve] get volumes: %v in %v", outputStr, cv.DirPath)
	volumes := make([]string, 0)
	for _, line := range strings.Split(string(output), "\n") {
		pLine := strings.TrimSpace(line)
		if pLine != "" {
			volumes = append(volumes, pLine)
		}
	}

	return volumes, nil
}

// curve-nbd map cbd:<user>/<filename_full_path>_<user>_
func (cv *CurveVolume) Map(ctx context.Context, disableInUseChecks bool) (string, error) {
	devicePath, found := waitForMapped(ctx, cv.FilePath, cv.User, 1)
	if found {
		ctxlog.V(4).Infof(ctx, "[curve-nbd] the curve file %s already mapped at %v", cv.FilePath, devicePath)
		return devicePath, nil
	}

	ctxlog.Infof(ctx, "[curve-nbd] starting to attach curve file: %s", cv.FilePath)

	// wait for curve image status available and able to mapped
	if err := waitForCurveFileReady(ctx, cv.FileName, cv.User, disableInUseChecks); err != nil {
		return "", fmt.Errorf("curve file %s may not be ready, err: %v", cv.FilePath, err)
	}

	// map device
	cbdMapPath := fmt.Sprintf("cbd:%s/%s_%s_", cv.User, cv.FilePath, cv.User)
	args := []string{"map", cbdMapPath, "--timeout", "86400"}
	go util.ExecCommand(curveNbdCmd, args)

	devicePath, found = waitForMapped(ctx, cv.FilePath, cv.User, 10)
	if !found {
		return "", fmt.Errorf("can not find devicePath after mapping successfully")
	}

	return devicePath, nil
}

// curve-nbd unmap
func (cv *CurveVolume) UnMap(ctx context.Context) error {
	devicePath, err := getNbdDevFromFileName(ctx, cv.FilePath, cv.User)
	if err != nil {
		return err
	}

	// unmap
	output, err := util.ExecCommand(curveNbdCmd, []string{"unmap", devicePath})
	if err != nil {
		return fmt.Errorf("curve: unmap file %s failed, err: %v, output: %v", cv.FilePath, err, string(output))
	}

	return nil
}
