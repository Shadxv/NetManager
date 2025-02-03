package model

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	"NetManager/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type PaperData struct {
	groupName   string
	version     string
	build       int
	minReplicas int
	mongodbURI  string
	redisURI    string
}

func NewPaperData(groupName string, version string, build int, minReplicas int, mongodbURI string, redisURI string) *PaperData {
	return &PaperData{
		groupName:   groupName,
		version:     version,
		build:       build,
		minReplicas: minReplicas,
		mongodbURI:  mongodbURI,
		redisURI:    redisURI,
	}
}
func (data *PaperData) GroupName() string {
	return data.groupName
}

func (data *PaperData) Version() string {
	return data.version
}

func (data *PaperData) BuildNumber() int {
	return data.build
}

func (data *PaperData) MinReplicas() int {
	return data.minReplicas
}

func (data *PaperData) MongoDBURI() string {
	return data.mongodbURI
}

func (data *PaperData) RedisURI() string {
	return data.redisURI
}

func (data *PaperData) Build(serviceModel interfaces.ServiceModel, printer interfaces.Printer, imageManager interfaces.ImageManager, serviceManager interfaces.ServiceManager) {
	imageManager.FullDeployImage(serviceModel, serviceManager)
}

func (data *PaperData) Update(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {

}

func (data *PaperData) Stop(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	wasStatefulSetDeloyed, wasServiceDeployed := false, false

	clusterManager.DeleteStatefulSet(serviceModel.Name())
	clusterManager.DeleteService(serviceModel.Name() + "-service")

	if !wasStatefulSetDeloyed || !wasServiceDeployed {
		printer.Print("Service "+serviceModel.Name()+" stopped.", printer.Service())
	}
}

func (data *PaperData) Start(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {

}

func (data *PaperData) Deploy(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	_, err := clusterManager.GetStatefulSetOrErr(serviceModel.Name())
	if err == nil {
		printer.PrintColored("StatefulSet has already been deployed.", printer.Service(), types.Yellow)
	} else {
		clusterManager.CreateStatefulSet(data.generateStatefulSet(serviceModel))
	}

	_, err = clusterManager.GetServiceOrErr(serviceModel.Name() + "-service")
	if err == nil {
		printer.PrintColored("Service has already been deployed.", printer.Service(), types.Yellow)
	} else {
		clusterManager.CreateService(data.generateService(serviceModel))
	}
}

func (data *PaperData) generateStatefulSet(serviceModel interfaces.ServiceModel) *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceModel.Name(),
			Labels: map[string]string{
				"app": serviceModel.Name(),
			},
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: serviceModel.Name() + "-service",
			Replicas:    util.IntToPtr(data.minReplicas),
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

func (data *PaperData) generateService(serviceModel interfaces.ServiceModel) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceModel.Name() + "-service",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": serviceModel.Name()},
			Ports: []corev1.ServicePort{
				{
					Port:       45565,
					TargetPort: intstr.IntOrString{IntVal: 25565},
					Protocol:   corev1.ProtocolTCP,
				},
			},
			Type: corev1.ServiceTypeLoadBalancer,
		},
	}
}
