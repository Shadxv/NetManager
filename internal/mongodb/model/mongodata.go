package model

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	"encoding/base64"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"strconv"
	"time"
)

type MongoData struct {
	port            int
	externalPort    int
	username        string
	password        string
	authRequired    bool
	authorization   bool
	internalMongoIp string
	externalMongoIp string
	internalURI     string
	externalURI     string
}

func NewMongoData(port int, externalPort int, username string, password string, authRequired bool, authorization bool) *MongoData {
	return &MongoData{
		port:          port,
		externalPort:  externalPort,
		username:      username,
		password:      password,
		authRequired:  authRequired,
		authorization: authorization,
	}
}

func (data *MongoData) Port() int {
	return data.port
}

func (data *MongoData) ExternalPort() int {
	return data.externalPort
}

func (data *MongoData) Username() string {
	return data.username
}

func (data *MongoData) Password() string {
	return data.password
}

func (data *MongoData) AuthRequired() bool {
	return data.authRequired
}

func (data *MongoData) Authorization() bool {
	return data.authorization
}

func (data *MongoData) InternalMongoIp() string {
	return data.internalMongoIp
}

func (data *MongoData) ExternalMongoIp() string {
	return data.externalMongoIp
}

func (data *MongoData) InternalURI() string {
	return data.internalURI
}

func (data *MongoData) ExternalURI() string {
	return data.externalURI
}

func (data *MongoData) buildURI(ip string, port int) string {
	auth := ""
	if data.authRequired {
		auth = fmt.Sprintf("%s:%s@", data.username, data.password)
	}

	uri := fmt.Sprintf("mongodb://%s%s:%d",
		auth,
		ip,
		port,
	)

	return uri
}

func (data *MongoData) Build(serviceModel interfaces.ServiceModel, printer interfaces.Printer, imageManager interfaces.ImageManager, serviceManager interfaces.ServiceManager) {
	printer.PrintColored("MongoDB service cannot be built.", printer.Service(), types.Yellow)
}

func (data *MongoData) Update(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	printer.PrintColored("MongoDB service cannot be updated.", printer.Service(), types.Yellow)
}

func (data *MongoData) Stop(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	printer.Print("Stopping MongoDB service...", printer.Service())
	clusterManager.DeleteStatefulSet(serviceModel.Name())
	clusterManager.DeleteConfigMap(serviceModel.Name() + "-config")
	clusterManager.DeleteSecret(serviceModel.Name() + "-root-credentials")
	clusterManager.DeleteService(serviceModel.Name() + "-internal")
	clusterManager.DeleteService(serviceModel.Name() + "-external")
	printer.Print("MongoDB service has been stopped", printer.Service())
}

func (data *MongoData) Start(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	printer.PrintColored("MongoDB service cannot be started.", printer.Service(), types.Yellow)
}

func (data *MongoData) Deploy(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	printer.Print("Checking MongoDB services statuses...", printer.Service())
	_, err := clusterManager.GetConfigMapOrErr(serviceModel.Name() + "-config")
	if err != nil {
		clusterManager.CreateConfigMap(data.generateConfigMap(serviceModel))
	}

	_, err = clusterManager.GetSecretOrErr(serviceModel.Name() + "-root-credentials")
	if err != nil {
		clusterManager.CreateSecret(data.generateCredentialsSecret(serviceModel))
	}

	_, err = clusterManager.GetServiceOrErr(serviceModel.Name() + "-internal")
	if err != nil {
		clusterManager.CreateService(data.generateClusterIPService(serviceModel))
	}

	_, err = clusterManager.GetServiceOrErr(serviceModel.Name() + "-external")
	if err != nil {
		clusterManager.CreateService(data.generateNodePortService(serviceModel))
	}

	_, err = clusterManager.GetStatefulSetOrErr(serviceModel.Name())
	if err != nil {
		clusterManager.CreateStatefulSet(data.generateStatefulSet(serviceModel))
	}

	var internalService *corev1.Service
	if internalService, _, err = data.waitForServicesReady(serviceModel, clusterManager); err != nil {
		printer.PrintColored("Error occured during waiting for MongoDB services.", printer.Service(), types.Red)
		printer.PrintColored(err.Error(), printer.Service(), types.Red)
		return
	}

	if data.internalMongoIp = internalService.Spec.ClusterIP; data.internalMongoIp == "" {
		printer.PrintColored("NetManager could not get MongoDB internal IP.", printer.Service(), types.Red)
		return
	}

	nodes, err := clusterManager.GetNodes()
	if err != nil {
		printer.PrintColored("NetManager could not get MongoDB external IP.", printer.Service(), types.Red)
		return
	}
	for _, node := range nodes {
		for _, addr := range node.Status.Addresses {
			if addr.Type == corev1.NodeExternalIP {
				data.externalMongoIp = addr.Address
				break
			} else if addr.Type == corev1.NodeInternalIP {
				data.externalMongoIp = addr.Address
				break
			}
		}
		if data.externalMongoIp != "" {
			break
		}
	}

	if data.externalMongoIp == "" {
		printer.PrintColored("Could not find any node IP address", printer.Service(), types.Red)
		return
	}

	serviceModel.SetStatus(types.Running)
	printer.Print("Internal address: "+data.internalMongoIp+":"+strconv.Itoa(data.port), printer.Service())
	printer.Print("External address: "+data.externalMongoIp+":"+strconv.Itoa(data.externalPort), printer.Service())
	data.internalURI = data.buildURI(data.internalMongoIp, data.port)
	data.externalURI = data.buildURI(data.externalMongoIp, data.externalPort)
	printer.Print("Internal URI: "+data.internalURI, printer.Service())
	printer.Print("External URI: "+data.externalURI, printer.Service())
	printer.Print("MongoDB service is running", printer.Service())

}

