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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	if err != nil {
		return ctrl.Result{}, err
	}
	log.Log.Info("version!!!!!!!!", "ver", wp.Spec)

	err = r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: "nginx-pod"}, &pod)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Log.Info("pod", "apiVersion", pod.APIVersion, "hostPID", pod.Spec.HostPID)

	// TODO(user): your logic here

	// r.GetServiceListAndPrintServiceName()
	//r.CreateNginxPod()
	r.CreateDeploymentAndService()

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

func (r *WpReconciler) GetServiceListAndPrintServiceName() {
	var serviceList corev1.ServiceList
	err := r.List(context.Background(), &serviceList)
	if err != nil {
		log.Log.Error(err, "failed to get service list")
	}

	for _, service := range serviceList.Items {
		log.Log.Info("service name", "name", service.Name)
	}
}

func (r *WpReconciler) CreateNginxPod() {
	// nginx-pod2という名前のpodがあれば削除する
	var pod corev1.Pod
	err := r.Get(context.Background(), client.ObjectKey{Namespace: "default", Name: "nginx-pod2"}, &pod)
	if err == nil {
		err = r.Delete(context.Background(), &pod)
		if err != nil {
			log.Log.Error(err, "failed to delete nginx pod")

		}
	}
	// 削除したpodがなくなるまで待つ
	for {
		err = r.Get(context.Background(), client.ObjectKey{Namespace: "default", Name: "nginx-pod2"}, &pod)
		if err != nil {
			break
		}
	}

	pod = corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-pod2",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx:latest",
				},
			},
		},
	}
	err = r.Create(context.Background(), &pod)
	if err != nil {
		log.Log.Error(err, "failed to create nginx pod")
	}
}

func (r *WpReconciler) CreateDeploymentAndService() {
	log.Log.Info("CreateDeploymentAndService")
	// nginxのdeploymentとserviceを作成する
	// deploymentがあれば削除する
	var deployment appsv1.Deployment
	err := r.Get(context.Background(), client.ObjectKey{Namespace: "default", Name: "nginx-deployment"}, &deployment)
	if err == nil {
		err = r.Delete(context.Background(), &deployment)
		if err != nil {
			log.Log.Error(err, "failed to delete nginx deployment")
		}
	}
	// 削除したdeploymentがなくなるまで待つ
	for {
		err = r.Get(context.Background(), client.ObjectKey{Namespace: "default", Name: "nginx-deployment"}, &deployment)
		if err != nil {
			break
		}
	}
	// deploymentを作成する
	deployment = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-deployment",
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: func(i int32) *int32 { return &i }(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "nginx"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "nginx"},
				},
				// nginxのlatestイメージを使う
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:latest",
						},
					},
				},
			},
		},
	}
	err = r.Create(context.Background(), &deployment)
	if err != nil {
		log.Log.Error(err, "failed to create nginx deployment")
	}

	// seviceを作成する
	var service corev1.Service
	err = r.Get(context.Background(), client.ObjectKey{Namespace: "default", Name: "nginx-service"}, &service)
	if err == nil {
		err = r.Delete(context.Background(), &service)
		if err != nil {
			log.Log.Error(err, "failed to delete nginx service")
		}
	}
	// 削除したserviceがなくなるまで待つ
	for {
		err = r.Get(context.Background(), client.ObjectKey{Namespace: "default", Name: "nginx-service"}, &service)
		if err != nil {
			break
		}
	}
	// serviceを作成する
	service = corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-service",
			Namespace: "default",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": "nginx"},
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: 80,
				},
			},
		},
	}
	err = r.Create(context.Background(), &service)
	if err != nil {
		log.Log.Error(err, "failed to create nginx service")
	}
}
