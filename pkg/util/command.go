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
	"os/exec"
)

func ExecCommandHost(command string, args []string) ([]byte, error) {
	// nsenter -t 1 -m -p -n -i -u command args...
	nsenterArgs := []string{"-t", "1", "-m", "-p", "-n", "-i", "-u", command}
	return ExecCommand("nsenter", append(nsenterArgs, args...))
}

func ExecCommand(command string, args []string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	return cmd.CombinedOutput()
}
