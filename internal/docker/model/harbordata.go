package model

type HarborData struct {
	domain      string
	projectName string
	username    string
	password    string
}

func NewHarborData(domain string, projectName string, username string, password string) *HarborData {
	return &HarborData{
		domain:      domain,
		projectName: projectName,
		username:    username,
		password:    password,
	}
}

func (data *HarborData) Domain() string {
	return data.domain
}

func (data *HarborData) ProjectName() string {
	return data.projectName
}

func (data *HarborData) Username() string {
	return data.username
}

func (data *HarborData) Password() string {
	return data.password
}
