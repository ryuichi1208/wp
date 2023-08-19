package controllers

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
}
