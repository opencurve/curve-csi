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
	"runtime"
)

var (
	BuildTime string
	GitCommit string
	Version   string
)

func GetVersion() string {
	format := `Version:    %s
Go version: %s
Platform: %s/%s
Git commit: %s
Built:      %s
`
	return fmt.Sprintf(format, Version, runtime.Version(), runtime.GOOS, runtime.GOARCH, GitCommit, BuildTime)
}
