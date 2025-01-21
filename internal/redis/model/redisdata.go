package model

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	"NetManager/pkg/util"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
	"time"
)

type RedisData struct {
	port            int
	externalPort    int
	password        string
	maxMemory       string
	maxMemoryPolicy string
	internalRedisIp string
	externalRedisIp string
}

func NewRedisData(port int, externalPort int, password string, maxMemory string, maxMemoryPolicy string) *RedisData {
	return &RedisData{
		port:            port,
		externalPort:    externalPort,
		password:        password,
		maxMemory:       maxMemory,
		maxMemoryPolicy: maxMemoryPolicy,
	}
}

func (data *RedisData) Port() int {
	return data.port
}

func (data *RedisData) ExternalPort() int {
	return data.externalPort
}

func (data *RedisData) Password() string {
	return data.password
}

func (data *RedisData) MaxMemory() string {
	return data.maxMemory
}

func (data *RedisData) MaxMemoryPolicy() string {
	return data.maxMemoryPolicy
}

func (data *RedisData) InternalRedisIp() string {
	return data.internalRedisIp
}

func (data *RedisData) ExternalRedisIp() string {
	return data.externalRedisIp
}

func (data *RedisData) Build(serviceModel interfaces.ServiceModel, printer interfaces.Printer, imageManager interfaces.ImageManager, serviceManager interfaces.ServiceManager) {
	printer.PrintColored("Redis service cannot be built.", printer.Service(), types.Yellow)
}

func (data *RedisData) Update(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	printer.PrintColored("Redis service cannot be updated.", printer.Service(), types.Yellow)
}

func (data *RedisData) Stop(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	printer.Print("Stopping Redis service...", printer.Service())
	clusterManager.DeleteDeployment(serviceModel.Name())
	clusterManager.DeleteConfigMap(serviceModel.Name() + "-config")
	clusterManager.DeleteService(serviceModel.Name() + "-internal")
	clusterManager.DeleteService(serviceModel.Name() + "-external")
	printer.Print("Redis service has been stopped", printer.Service())
}

func (data *RedisData) Start(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	printer.PrintColored("Redis service cannot be started.", printer.Service(), types.Yellow)
}

func (data *RedisData) Deploy(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	printer.Print("Checking Redis services statuses...", printer.Service())
	_, err := clusterManager.GetConfigMapOrErr(serviceModel.Name() + "-config")
	if err != nil {
		clusterManager.CreateConfigMap(data.generateConfigMap(serviceModel))
	}

	_, err = clusterManager.GetDeploymentOrErr(serviceModel.Name())
	if err != nil {
		clusterManager.CreateDeployment(data.generateDeployment(serviceModel))
	}

	_, err = clusterManager.GetServiceOrErr(serviceModel.Name() + "-internal")
	if err != nil {
		clusterManager.CreateService(data.generateClusterIPService(serviceModel))
	}

	_, err = clusterManager.GetServiceOrErr(serviceModel.Name() + "-external")
	if err != nil {
		clusterManager.CreateService(data.generateNodePortService(serviceModel))
	}

	var internalService *corev1.Service
	if internalService, _, err = data.waitForServicesReady(serviceModel, clusterManager); err != nil {
		printer.PrintColored("Error occured during waiting for Redis services.", printer.Service(), types.Red)
		printer.PrintColored(err.Error(), printer.Service(), types.Red)
		return
	}

	if data.internalRedisIp = internalService.Spec.ClusterIP; data.internalRedisIp == "" {
		printer.PrintColored("NetManager could not get Redis internal IP.", printer.Service(), types.Red)
		return
	}

	nodes, err := clusterManager.GetNodes()
	if err != nil {
		printer.PrintColored("NetManager could not get Redis external IP.", printer.Service(), types.Red)
		return
	}
	for _, node := range nodes {
		for _, addr := range node.Status.Addresses {
			if addr.Type == corev1.NodeExternalIP {
				data.externalRedisIp = addr.Address
				break
			} else if addr.Type == corev1.NodeInternalIP {
				data.externalRedisIp = addr.Address
				break
			}
		}
		if data.externalRedisIp != "" {
			break
		}
	}

	if data.externalRedisIp == "" {
		printer.PrintColored("Could not find any node IP address", printer.Service(), types.Red)
		return
	}

	serviceModel.SetStatus(types.Running)
	printer.Print("Internal address: "+data.internalRedisIp+":"+strconv.Itoa(data.port), printer.Service())
	printer.Print("External address: "+data.externalRedisIp+":"+strconv.Itoa(data.externalPort), printer.Service())
	printer.Print("Redis service is running", printer.Service())
}

