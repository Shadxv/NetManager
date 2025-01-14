package model

type HarborData struct {
	httpPort     int
	projectName  string
	username     string
	userMail     string
	userPassword string
	userRole     string
	disableGuest bool
}

func NewHarborData(httpPort int, projectName string, username string, userMail string, userPassword string, userRole string, disableGuest bool) *HarborData {
	return &HarborData{
		httpPort: httpPort,
		projectName: projectName,
		username: username,
		userMail: userMail,
		userPassword: userPassword,
		userRole: userRole,
		disableGuest: disableGuest,
	}
}

func (data *HarborData) HttpPort() int {
	return data.httpPort
}

func (data *HarborData) ProjectName() string {
	return data.projectName
}

func (data *HarborData) Username() string {
	return data.username
}

func (data *HarborData) UserMail() string {
	return data.userMail
}

func (data *HarborData) UserPassword() string {
	return data.userPassword
}

func (data *HarborData) UserRole() string {
	return data.userRole
}

func (data *HarborData) DisableGuest() bool {
	return data.disableGuest
}

