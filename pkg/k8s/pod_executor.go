package k8s

import (
	"bufio"
	"bytes"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type PodExecutor struct {
	client corev1client.CoreV1Interface
	config *rest.Config
}

func NewPodExecutor() (*PodExecutor, error) {
	client, config, err := ClientGenerator()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &PodExecutor{
		client: client,
		config: config,
	}, nil
}

func (p *PodExecutor) Exec(namespace, podName, containerName string, command ...string) ([]byte, []byte, error) {
	request := p.client.RESTClient().
		Post().
		Resource("pods").
		Namespace(namespace).
		Name(podName).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   command,
			Stdout:    true,
			Stderr:    true,
			Stdin:     false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(p.config, "POST", request.URL())
	if err != nil {
		return nil, nil, err
	}

	stdOut := bytes.Buffer{}
	stdErr := bytes.Buffer{}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: bufio.NewWriter(&stdOut),
		Stderr: bufio.NewWriter(&stdErr),
		Stdin:  nil,
		Tty:    false,
	})

	return stdOut.Bytes(), stdErr.Bytes(), err
}

func (p *PodExecutor) SetGlobalSysVar(namespace, podName string, query string) error {
	cmd := []string{"xenoncli", "mysql", "sysvar", query}
	_, stderr, err := p.Exec(namespace, podName, "xenon", cmd...)
	if err != nil {
		return err
	}
	if len(stderr) != 0 {
		return fmt.Errorf("run command %s in xenon failed: %s", cmd, stderr)
	}
	return nil
}
