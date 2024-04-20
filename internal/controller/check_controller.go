/*
Copyright 2024.

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

package controller

import (
	"context"
	"fmt"
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CheckReconciler reconciles a Check object
type CheckReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const HostPathType = "DirectoryOrCreate"

//+kubebuilder:rbac:groups=ctf.check.com,resources=checks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ctf.check.com,resources=checks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ctf.check.com,resources=checks/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Check object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile

func (r *CheckReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	podNames := []string{"unit216", "x_ta", "cyber.spawn()", "to_netsaref", "hello world", "la7beb", "com'tech", "tsukuyomi", "bsissa"}

	for {
		select {
		case <-ctx.Done():
			// If the context is cancelled, stop the ticker and return
			return ctrl.Result{}, nil
		case <-ticker.C:
			// Every minute, check the status of the pods
			for _, podName := range podNames {
				pod := &corev1.Pod{}
				err := r.Get(ctx, client.ObjectKey{Namespace: "koth", Name: podName}, pod)
				if err != nil { // Fix: Added formatting directive for the error message
					continue
				}

				creationTime := pod.GetCreationTimestamp()
				timeSinceCreation := time.Since(creationTime.Time)
				minutesSinceCreation := timeSinceCreation.Minutes()
				score := int(minutesSinceCreation) * 10
				err = r.writePodScoreToFile(podName, score, ctx)
				if err != nil {
					return ctrl.Result{}, err
				}
				if err := r.deletepod(ctx, podName); err != nil {
					return ctrl.Result{}, err
				}

			}
		}
	}
}

func (r *CheckReconciler) writePodScoreToFile(podName string, score int, ctx context.Context) error {
	Name := fmt.Sprintf("/tmp/scores/%s.txt", podName)

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName + "-check",
			Namespace: "koth",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "check",
					Image: "youffes/check",
					Env: []corev1.EnvVar{
						{
							Name:  "POD_NAME",
							Value: podName,
						},
						{
							Name:  "SCORE",
							Value: strconv.Itoa(score),
						},
						{
							Name:  "FILENAME",
							Value: Name,
						},
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "scores",
							MountPath: "/tmp/scores",
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "scores",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/tmp/scores",
							Type: new(corev1.HostPathType),
						},
					},
				},
			},
		},
	}
	err := r.Create(ctx, pod)
	if err != nil {
		return err
	}
	return nil
}

func (r *CheckReconciler) deletepod(ctx context.Context, podName string) error {
	pod := corev1.Pod{}
	podname := types.NamespacedName{
		Name:      podName + "-check",
		Namespace: "default",
	}
	if err := r.Get(ctx, podname, &pod); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
	}
	if err := r.Delete(ctx, &pod); err != nil {
		return err
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CheckReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		Complete(r)
}
