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

type PodLogger struct {
	client corev1client.CoreV1Interface
	config *rest.Config
}

func NewPodLogger() (*PodLogger, error) {
	client, config, err := ClientGenerator()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &PodLogger{
		client: client,
		config: config,
	}, nil
}

func (p *PodLogger) GetPodLogs(podName, namespace string) ([]byte, []byte, error) {
	request := p.client.RESTClient().
		Get().
		Resource("pods").
		Namespace(namespace).
		Name(podName).
		SubResource("log").
		VersionedParams(&corev1.PodLogOptions{
			InsecureSkipTLSVerifyBackend: true,
			Follow:                       true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(p.config, "GET", request.URL())
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
