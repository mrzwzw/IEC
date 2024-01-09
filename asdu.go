package IEC

import (
	"encoding/binary"
	"fmt"
	"log"
	"time"
)

// ASDU 应用服务数据单元
type ASDU struct {
	TypeID        byte    // 类型标识
	Sequence      bool    // 是否连续
	Length        byte    // 可变结构限定词
	Cause         uint16  // 传输原因
	PublicAddress uint16  // 公共地址
	RecordAddress uint    // 记录地址
	Ts            float64 // 毫秒级时间戳
}

// 数据类型
const (
	// MSpNa1 不带游标的单点遥信，3个字节的地址，1个字节的值
	MSpNa1 = 1
	// MDpNa1 不带时标的双点遥信，每个遥信占1个字节
	MDpNa1 = 3
	// MMeNc1 带品质描述的测量值，每个遥测值占3个字节
	MMeNa1 = 9
	// MMeNc1 带品质描述的浮点值，每个遥测值占5个字节
	MMeNc1 = 13
	// MItNa1 电度总量,每个遥脉值占5个字节
	// MItNa1 = 15
	// MSpTb1 带游标的单点遥信，3个字节的地址，1个字节的值，7个字节短时标
	MSpTb1 = 30
	// MEiNA1 初始化结束
	MEiNA1 = 70
	// CIcNa1 总召唤
	CIcNa1 = 100
	// CCiNa1 电度总召唤
	CCiNa1 = 101
)

// 子站侧数据类型
const (
	// 带时标的但单点信息 常用
	SUBPoint = 1
	// 积分电能量-表底值，4字节 常用
	SUBElecBottom = 2
	// 积分电能量-增量值，4字节
	SUBElecIncr = 5
	// 传送某几路电表的电量实时数据。p102
	SUBDBCur = 15
	// 初始化结束 常用
	SUBInitEnd = 70
	// 采集器的制造厂和产品规范
	SUBManuSpeci = 71
	// 采集器当前系统时间
	SUBNCurTime = 72
	// 时钟同步
	SUBClockSync = 128
	// 复费率积分电能量-表底值，4字节
	SUBCompRateCreditsBottom = 160
	// 遥测量当前值
	SUBCurYC = 161
	// 要测量历史值
	SUBHisYC = 162
)

// 主站侧数据类型
const (
	// 读制造厂和产品规范
	MAINManuSpeci = 100
	// 读带时标的单点信息的记录
	MAINPoint = 101
	// 读一个选定时间的带时标的单点信息的记录
	MAINSTPoint = 102
	// 读采集器的当前系统时间
	MAINCurTime = 104
	// 读最早累计时段的积分电能量一表底值
	MAINElecBottom = 105
	// 读选定时间范用、选定地址范用的积分电能量一表底值
	MAINSTElecBottom = 120
	// 读选定时间范围、选定地址范用的积分电能量一增量值
	MAINSTElecIncr = 121
	// 时钟同步
	MAINClockSync = 128
	// 读指定地址范围和时间范围的复费率积分电能量一表底值
	MAINSTCompRateCreditsBottom = 160
	// 读指定地址范围的遥测量当前值
	MAINCurYC = 171
)

