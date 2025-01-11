package service

type Model interface {
	Name() string
	ServiceType() string
	Status() string
	SetStatus(status string)
	ImageName() string
	CurrentVersion() string
	SetCurrentVersion(currentVersion string) bool
	AvailableVersions() []string
	AddVersion(version string)
	ServiceData() *interface{}
}
