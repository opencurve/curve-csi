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
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation"

	"github.com/opencurve/curve-csi/pkg/util/ctxlog"
)

// ValidateDriverName validates the driver name
func ValidateDriverName(driverName string) error {
	if driverName == "" {
		return errors.New("driver name is empty")
	}

	if len(driverName) > 63 {
		return errors.New("driver name length should be less than 63 chars")
	}
	var err error
	for _, msg := range validation.IsDNS1123Subdomain(strings.ToLower(driverName)) {
		if err == nil {
			err = errors.New(msg)
			continue
		}
		err = errors.Wrap(err, msg)
	}
	return err
}

func SystemMapOnHost(ctx context.Context, serviceName string, mapCommands []string) (err error) {
	ctxlog.Infof(ctx, "starting to run %s.service", serviceName)
	systemMapArgs := []string{"--description=k8scsi", "--unit", serviceName, "-r", "--"}
	systemMapArgs = append(systemMapArgs, mapCommands...)

	var output []byte
	defer func() {
		// tear down
		if err != nil {
			output, _ = ExecCommandHost("systemctl", []string{"status", serviceName})
			ctxlog.Warningf(ctx, "systemctl status %s, output: %s", serviceName, string(output))
			_, _ = ExecCommandHost("systemctl", []string{"stop", serviceName})
			_, _ = ExecCommandHost("systemctl", []string{"reset-failed", serviceName})
		}
	}()

	output, err = ExecCommandHost("systemd-run", systemMapArgs)
	if err != nil {
		// service already exists, reset it and try again
		if !strings.Contains(string(output), "already exists") {
			return fmt.Errorf("failed to map, output: %s", string(output))
		}
		ctxlog.Warningf(ctx, "systemctl reset-failed %s.service and try mapping again", serviceName)
		_, _ = ExecCommandHost("systemctl", []string{"reset-failed", serviceName})
		_, _ = ExecCommandHost("systemd-run", systemMapArgs)
	}
	// check service status
	output, err = ExecCommandHost("systemctl", []string{"show", serviceName, "-p", "ExecMainStatus"})
	if err != nil {
		return fmt.Errorf("systemctl show %s.service failed, output: %s", serviceName, string(output))
	}
	if !strings.Contains(string(output), "ExecMainStatus=0") {
		return fmt.Errorf("%s.service started successfully, but map failed, %s", serviceName, string(output))
	}
	ctxlog.Infof(ctx, "map successfully, running as %s.service", serviceName)
	return nil
}

type NotFoundErr struct {
	Id string
}

func NewNotFoundErr(id ...string) *NotFoundErr {
	if len(id) != 0 {
		return &NotFoundErr{Id: strings.Join(id, ",")}
	}
	return &NotFoundErr{}
}

func (e *NotFoundErr) Error() string {
	if e.Id == "" {
		return "Not found"
	}
	return fmt.Sprintf("Object with Id(%v) not found", e.Id)
}

func IsNotFoundErr(err error, id ...string) bool {
	notFoundErr := NewNotFoundErr(id...)
	return err.Error() == notFoundErr.Error()
}

func GenVolumePodName(volumeName string) string {
	return fmt.Sprintf("%s-%s", volumeName, os.Getenv("NODE_ID"))
}
