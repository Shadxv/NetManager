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

type PaperData struct {
	GroupNameField   string
	VersionField     string
	BuildField       int
	MinReplicasField int
	MongodbURIField  string
	RedisURIField    string
}

func NewPaperData(groupName string, version string, build int, minReplicas int, mongodbURI string, redisURI string) *PaperData {
	return &PaperData{
		GroupNameField:   groupName,
		VersionField:     version,
		BuildField:       build,
		MinReplicasField: minReplicas,
		MongodbURIField:  mongodbURI,
		RedisURIField:    redisURI,
	}
}
func (data *PaperData) GroupName() string {
	return data.GroupNameField
}

func (data *PaperData) Version() string {
	return data.VersionField
}

func (data *PaperData) BuildNumber() int {
	return data.BuildField
}

func (data *PaperData) MinReplicas() int {
	return data.MinReplicasField
}

func (data *PaperData) MongoDBURI() string {
	return data.MongodbURIField
}

func (data *PaperData) RedisURI() string {
	return data.RedisURIField
}

func (data *PaperData) Build(serviceModel interfaces.ServiceModel, printer interfaces.Printer, imageManager interfaces.ImageManager, serviceManager interfaces.ServiceManager) {
	imageManager.FullDeployImage(serviceModel, serviceManager)
}

func (data *PaperData) Update(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {

}

func (data *PaperData) Stop(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	wasStatefulSetDeloyed, wasServiceDeployed := false, false

	clusterManager.DeleteStatefulSet(serviceModel.Name(), serviceModel.Namespace())
	clusterManager.DeleteService(serviceModel.Name()+"-service", serviceModel.Namespace())

	if !wasStatefulSetDeloyed || !wasServiceDeployed {
		printer.Print("Service "+serviceModel.Name()+" stopped.", printer.Service())
	}
}

func (data *PaperData) Start(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {

}

func (data *PaperData) Deploy(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
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

func (data *PaperData) generateStatefulSet(serviceModel interfaces.ServiceModel) *appsv1.StatefulSet {
	statefulset := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceModel.Name(),
			Labels: map[string]string{
				"app": serviceModel.Name(),
			},
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: serviceModel.Name() + "-service",
			Replicas:    util.IntToPtr(data.MinReplicasField),
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
									Value: "dreammc",
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

func (data *PaperData) generateService(serviceModel interfaces.ServiceModel) *corev1.Service {
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
					Protocol:   corev1.ProtocolTCP,
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
}
