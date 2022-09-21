/*
Copyright 2022.

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

package controllers

import (
	"context"
	"encoding/json"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
	opsv1alpha1 "github.com/mimani68/db-operator/api/v1alpha1"
)

// BackupReconciler reconciles a Backup object
type BackupReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ops.db.io,resources=backups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ops.db.io,resources=backups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ops.db.io,resources=backups/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Backup object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *BackupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	dbBackup := &opsv1alpha1.Backup{}
	err := r.Get(ctx, req.NamespacedName, dbBackup)
	if err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("dbBackup resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		r.Log.Error(err, "Failed to get dbBackup")
		return ctrl.Result{}, err
	} else {
		r.Log.Info("dbBackup resource found.")
		value, err := json.Marshal(dbBackup)
		if err == nil {
			r.Log.Info(string(value))
		}
	}

	// found := &appsv1.Deployment{}
	// svc := &corev1.Service{}

	// //
	// // Exec sample
	// //
	// command := []string{"uname", "-a"}
	// // status := corev1.ConditionUnknown
	// executor, err := exec.NewPodExecutor()
	// if err != nil {
	// 	r.Log.Error(err, err.Error())
	// }

	// stdout, stderr, err := executor.Exec("echo", "echo-57cd9f4f98-7c2v7", "", command...)
	// if err != nil {
	// 	r.Log.Error(err, err.Error())
	// }
	// if stderr != nil {
	// 	r.Log.Error(nil, string(stderr))
	// }
	// r.Log.Info(string(stdout))

	//
	//
	// === Get pods ===
	//
	//
	// pod := &corev1.Pod{}
	// err = r.Get(ctx, types.NamespacedName{
	// 	Name:      "echo-57cd9f4f98-7c2v7",
	// 	Namespace: "echo",
	// }, pod)

	// r.Log.Info("Get pods of \"echo\" namespace")
	// if err != nil && errors.IsNotFound(err) {
	// 	r.Log.Error(err, err.Error())
	// } else {
	// 	value, _ := json.Marshal(pod)
	// 	r.Log.Info(string(value))
	// }

	//
	// === X ===
	//
	// // Check if the deployment already exists, if not create a new one
	// found := &appsv1.Deployment{}
	// err = r.Get(ctx, req.NamespacedName, found)
	// if err != nil && errors.IsNotFound(err) {
	// 	// Define a new deployment
	// 	dep := r.deploymentFordbBackup(dbBackup)
	// 	log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
	// 	err = r.Create(ctx, dep)
	// 	if err != nil {
	// 		log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
	// 		return ctrl.Result{}, err
	// 	}
	// 	// Deployment created successfully - return and requeue
	// 	return ctrl.Result{Requeue: true}, nil
	// } else if err != nil {
	// 	log.Error(err, "Failed to get Deployment")
	// 	return ctrl.Result{}, err
	// }

	// // Ensure the deployment size is the same as the spec
	// size := dbBackup.Spec.DbUrl
	// if *found.Spec.Replicas != size {
	// 	found.Spec.Replicas = &size
	// 	err = r.Update(ctx, found)
	// 	if err != nil {
	// 		log.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
	// 		return ctrl.Result{}, err
	// 	}
	// 	// Spec updated - return and requeue
	// 	return ctrl.Result{Requeue: true}, nil
	// }

	// // Update the dbBackup status with the pod names
	// // List the pods for this dbBackup's deployment
	// podList := &corev1.PodList{}
	// listOpts := []client.ListOption{
	// 	client.InNamespace(dbBackup.Namespace),
	// 	client.MatchingLabels(labelsFordbBackup(dbBackup.Name)),
	// }
	// if err = r.List(ctx, podList, listOpts...); err != nil {
	// 	log.Error(err, "Failed to list pods", "dbBackup.Namespace", dbBackup.Namespace, "dbBackup.Name", dbBackup.Name)
	// 	return ctrl.Result{}, err
	// }
	// podNames := getPodNames(podList.Items)

	// // Update status.Nodes if needed
	// if !reflect.DeepEqual(podNames, dbBackup.Status.Nodes) {
	// 	dbBackup.Status.Nodes = podNames
	// 	err := r.Status().Update(ctx, dbBackup)
	// 	if err != nil {
	// 		log.Error(err, "Failed to update dbBackup status")
	// 		return ctrl.Result{}, err
	// 	}
	// }

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BackupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&opsv1alpha1.Backup{}).
		Complete(r)
}
