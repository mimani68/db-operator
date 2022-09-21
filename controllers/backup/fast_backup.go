package backup

import (
	"github.com/go-logr/logr"
	"github.com/mimani68/db-operator/pkg/exec"
)

func ImidiateBackup(dbUrlConnection string, log logr.Logger) {
	// status := corev1.ConditionUnknown
	executor, err := exec.NewPodExecutor()
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

}