// ParseASDU 解析asdu
func (asdu *ASDU) ParseASDU(asduBytes []byte) (signals []*Signal, err error) {
	signals = make([]*Signal, 0)
	// 不太理解为什么和4做判断
	if asduBytes == nil || len(asduBytes) < 4 {
		err = fmt.Errorf("asdu[%X]非法", asduBytes)
		return
	}
	asdu.TypeID = asduBytes[0]
	// 数据是否连续（没有做改变，我觉得和104是一样的）-------------------------------------------
	asdu.Sequence, asdu.Length = asdu.ParseVariable(asduBytes[1])
	var firstAddress uint16

	asdu.Cause = uint16(asduBytes[2])
	// asdu.PublicAddress = binary.LittleEndian.Uint16([]byte{asduBytes[3], asduBytes[4]})
	asdu.PublicAddress = binary.LittleEndian.Uint16(asduBytes[3:5])
	asdu.RecordAddress = uint(asduBytes[5])

	// 如果连续
	// if asdu.Sequence {
	// 	firstAddress = binary.LittleEndian.Uint32([]byte{asduBytes[5], asduBytes[6], asduBytes[7], 0x00})
	// }
	for i := 0; i < int(asdu.Length); i++ {
		s := new(Signal)
		s.TypeID = uint(asdu.TypeID)
		if asdu.Sequence {
			s.Address = firstAddress
			firstAddress++
		}
		switch asdu.TypeID {
		case SUBDBCur:
			// 信息体地址（２字节）＋数据内容（4字节）＋质量码（１字节）

			size := 7
			if asdu.Sequence {
				size := 5
				s.Value = float64(binary.LittleEndian.Uint32([]byte{
					asduBytes[9+i*size], asduBytes[9+i*size+1],
					asduBytes[9+i*size+2], asduBytes[9+i*size+3],
				}))
			} else {
				s.Address = binary.LittleEndian.Uint16([]byte{asduBytes[6+i*size], asduBytes[6+i*size+1]})
				// 不同的记录地址解析的数据类型不同
				switch asdu.RecordAddress {
				// 83：智能电表电能量、84：最大需量
				case 83, 84:
					s.Value = float64(binary.LittleEndian.Uint16([]byte{
						asduBytes[6+i*size+2], asduBytes[6+i*size+3],
						asduBytes[6+i*size+4], asduBytes[6+i*size+5],
					}))
				}
				s.Quality = asduBytes[12+i*size]
			}

		case MSpNa1, MDpNa1:
			size := 4
			if asdu.Sequence {
				s.Value = float64(asduBytes[9+i])
			} else {
				s.Address = binary.LittleEndian.Uint16([]byte{asduBytes[6+i*size], asduBytes[6+i*size+1], asduBytes[6+i*size+2], 0x00})
				s.Value = float64(asduBytes[6+i*size+3])
			}
		case MMeNa1:
			size := 6
			if asdu.Sequence {
				size := 3
				s.Value = float64(binary.LittleEndian.Uint16([]byte{asduBytes[9+i*size], asduBytes[9+i*size+1]}))
				s.Quality = asduBytes[9+i*size+2]
			} else {
				s.Address = binary.LittleEndian.Uint16([]byte{asduBytes[6+i*size], asduBytes[6+i*size+1], asduBytes[6+i*size+2], 0x00})
				s.Value = float64(binary.LittleEndian.Uint16([]byte{asduBytes[6+i*size+3], asduBytes[6+i*size+4]}))
				s.Quality = asduBytes[6+i*size+5]
			}

		// case MItNa1:
		// size := 8
		// if asdu.Sequence {
		// 	size := 5
		// 	s.Value = float64(binary.LittleEndian.Uint32([]byte{
		// 		asduBytes[9+i*size], asduBytes[9+i*size+1],
		// 		asduBytes[9+i*size+2], asduBytes[9+i*size+3],
		// 	}))
		// } else {
		// 	s.Address = binary.LittleEndian.Uint32([]byte{asduBytes[6+i*size], asduBytes[6+i*size+1], asduBytes[6+i*size+2], 0x00})
		// 	s.Value = float64(binary.LittleEndian.Uint32([]byte{
		// 		asduBytes[6+i*size+3], asduBytes[9+i*size+4],
		// 		asduBytes[9+i*size+5], asduBytes[9+i*size+6],
		// 	}))
		// }

		default:
			log.Fatalln("暂不支持的数据类型:", asdu.TypeID)
		}
		signals = append(signals, s)
	}
	return
}

// ParseVariable 解析asdu可变结构限定词
func (asdu *ASDU) ParseVariable(b byte) (sq bool, length byte) {
	// 最高位是否为1
	sq = b&128>>7 == 1
	if sq {
		length = b - 1<<7
		return
	}
	length = b
	return
}

// ParseTime 解析asdu中7个字节时表,转为带毫秒的时间戳
func (asdu *ASDU) ParseTime(asduBytes []byte) float64 {
	if len(asduBytes) != 7 {
		return 0
	}
	milliseconds := binary.LittleEndian.Uint16([]byte{asduBytes[0], asduBytes[1]})
	nanosecond := (int(milliseconds) % 1000) * 1000000
	second := int(milliseconds / 1000)
	minute := int(asduBytes[2])
	hour := int(asduBytes[3])
	day := int(asduBytes[4])
	month := int(asduBytes[5])
	year := int(asduBytes[6]) + 2000
	return float64(time.Date(year, time.Month(month), day, hour, minute, second, nanosecond, time.Local).Unix()) + float64(nanosecond)/1000000000.0
}
