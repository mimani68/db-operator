package backup

import (
	"github.com/go-logr/logr"
	"github.com/mimani68/db-operator/pkg/k8s"
)

func SchedualBackup(dbUrlConnection, typeOfDb string, log logr.Logger) (bool, error) {
	// status := corev1.ConditionUnknown

	logger, err := k8s.NewPodLogger()
	if err != nil {
		log.Error(err, err.Error())
	}
	stdOut, stdError, err := logger.GetPodLogs("echo-57cd9f4f98-pgdv4", "echo")
	if err != nil {
		log.Error(err, err.Error())
	}
	if stdError != nil {
		log.Info(string(stdError))
	}
	if err != nil {
		log.Info(string(stdOut))
	}

	// executor, err := k8s.NewPodExecutor()
	// if err != nil {
	// 	log.Error(err, err.Error())
	// }

	// command := []string{"uname", "-a"}
	// stdout, stderr, err := executor.Exec("echo", "echo-57cd9f4f98-7c2v7", "", command...)
	// if err != nil {
	// 	log.Error(err, err.Error())
	// }
	// if stderr != nil {
	// 	log.Error(nil, string(stderr))
	// }
	// log.Info(string(stdout))

	return true, nil
}

func ImidiateBackup(dbUrlConnection, typeOfDb string, log logr.Logger) (bool, error) {
	// status := corev1.ConditionUnknown
	executor, err := k8s.NewPodExecutor()
	if err != nil {
		log.Error(err, err.Error())
	}

	command := []string{"uname", "-a"}
	stdout, stderr, err := executor.Exec("echo", "echo-57cd9f4f98-7c2v7", "", command...)
	if err != nil {
		log.Error(err, err.Error())
	}
	if stderr != nil {
		log.Error(nil, string(stderr))
	}
	log.Info(string(stdout))
	return true, nil
}