func (data *MongoData) generateCredentialsSecret(serviceModel interfaces.ServiceModel) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceModel.Name() + "-root-credentials",
			Labels: map[string]string{
				"app": serviceModel.Name(),
			},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"MONGO_INITDB_ROOT_USERNAME": []byte(base64.StdEncoding.EncodeToString([]byte(data.Username()))),
			"MONGO_INITDB_ROOT_PASSWORD": []byte(base64.StdEncoding.EncodeToString([]byte(data.Password()))),
		},
	}
}

func (data *MongoData) generateClusterIPService(serviceModel interfaces.ServiceModel) *corev1.Service {
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

func (data *MongoData) generateNodePortService(serviceModel interfaces.ServiceModel) *corev1.Service {
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

func (data *MongoData) generateConfigMap(serviceModel interfaces.ServiceModel) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceModel.Name() + "-config",
			Labels: map[string]string{
				"app": serviceModel.Name(),
			},
		},
		Data: map[string]string{
			"mongod.conf": fmt.Sprintf(`
net:
    port: %d
security:
    authorization: %t
    authentication: %t
`, data.port, data.authorization, data.authRequired),
			"init.js": `
db.createUser({
    user: process.env.MONGO_INITDB_ROOT_USERNAME,
    pwd: process.env.MONGO_INITDB_ROOT_PASSWORD,
    roles: [
        { role: "userAdminAnyDatabase", db: "admin" },
        { role: "readWriteAnyDatabase", db: "admin" }
    ]
})`,
		},
	}
}

func (data *MongoData) generateStatefulSet(serviceModel interfaces.ServiceModel) *appsv1.StatefulSet {
	labels := map[string]string{
		"app": serviceModel.Name(),
	}

	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:   serviceModel.Name(),
			Labels: labels,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: serviceModel.Name() + "-internal",
			Replicas:    pointer.Int32(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
					Annotations: map[string]string{
						"longhornio/disable-auto-mount": "false",
						"longhornio/node-selector":      "kubernetes.io/hostname: minikube",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "mongodb",
							Image: "mongo:latest",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: int32(data.Port()),
									Protocol:      corev1.ProtocolTCP,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "mongodb-data",
									MountPath: "/data/db",
									SubPath:   "mongodb",
								},
								{
									Name:      "mongodb-config",
									MountPath: "/etc/mongod.conf",
									SubPath:   "mongod.conf",
								},
								{
									Name:      "mongodb-config",
									MountPath: "/docker-entrypoint-initdb.d/init.js",
									SubPath:   "init.js",
								},
							},
							Command: []string{
								"mongod",
								"--config=/etc/mongod.conf",
							},
							EnvFrom: []corev1.EnvFromSource{
								{
									SecretRef: &corev1.SecretEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: serviceModel.Name() + "root-credentials",
										},
									},
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "mongodb-config",
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
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "mongodb-data",
						Annotations: map[string]string{
							"longhornio/disable-frontend":  "false",
							"longhornio/volume-scheduling": "true",
							"longhornio/node-selector":     "kubernetes.io/hostname: minikube",
							"longhornio/replica-count":     "1",
						},
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadWriteOnce,
						},
						StorageClassName: pointer.String("longhorn"),
						Resources: corev1.VolumeResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: resource.MustParse("5Gi"),
							},
						},
					},
				},
			},
		},
	}
}

func (data *MongoData) waitForServicesReady(serviceModel interfaces.ServiceModel, clusterManager interfaces.ClusterManager) (*corev1.Service, *corev1.Service, error) {
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

			statefulsetReady, err := data.isStatefulSetReady(serviceModel.Name(), clusterManager)
			if err != nil {
				return nil, nil, fmt.Errorf("error checking statefulset: %w", err)
			}

			if internalReady && externalReady && statefulsetReady {
				return internalService, externalService, nil
			}
		}
	}
}

func (data *MongoData) isServiceReady(name string, clusterManager interfaces.ClusterManager) (*corev1.Service, bool, error) {
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

func (data *MongoData) isStatefulSetReady(name string, clusterManager interfaces.ClusterManager) (bool, error) {
	statefulset, err := clusterManager.GetStatefulSetOrErr(name)
	if err != nil {
		return false, err
	}

	if statefulset.Status.ReadyReplicas != *statefulset.Spec.Replicas {
		return false, nil
	}

	return data.arePodsReady(name, clusterManager)
}

func (data *MongoData) arePodsReady(name string, clusterManager interfaces.ClusterManager) (bool, error) {
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
