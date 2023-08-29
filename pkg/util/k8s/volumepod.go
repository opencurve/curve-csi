package k8s

import (
	"bytes"
	"context"
	"fmt"
	"os"

	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func MapVolumeInPod(ctx context.Context, podName string, command string) (output string, err error) {
	namespace := os.Getenv("POD_NAMESPACE")
	imageName := os.Getenv("CLIENT_IMAGE")
	// check volume pod exists
	pod, err := getPod(ctx, podName, namespace)
	if err != nil {
		return "", err
	}
	// pod not exist
	if pod == nil {
		// create volume pod
		if err := createVolumePod(ctx, podName, namespace, imageName); err != nil {
			return "", err
		}
		return "", fmt.Errorf("waiting for the pod running, pod:%s", podName)
	}
	if pod.Status.Phase != v1.PodRunning {
		return "", fmt.Errorf("waiting for the pod running, pod:%s, phase:%v ", podName, pod.Status.Phase)
	}
	// map volume
	return execCmd(ctx, podName, namespace, command)
}

func UnMapVolumeInPod(ctx context.Context, podName string, command string) (err error) {
	namespace := os.Getenv("POD_NAMESPACE")
	// check volume pod exists
	pod, err := getPod(ctx, podName, namespace)
	if err != nil {
		return err
	}
	// pod not exist, return unmap succeeded
	if pod == nil {
		return nil
	}
	// unmap volume
	_, err = execCmd(ctx, podName, namespace, command)
	if err != nil {
		return err
	}
	// delete volume pod
	if err := deletePod(ctx, podName, namespace); err != nil {
		return err
	}
	return
}

// get pod info
func getPod(ctx context.Context, podName, namespace string) (pod *v1.Pod, err error) {
	client, err := NewK8sClient()
	if err != nil {
		return nil, fmt.Errorf("can not get pod %s information, failed to connect to Kubernetes: %w", podName, err)
	}

	pod, err = client.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get pod %s information: %w", podName, err)
	}
	return pod, nil
}

// create volume pod
func createVolumePod(ctx context.Context, podName, namespace, imageName string) (err error) {
	client, err := NewK8sClient()
	if err != nil {
		return fmt.Errorf("failed to connect to Kubernetes: %w", err)
	}
	pod, err := genVolumePod(podName, namespace, imageName)
	if err != nil {
		return fmt.Errorf("failed to create volume pod: %w", err)
	}
	_, err = client.CoreV1().Pods(namespace).Create(ctx, pod, metav1.CreateOptions{FieldManager: os.Getenv("HOSTNAME")})
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			return nil
		}
		return err
	}
	return
}

func deletePod(ctx context.Context, podName, namespace string) (err error) {
	client, err := NewK8sClient()
	if err != nil {
		return fmt.Errorf("failed to connect to Kubernetes: %w", err)
	}
	err = client.CoreV1().Pods(namespace).Delete(ctx, podName, metav1.DeleteOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}
	return
}

// exec command on specific pod and wait the command's output.
func execCmd(ctx context.Context, podName, namespace string, command string) (out string, err error) {
	client, err := NewK8sClient()
	if err != nil {
		return "", fmt.Errorf("failed to connect to Kubernetes: %w", err)
	}
	cmd := []string{
		"sh",
		"-c",
		command,
	}
	req := client.CoreV1().RESTClient().Post().Resource("pods").Name(podName).
		Namespace(namespace).SubResource("exec")
	option := &v1.PodExecOptions{
		Command: cmd,
		Stdin:   false,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}
	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)
	var stdout, stderr bytes.Buffer
	cfg, _ := rest.InClusterConfig()
	exec, err := remotecommand.NewSPDYExecutor(cfg, "POST", req.URL())
	if err != nil {
		return "", err
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stderr: &stderr,
		Stdout: &stdout,
	})
	if err != nil {
		return "", err
	}
	if stderr.String() != "" {
		return stdout.String(), fmt.Errorf("exec [%s] failed, stderr:%s", command, stderr.String())
	}
	return stdout.String(), nil
}