func (data *RedisData) generateDeployment(serviceModel interfaces.ServiceModel) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceModel.Name(),
			Labels: map[string]string{
				"app": serviceModel.Name(),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: util.IntToPtr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": serviceModel.Name(),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": serviceModel.Name(),
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  serviceModel.Name() + "-container",
							Image: serviceModel.ImageName() + ":" + serviceModel.CurrentVersion(),
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: int32(data.port),
								},
							},
							Command: []string{
								"redis-server",
								"/etc/redis/redis.conf",
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "redis-volume",
									MountPath: "/etc/redis",
								},
							},
							ImagePullPolicy: "Always",
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "redis-volume",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: serviceModel.Name() + "-config",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (data *RedisData) generateClusterIPService(serviceModel interfaces.ServiceModel) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceModel.Name() + "-internal",
			Labels: map[string]string{
				"app": serviceModel.Name(),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": serviceModel.Name(),
			},
			Ports: []corev1.ServicePort{
				{
					Port: int32(data.port),
					TargetPort: intstr.IntOrString{
						IntVal: int32(data.port),
					},
					Protocol: corev1.ProtocolTCP,
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
}

func (data *RedisData) generateNodePortService(serviceModel interfaces.ServiceModel) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceModel.Name() + "-external",
			Labels: map[string]string{
				"app": serviceModel.Name(),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": serviceModel.Name(),
			},
			Ports: []corev1.ServicePort{
				{
					Port: int32(data.port),
					TargetPort: intstr.IntOrString{
						IntVal: int32(data.port),
					},
					NodePort: int32(data.externalPort),
					Protocol: corev1.ProtocolTCP,
				},
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}
}

func (data *RedisData) generateConfigMap(serviceModel interfaces.ServiceModel) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceModel.Name() + "-config",
			Labels: map[string]string{
				"app": serviceModel.Name(),
			},
		},
		Data: map[string]string{
			"redis.conf": fmt.Sprintf(`maxmemory 256mb
maxmemory-policy allkeys-lru
port %d
bind 0.0.0.0
appendonly yes
appendfsync everysec
timeout 0
tcp-keepalive 300
maxclients 10000`, data.port),
		},
	}
}

func (data *RedisData) waitForServicesReady(serviceModel interfaces.ServiceModel, clusterManager interfaces.ClusterManager) (*corev1.Service, *corev1.Service, error) {
	timeout := time.After(5 * time.Minute)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return nil, nil, fmt.Errorf("timeout waiting for Redis services to be ready")
		case <-ticker.C:
			internalService, internalReady, err := data.isServiceReady(serviceModel.Name()+"-internal", clusterManager)
			if err != nil {
				return nil, nil, fmt.Errorf("error checking internal service: %w", err)
			}

			externalService, externalReady, err := data.isServiceReady(serviceModel.Name()+"-external", clusterManager)
			if err != nil {
				return nil, nil, fmt.Errorf("error checking external service: %w", err)
			}

			deploymentReady, err := data.isDeploymentReady(serviceModel.Name(), clusterManager)
			if err != nil {
				return nil, nil, fmt.Errorf("error checking deployment: %w", err)
			}

			if internalReady && externalReady && deploymentReady {
				return internalService, externalService, nil
			}
		}
	}
}

func (data *RedisData) isServiceReady(name string, clusterManager interfaces.ClusterManager) (*corev1.Service, bool, error) {
	service, err := clusterManager.GetServiceOrErr(name)
	if err != nil {
		return nil, false, err
	}

	if service.Spec.ClusterIP == "" {
		return nil, false, nil
	}

	if service.Spec.Type == corev1.ServiceTypeNodePort {
		for _, port := range service.Spec.Ports {
			if port.NodePort == 0 {
				return nil, false, nil
			}
		}
	}

	return service, true, nil
}

func (data *RedisData) isDeploymentReady(name string, clusterManager interfaces.ClusterManager) (bool, error) {
	deployment, err := clusterManager.GetDeploymentOrErr(name)
	if err != nil {
		return false, err
	}

	if deployment.Status.ReadyReplicas != *deployment.Spec.Replicas {
		return false, nil
	}

	return data.arePodsReady(name, clusterManager)
}

func (data *RedisData) arePodsReady(name string, clusterManager interfaces.ClusterManager) (bool, error) {
	pods, err := clusterManager.GetPodsOrErr("app=" + name)
	if err != nil {
		return false, err
	}

	for _, pod := range pods {
		if pod.Status.Phase != corev1.PodRunning {
			return false, nil
		}

		for _, containerStatus := range pod.Status.ContainerStatuses {
			if !containerStatus.Ready {
				return false, nil
			}
		}
	}

	return true, nil
}
