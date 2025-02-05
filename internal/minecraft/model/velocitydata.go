package model

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type VelocityData struct {
	groupName      string
	version        string
	build          int
	port           int
	replicasAmount int
	mongodbURI     string
	redisURI       string
}

func NewVelocityData(groupName string, version string, build int, port int, replicasAmount int, mongodbURI string, redisURI string) *VelocityData {
	return &VelocityData{
		groupName:      groupName,
		version:        version,
		build:          build,
		port:           port,
		replicasAmount: replicasAmount,
		mongodbURI:     mongodbURI,
		redisURI:       redisURI,
	}
}

func (data *VelocityData) GroupName() string {
	return data.groupName
}

func (data *VelocityData) Version() string {
	return data.version
}

func (data *VelocityData) BuildNumber() int {
	return data.build
}

func (data *VelocityData) Port() int {
	return data.port
}

func (data *VelocityData) ReplicasAmount() int {
	return data.replicasAmount
}

func (data *VelocityData) MongoDBURI() string {
	return data.mongodbURI
}

func (data *VelocityData) RedisURI() string {
	return data.redisURI
}

func (data *VelocityData) Build(serviceModel interfaces.ServiceModel, printer interfaces.Printer, imageManager interfaces.ImageManager, serviceManager interfaces.ServiceManager) {
	imageManager.FullDeployImage(serviceModel, serviceManager)
}

func (data *VelocityData) Update(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
}

func (data *VelocityData) Stop(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	wasStatefulSetDeloyed, wasServiceDeployed := false, false

	clusterManager.DeleteStatefulSet(serviceModel.Name())
	clusterManager.DeleteService(serviceModel.Name() + "-service")

	if !wasStatefulSetDeloyed || !wasServiceDeployed {
		printer.Print("Service "+serviceModel.Name()+" stopped.", printer.Service())
	}
}

func (data *VelocityData) Start(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {

}

func (data *VelocityData) Deploy(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {

}

func (data *VelocityData) generateStatefulSet(serviceModel interfaces.ServiceModel) *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceModel.Name(),
			Labels: map[string]string{
				"app": serviceModel.Name(),
			},
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: serviceModel.Name() + "-service",
			Replicas:    util.IntToPtr(data.replicasAmount),
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
							Image: "registry.dreammc.pl/dreammc/" + serviceModel.ImageName() + ":" + serviceModel.CurrentVersion(),
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: int32(25565),
								},
							},
							ImagePullPolicy: "Always",
							Env: []corev1.EnvVar{
								{
									Name:  "MONGODB_URI",
									Value: data.mongodbURI,
								},
								{
									Name:  "REDIS_URI",
									Value: data.redisURI,
								},
								{
									Name:  "GROUP_NAME",
									Value: data.groupName,
								},
								{
									Name:  "SERVICE_NAME",
									Value: serviceModel.Name(),
								},
								{
									Name: "SERVER_ID",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
							},
						},
					},
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: "harbor-credentials-secret",
						},
					},
				},
			},
		},
	}
}

func (data *VelocityData) generateService(serviceModel interfaces.ServiceModel) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceModel.Name() + "-service",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": serviceModel.Name()},
			Ports: []corev1.ServicePort{
				{
					Port:       25565,
					TargetPort: intstr.IntOrString{IntVal: 25565},
					NodePort:   int32(data.port),
					Protocol:   corev1.ProtocolTCP,
				},
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}
}
