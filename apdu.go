package IEC

import "fmt"

// APDU 104数据包
type APDU struct {
	ControlDomains byte
	Address        *Address
	ASDU           *ASDU
	Len            int
	// ASDULen        int
	// CtrType        byte
	// CtrFrame       interface{}
	Signals []*Signal
}
type Address struct {
	LowAddress  int
	HighAddress int
}

// parseAPDU 解析APDU
func (apdu *APDU) parseAPDU(input []byte) error {
	// 这里的数字还需要修改
	if input == nil || len(input) < 10 {
		return fmt.Errorf("APDU报文[%X]非法", input)
	}
	Address := &Address{
		LowAddress:  int(input[1]),
		HighAddress: int(input[2]),
	}

	apdu.ControlDomains = input[0]
	asdu := new(ASDU)

	signals, err := asdu.ParseASDU(input[3:])
	if err != nil {
		return fmt.Errorf("APDU报文[%X]解析ASDU域[%X]异常: %v", input, input[3:], err)
	}

	apdu.Address = Address
	apdu.ASDU = asdu
	// apdu.CtrType = fType
	// apdu.CtrFrame = ctrFrame
	apdu.Signals = signals
	return nil
}
