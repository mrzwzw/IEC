package main

import (
	"github.com/mrzwzw/IEC"

	"github.com/sirupsen/logrus"
)

func main() {
	p := IEC.NewRTUClientProvider()

	p.Address = "/dev/ttyUSB0"
	p.BaudRate = 9600
	p.DataBits = 8
	p.Parity = "N"
	p.StopBits = 1

	var logger *logrus.Logger
	myclli := &myClient{}
	client := IEC.NewClient(myclli, logger, p)

	client.Run()
	// 发送链路复位帧
	client.Reset()

	// 发送确定帧
	client.Resetframe()

	// 发送读取指定地址范围内的遥测量
	var start byte = 0x11
	var end byte = 0x13
	client.SendYC(start, end)
}
