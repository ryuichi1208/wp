/*
Copyright 2023.

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
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	wpv1 "github.com/ryuichi1208.wp/api/v1"
)

// WpReconciler reconciles a Wp object
type WpReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=wp.gurasan.github.io,resources=wps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=wp.gurasan.github.io,resources=wps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=wp.gurasan.github.io,resources=wps/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=events,verbs=create;update;patch

func (r *WpReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Log.Info("Reconciling Wp", "namespace", req.Namespace, "name", req.Name)
	_ = log.FromContext(ctx)

	var pod corev1.Pod

	wp := &wpv1.Wp{}
	err := r.Get(ctx, req.NamespacedName, wp)

	err = r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: "nginx-pod"}, &pod)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Log.Info("pod", "apiVersion", pod.APIVersion, "hostPID", pod.Spec.HostPID)

	r.updateStatus(wp)
	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: 10 * time.Second,
	}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *WpReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&wpv1.Wp{}).
		Complete(r)
}

func (r *WpReconciler) updateStatus(wp *wpv1.Wp) error {
	wp.Status.Ready = true
	return r.Status().Update(context.Background(), wp)
}
