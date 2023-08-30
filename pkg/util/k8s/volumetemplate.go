package k8s

import (
	"fmt"
	"os"

	v1 "k8s.io/api/core/v1"
)

func genVolumePod(name, namespace, imageName string) (*v1.Pod, error) {
	mds := os.Getenv("MDSADDR")
	nodeName := os.Getenv("NODE_ID")
	if name == "" {
		return nil, fmt.Errorf("name missing")
	}
	if namespace == "" {
		return nil, fmt.Errorf("namespace missing")
	}
	if mds == "" {
		return nil, fmt.Errorf("mds missing")
	}
	if imageName == "" {
		return nil, fmt.Errorf("imageName missing")
	}
	if nodeName == "" {
		return nil, fmt.Errorf("nodeName missing")
	}

	pod := &v1.Pod{}
	privileged := true
	allowPrivilegeEscalation := true
	pod.Name = name
	pod.Namespace = namespace
	pod.Spec = v1.PodSpec{
		Containers: []v1.Container{
			{
				Name:  "csi-volume",
				Image: imageName,
				Args: []string{
					"--endpoint=tcp://127.0.0.1:8000",
					"--drivername=curve.csi.netease.com",
					"--nodeid=$(NODE_ID)",
					"--client=true",
					"--logtostderr=true",
					"-v=1",
				},
				SecurityContext: &v1.SecurityContext{
					Capabilities: &v1.Capabilities{
						Add: []v1.Capability{
							"SYS_ADMIN",
						},
					},
					Privileged:               &privileged,
					AllowPrivilegeEscalation: &allowPrivilegeEscalation,
				},
				Env: []v1.EnvVar{
					{
						Name: "POD_IP",
						ValueFrom: &v1.EnvVarSource{
							FieldRef: &v1.ObjectFieldSelector{
								APIVersion: "v1",
								FieldPath:  "status.podIP",
							},
						},
					},
					{
						Name: "NODE_ID",
						ValueFrom: &v1.EnvVarSource{
							FieldRef: &v1.ObjectFieldSelector{
								APIVersion: "v1",
								FieldPath:  "spec.nodeName",
							},
						},
					},
					{
						Name: "POD_NAMESPACE",
						ValueFrom: &v1.EnvVarSource{
							FieldRef: &v1.ObjectFieldSelector{
								APIVersion: "v1",
								FieldPath:  "metadata.namespace",
							},
						},
					},
					{
						Name:  "MDSADDR",
						Value: mds,
					},
				},
				VolumeMounts: []v1.VolumeMount{
					{
						Name:      "host-dev",
						MountPath: "/dev",
					},
					{
						Name:      "host-sys",
						MountPath: "/sys",
					},
					{
						Name:      "lib-modules",
						MountPath: "/lib/modules",
						ReadOnly:  true,
					},
					{
						Name:      "localtime",
						MountPath: " /etc/localtime",
					},
				},
			},
		},
		Volumes: []v1.Volume{
			{
				Name: "host-dev",
				VolumeSource: v1.VolumeSource{
					HostPath: &v1.HostPathVolumeSource{
						Path: "/dev",
					},
				},
			},
			{
				Name: "host-sys",
				VolumeSource: v1.VolumeSource{
					HostPath: &v1.HostPathVolumeSource{
						Path: "/sys",
					},
				},
			},
			{
				Name: "lib-modules",
				VolumeSource: v1.VolumeSource{
					HostPath: &v1.HostPathVolumeSource{
						Path: "/lib/modules",
					},
				},
			},
			{
				Name: "localtime",
				VolumeSource: v1.VolumeSource{
					HostPath: &v1.HostPathVolumeSource{
						Path: "/etc/localtime",
					},
				},
			},
		},
		NodeName:    nodeName,
		HostNetwork: true,
		HostPID:     true,
	}
	return pod, nil
}
