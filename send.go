package IEC

// 复位链路
func (c *Client) Reset() {
	data := []byte{0x10, 0x40, 0x01, 0x00, 0x41, 0x16}
	c.sendChan <- data
}

// 确定帧
func (c *Client) Resetframe() {
	data := []byte{0x10, 0x53, 0x01, 0x00, 0x54, 0x16}
	c.sendChan <- data
}

// sendTotalCall 读取指定地址范围内的遥测量
func (c *Client) SendYC(a, b byte) {
	data1 := []byte{0x73, 0x01, 0x00, 0xAB, 0x01, 0x06, 0x01, 0x00, 0x00, a, b}
	data2 := addMod(data1[0:3], 256)

	data := convertBytes(data1)
	data = append(data, data2, 0x16)
	c.Logger.Debugf("读取指定地址范围内的遥测量: [% X]", data)
}

// 标识:82 脉冲表电能量
func (c *Client) SendPulseElec() {
	data1 := []byte{0x7B, 0x01, 0x7C, 0x00, 0x05, 0x01, 0x00, 0x82, 0x00, 0x00, 0x01, 0x00}
	data2 := addMod(data1[0:3], 256)

	data := convertBytes(data1)
	data = append(data, data2, 0x16)
	c.Logger.Debugf("读取指定地址范围内的遥测量: [% X]", data)
}

// 标识:83 智能表电能值
func (c *Client) SendSmartElec() {
	data1 := []byte{0x7B, 0x01, 0x7C, 0x00, 0x05, 0x01, 0x00, 0x83, 0x00, 0x00, 0x01, 0x00}
	data2 := addMod(data1[0:3], 256)

	data := convertBytes(data1)
	data = append(data, data2, 0x16)
	c.Logger.Debugf("读取指定地址范围内的遥测量: [% X]", data)
}

// 标识:84 最大需量
func (c *Client) SendMaxDemand() {
	data1 := []byte{0x7B, 0x01, 0x7C, 0x00, 0x05, 0x01, 0x00, 0x84, 0x00, 0x00, 0x01, 0x00}
	data2 := addMod(data1[0:3], 256)

	data := convertBytes(data1)
	data = append(data, data2, 0x16)
	c.Logger.Debugf("读取指定地址范围内的遥测量: [% X]", data)
}

// 标识:85 瞬时值
func (c *Client) SendInstant() {
	data1 := []byte{0x7B, 0x01, 0x7C, 0x00, 0x05, 0x01, 0x00, 0x85, 0x00, 0x00, 0x01, 0x00}
	data2 := addMod(data1[0:3], 256)

	data := convertBytes(data1)
	data = append(data, data2, 0x16)
	c.Logger.Debugf("读取指定地址范围内的遥测量: [% X]", data)
}

// 标识:88 电压合格率
func (c *Client) SendVoltagePassRate() {
	data1 := []byte{0x7B, 0x01, 0x7C, 0x00, 0x05, 0x01, 0x00, 0x88, 0x00, 0x00, 0x01, 0x00}
	data2 := addMod(data1[0:3], 256)

	data := convertBytes(data1)
	data = append(data, data2, 0x16)
	c.Logger.Debugf("读取指定地址范围内的遥测量: [% X]", data)
}
