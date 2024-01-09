package IEC

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

// Client 102客户端
type Client struct {
	provider *IEC102Provider
	cancel   context.CancelFunc
	Logger   *logrus.Logger
	DataHandlerInterface
	dataChan chan *APDU
	sendChan chan []byte
	wg       *sync.WaitGroup
	// task func(c *APDU)
}

type FixedData struct{}

// 1、根据启动符号来判断是什么基本帧
// 固定帧长帧（启动字符为10H）
// 可变帧长帧（启动字符为68H）

func (c *Client) parseData(ctx context.Context) error {
	handleErr := func(tag string, err error) {
		c.Logger.Errorf("%s read socket读操作异常: %v", tag, err)
	}

	// 前四个字节：68h，L，L,68h
	buf := make([]byte, 1)
	// 读取启动符和长度
	n, err := c.provider.port.Read(buf)
	if err != nil {
		handleErr("读取启动符和长度", err)
		return err
	}
	// c.provider.port.SetDeadline(time.Now().Add(contextTimeout))
	if n == 0 {
		c.Logger.Info("读取到空数据,3s后继续读取数据")
		time.Sleep(3 * time.Second)
		return nil
	}
	reset, err := isFixedFrames(buf[0])
	if err != nil {
		return err
	}

	if reset {
		parseFixedFrames(ctx, c, handleErr)
	} else {
		err := parseChangeFrames(ctx, c, handleErr)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

// 解析数据根据数据 自动识别数据大小
// 解析可变帧数据
func parseChangeFrames(ctx context.Context, c *Client, handleErr func(tag string, err error)) error {
	// 前四个字节：68h，L，L,68h
	buf := make([]byte, 3)
	// 读取启动符和长度
	n, err := c.provider.port.Read(buf)
	if err != nil {
		handleErr("读取启动符和长度", err)
		return err
	}
	if n == 0 {
		c.Logger.Info("读取到空数据,3s后继续读取数据")
		time.Sleep(3 * time.Second)
		return nil
	}
	// 先拿出两个字节，第二个字节为数据的长度
	length := int(buf[0])
	// 读取正文
	contentBuf := make([]byte, length)
	n, err = c.provider.port.Read(contentBuf)
	if err != nil {
		handleErr("读取正文", err)
		return err
	}
	// 长度不够继续读取，直至达到期望长度
	i := 1
	// 其总长度为L+6，前面已经获取四个字节了，so再减去4
	for n < (length + 2) {
		i++
		nextLength := length - n
		nextBuf := make([]byte, nextLength)
		m, err := c.provider.port.Read(nextBuf)
		if err != nil {
			handleErr("循环读取正文", err)
			return err
		}
		contentBuf = append(contentBuf[:n], nextBuf[:m]...)
		n = len(contentBuf)
		c.Logger.Debugf("循环读取数据，当前为第%d次读取，期望长度:%d,本次长度:%d,当前总长度:%d", i, (length + 2), m, n)
	}

	// 前面已经将可变部分拿出来了，
	// 再之后需要基于先前部分做一个解析

	// c.Logger.Debugf("收到原始数据: [% X],rsn:%d,ssn:%d,长度:%d", append(buf, contentBuf[:n]...), c.rsn, c.ssn, 2+len(contentBuf[:n]))
	c.Logger.Debugf("收到原始数据: [% X],长度:%d", append(buf, contentBuf[:n]...), (4 + len(contentBuf[:n])))
	apdu := new(APDU)
	err = apdu.parseAPDU(contentBuf[:n])
	if err != nil {
		c.Logger.Warnf("解析APDU异常: %v", err)
		c.Logger.Panicln("退出程序")
		return err
	}

	return nil
}

// 是否为固定帧
func isFixedFrames(data byte) (bool, error) {
	d := int(data)
	if d == 16 {
		return true, nil
	} else if d == 104 {
		return false, nil
	} else {
		return false, fmt.Errorf("解析启动符出错")
	}
}

// var (
// 	reset bool
// 	ok    bool
// )

func parseFixedFrames(ctx context.Context, c *Client, handleErr func(tag string, err error)) error {
	buf := make([]byte, 5)

	// 读取启动符和长度
	n, err := c.provider.port.Read(buf)
	if err != nil {
		handleErr("读取启动符和长度", err)
		return err
	}

	if n == 0 {
		c.Logger.Info("读取到空数据,3s后继续读取数据")
		time.Sleep(3 * time.Second)
		return nil
	}

	typeID := buf[0]

	switch typeID {
	case 0:
		// 复位链路回复帧
		log.Println("复位链路应答成功")
	case 9:
		// 确认帧回复帧
		log.Println("确认帧应答成功")
	}

	return nil
}

// Write 写数据
func (c *Client) write(ctx context.Context) {
	c.Logger.Info("socket写协程启动")
	defer func() {
		c.cancel()
		c.wg.Done()
		c.Logger.Info("socket写协程停止")
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-c.sendChan:
			_, err := c.provider.port.Write(data)
			if err != nil {
				return
			}
		}
	}
}

// Read 读数据
func (c *Client) read(ctx context.Context) {
	c.Logger.Info("socket读协程启动")
	defer func() {
		c.cancel()
		c.wg.Done()
		c.Logger.Info("socket读协程停止")
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := c.parseData(ctx)
			if err != nil {
				return
			}
		}
	}
}

// NewClient 初始化客户端,连接失败，每隔10秒重试
func NewClient(handler DataHandlerInterface, logger *logrus.Logger, p *IEC102Provider) *Client {
	cli := &Client{
		provider: p,
		dataChan: make(chan *APDU, 1),
		sendChan: make(chan []byte, 1),
		Logger:   logger,
		wg:       new(sync.WaitGroup),
	}
	cli.DataHandlerInterface = handler
	return cli
}

// Run 运行
func (c *Client) Run() {
	go c.handleSignal()
	// 定时器，每15分钟发送一次总召唤
	for {

		ctx, cancel := context.WithCancel(context.Background())
		c.cancel = cancel
		c.wg.Add(3)
		go c.read(ctx)
		go c.write(ctx)
		go c.handler(ctx)

		<-ctx.Done()
		c.Logger.Info("等待goroutine退出")
		c.wg.Wait()
		if c.provider.port != nil {
			c.provider.port.Close()
		}

	}
}

// Close 结束程序
func (c *Client) handleSignal() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	c.cancel()
	c.wg.Wait()
	if c.provider.port != nil {
		err := c.provider.Close()
		log.Println("初始断开连接报错:", err)
	}
	c.Logger.Println("断开服务器连接，程序关闭")
	os.Exit(0)
}
