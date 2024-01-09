# IEC102

iec102主站golang实现

## 例

```go
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

```


## 任务接口

```go
type myClient struct{}

// Task 数据处理任务
func (c *myClient) Datahandler(data *IEC.APDU) error {
	// TODO 自定义数据处理
	println("do task")
	return nil
}
```



## 特性

- 连接不上3秒后重连
- 快速编码,解码
- interface设计,提供扩展性
- 简单的丰富的API



## 实现功能：

- 支持读取电能量

- 支持读取需量

- 支持读取自定义地址的遥测值

- 支持读取脉冲表电能量

- 支持读取智能表电能量

- 支持读取最大需量

- 支持读取瞬时量

- 支持读取电压合格率

