package model

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	"encoding/base64"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HarborData struct {
	DomainField      string
	ProjectNameField string
	UsernameField    string
	PasswordField    string
}

func NewHarborData(domain string, projectName string, username string, password string) *HarborData {
	return &HarborData{
		DomainField:      domain,
		ProjectNameField: projectName,
		UsernameField:    username,
		PasswordField:    password,
	}
}

func (data *HarborData) Domain() string {
	return data.DomainField
}

func (data *HarborData) ProjectName() string {
	return data.ProjectNameField
}

func (data *HarborData) Username() string {
	return data.UsernameField
}

func (data *HarborData) Password() string {
	return data.PasswordField
}

func (data *HarborData) Build(serviceModel interfaces.ServiceModel, printer interfaces.Printer, imageManager interfaces.ImageManager, serviceManager interfaces.ServiceManager) {
	printer.PrintColored("Harbor service cannot be built.", printer.Service(), types.Yellow)
}

func (data *HarborData) Update(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	printer.PrintColored("Harbor service cannot be updated.", printer.Service(), types.Yellow)
}

func (data *HarborData) Stop(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	printer.PrintColored("Harbor service cannot be stopped.", printer.Service(), types.Yellow)
}

func (data *HarborData) Start(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	printer.PrintColored("Harbor service cannot be started.", printer.Service(), types.Yellow)
}

func (data *HarborData) Deploy(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	_, err := clusterManager.GetSecretOrErr("harbor-credentials-secret")
	if err == nil {
		return
	}

	clusterManager.CreateSecret(data.createCredentialsSecret())
	printer.Print("Created Harbor credentials secret.", printer.Service())
}

func (data *HarborData) createCredentialsSecret() *corev1.Secret {
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", data.UsernameField, data.PasswordField)))
	dockerConfigJson := fmt.Sprintf(`{
		"auths": {
			"%s": {
				"username": "%s",
				"password": "%s",
				"email": "%s",
				"auth": "%s"
			}
		}
	}`, data.DomainField, data.UsernameField, data.PasswordField, "netmanager@dreammc.pl", auth)

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "harbor-credentials-secret",
		},
		Type: corev1.SecretTypeDockerConfigJson,
		Data: map[string][]byte{
			".dockerconfigjson": []byte(dockerConfigJson),
		},
	}
}
