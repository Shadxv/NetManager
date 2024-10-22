package config

import (
	"NetManager/internal/config/model"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func generateRedisConfigMap(config *model.RedisConfig) *corev1.ConfigMap {
	data := map[string]string{
		"maxmemory":        config.MaxMemory,
		"maxmemory-policy": config.MaxMemoryPolicy,
	}

	if config.Password != "" {
		data["requirepass"] = config.Password
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-redis-config",
		},
		Data: data,
	}
}

func generateRedisDeployment(config *model.RedisConfig) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "redis-deployment",
			Labels: map[string]string{
				"app": "redis",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: intToPtr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "redis",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "redis",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "redis-container",
							Image: fmt.Sprintf("%s:%s", config.DockerImage, config.Version),
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: int32(config.Port),
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU: resource.MustParse("1m"),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "data",
									MountPath: "/redis-master-data",
								},
								{
									Name:      "redis-config",
									MountPath: "/redis-master",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "redis-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "my-redis-config",
									},
								},
							},
						},
						{
							Name: "data",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}
}

func GenerateRedisConfig(config *model.RedisConfig) (*corev1.ConfigMap, *appsv1.Deployment) {
	return generateRedisConfigMap(config), generateRedisDeployment(config)
}

func intToPtr(i int32) *int32 {
	return &i
}
