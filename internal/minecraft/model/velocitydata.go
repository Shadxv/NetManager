package model

type VelocityData struct {
	version        string
	build          int
	port           int
	replicasAmount int
}

func NewVelocityData(version string, build int, port int, replicasAmount int) *VelocityData {
	return &VelocityData{
		version:        version,
		build:          build,
		port:           port,
		replicasAmount: replicasAmount,
	}
}

func (data *VelocityData) Version() string {
	return data.version
}

func (data *VelocityData) Build() int {
	return data.build
}

func (data *VelocityData) Port() int {
	return data.port
}

func (data *VelocityData) ReplicasAmount() int {
	return data.replicasAmount
}
