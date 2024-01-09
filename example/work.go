package main

import (
	"github.com/mrzwzw/IEC"
)

type myClient struct{}

// Task 数据处理任务
func (c *myClient) Datahandler(data *IEC.APDU) error {
	// TODO 自定义数据处理
	println("do task")
	return nil
}
