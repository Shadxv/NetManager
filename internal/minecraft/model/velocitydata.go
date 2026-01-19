package model

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	"NetManager/pkg/util"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type VelocityData struct {
	GroupNameField      string
	VersionField        string
	BuildField          int
	PortField           int
	ReplicasAmountField int
	MongodbURIField     string
	RedisURIField       string
}

func NewVelocityData(groupName string, version string, build int, port int, replicasAmount int, mongodbURI string, redisURI string) *VelocityData {
	return &VelocityData{
		GroupNameField:      groupName,
		VersionField:        version,
		BuildField:          build,
		PortField:           port,
		ReplicasAmountField: replicasAmount,
		MongodbURIField:     mongodbURI,
		RedisURIField:       redisURI,
	}
}

func (data *VelocityData) GroupName() string {
	return data.GroupNameField
}

func (data *VelocityData) Version() string {
	return data.VersionField
}

func (data *VelocityData) BuildNumber() int {
	return data.BuildField
}

func (data *VelocityData) Port() int {
	return data.PortField
}

func (data *VelocityData) ReplicasAmount() int {
	return data.ReplicasAmountField
}

func (data *VelocityData) MongoDBURI() string {
	return data.MongodbURIField
}

func (data *VelocityData) RedisURI() string {
	return data.RedisURIField
}

func (data *VelocityData) Build(serviceModel interfaces.ServiceModel, printer interfaces.Printer, imageManager interfaces.ImageManager, serviceManager interfaces.ServiceManager) {
	imageManager.FullDeployImage(serviceModel, serviceManager)
}

func (data *VelocityData) Update(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
}

func (data *VelocityData) Stop(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	wasStatefulSetDeloyed, wasServiceDeployed := false, false

	clusterManager.DeleteStatefulSet(serviceModel.Name(), serviceModel.Namespace())
	clusterManager.DeleteService(serviceModel.Name()+"-service", serviceModel.Namespace())

	if !wasStatefulSetDeloyed || !wasServiceDeployed {
		printer.Print("Service "+serviceModel.Name()+" stopped.", printer.Service())
	}
}

func (data *VelocityData) Start(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {

}

func (data *VelocityData) Deploy(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	_, err := clusterManager.GetStatefulSetOrErr(serviceModel.Name(), serviceModel.Namespace())
	if err == nil {
		printer.PrintColored("StatefulSet has already been deployed.", printer.Service(), types.Yellow)
	} else {
		clusterManager.CreateStatefulSet(data.generateStatefulSet(serviceModel), serviceModel.Namespace())
	}

	_, err = clusterManager.GetServiceOrErr(serviceModel.Name()+"-service", serviceModel.Namespace())
	if err == nil {
		printer.PrintColored("Service has already been deployed.", printer.Service(), types.Yellow)
	} else {
		clusterManager.CreateService(data.generateService(serviceModel), serviceModel.Namespace())
	}
}

func (data *VelocityData) generateStatefulSet(serviceModel interfaces.ServiceModel) *appsv1.StatefulSet {
	statefulset := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceModel.Name(),
			Labels: map[string]string{
				"app": serviceModel.Name(),
			},
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: serviceModel.Name() + "-service",
			Replicas:    util.IntToPtr(data.ReplicasAmountField),
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
									Value: data.MongodbURIField,
								},
								{
									Name:  "REDIS_URI",
									Value: data.RedisURIField,
								},
								{
									Name:  "GROUP_NAME",
									Value: "dreammc", // Later change to config value
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

	memory, err := resource.ParseQuantity("4Gi")
	if err != nil {
		return statefulset
	}

	statefulset.Spec.Template.Spec.Containers[0].Resources = corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceMemory: memory,
		},
	}

	return statefulset
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
					NodePort:   int32(data.PortField),
					Protocol:   corev1.ProtocolTCP,
				},
			},
			Type:                  corev1.ServiceTypeNodePort,
			ExternalTrafficPolicy: corev1.ServiceExternalTrafficPolicyTypeLocal,
		},
	}
}
