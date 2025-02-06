package interfaces

import "fmt"

type PacketData struct {
	SenderServiceGroup string `json:"senderServiceGroup"`
	SenderServiceName  string `json:"senderServiceName"`
	SenderServiceId    string `json:"senderServiceId"`
}

func NewPacketData() PacketData {
	return PacketData{
		SenderServiceGroup: "netmanager",
		SenderServiceName:  "netmanager",
		SenderServiceId:    "*",
	}
}

func GetListenerChannel(packetType string) string {
	return fmt.Sprintf("netmanager:netmanager:*:%s", packetType)
}
