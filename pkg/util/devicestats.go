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

package util

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/sys/unix"
)

type DeviceStats struct {
	Block bool

	AvailableBytes  int64
	TotalBytes      int64
	UsedBytes       int64
	AvailableInodes int64
	TotalInodes     int64
	UsedInodes      int64
}

func GetDeviceStats(path string) (*DeviceStats, error) {
	isBlock, err := isBlockDevice(path)
	if isBlock {
		size, err := getBlockDeviceSize(path)
		if err != nil {
			return nil, err
		}

		return &DeviceStats{
			Block:      true,
			TotalBytes: size,
		}, nil
	}

	var statfs unix.Statfs_t
	err = unix.Statfs(path, &statfs)
	if err != nil {
		return nil, fmt.Errorf("failed to statfs() %q: %s", path, err)
	}

	deviceStats := &DeviceStats{
		Block: false,

		AvailableBytes: int64(statfs.Bavail) * int64(statfs.Bsize),
		TotalBytes:     int64(statfs.Blocks) * int64(statfs.Bsize),

		AvailableInodes: int64(statfs.Ffree),
		TotalInodes:     int64(statfs.Files),
	}
	deviceStats.UsedBytes = deviceStats.TotalBytes - deviceStats.AvailableBytes
	deviceStats.UsedInodes = deviceStats.TotalInodes - deviceStats.AvailableInodes
	return deviceStats, nil
}

func isBlockDevice(path string) (bool, error) {
	var stat unix.Stat_t
	err := unix.Stat(path, &stat)
	if err != nil {
		return false, fmt.Errorf("failed to stat() %q: %s", path, err)
	}

	return (stat.Mode & unix.S_IFMT) == unix.S_IFBLK, nil
}

func getBlockDeviceSize(path string) (int64, error) {
	fd, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("failed to open path %s", path)
	}
	defer fd.Close()
	pos, err := fd.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, fmt.Errorf("error seeking to end of %s: %s", path, err)
	}
	return pos, nil
}
