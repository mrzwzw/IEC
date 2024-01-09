package IEC

import "context"

type DataHandlerInterface interface {
	Datahandler(*APDU) error
}

// handler 处理接收到的已解析数据
func (c *Client) handler(ctx context.Context) {
	c.Logger.Info("数据处理协程启动")
	defer func() {
		c.cancel()
		c.wg.Done()
		c.Logger.Info("数据接收协程停止")
	}()
	for {
		select {
		case resp := <-c.dataChan:
			c.Logger.Debugf("接收到数据类型:%d,原因:%d,长度:%d", resp.ASDU.TypeID, resp.ASDU.Cause, len(resp.Signals))
			go c.Datahandler(resp)
		case <-ctx.Done():
			return
		}
	}
}
