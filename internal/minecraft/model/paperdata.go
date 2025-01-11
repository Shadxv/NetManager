package model

type PaperData struct {
	version     string
	build       int
	minReplicas int
}

func NewPaperData(version string, build int, minReplicas int) *PaperData {
	return &PaperData{
		version:     version,
		build:       build,
		minReplicas: minReplicas,
	}
}

func (data *PaperData) Version() string {
	return data.version
}

func (data *PaperData) Build() int {
	return data.build
}

func (data *PaperData) MinReplicas() int {
	return data.minReplicas
}
