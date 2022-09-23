package k8s

import (
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func ClientGenerator() (*corev1client.CoreV1Client, *rest.Config, error) {
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	config, err := kubeconfig.ClientConfig()
	if err != nil {
		return nil, nil, err
	}

	client, err := corev1client.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	return client, config, nil
}
