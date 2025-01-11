package config

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GeneratePaperDeployment(serviceName string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName + "-deployment",
			Labels: map[string]string{
				"app": serviceName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: intToPtr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": serviceName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": serviceName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  serviceName + "-container",
							Image: "shadxvv/dreammc:latest",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: int32(25565),
								},
							},
							ImagePullPolicy: "Always",
						},
					},
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: "dockerhub-secret",
						},
					},
				},
			},
		},
	}
}
