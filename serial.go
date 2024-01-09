package IEC

import (
	"io"
	"sync"
	"time"

	"github.com/goburrow/serial"
)

const (
	// SerialDefaultTimeout Serial Default timeout
	SerialDefaultTimeout = 1 * time.Second
	// SerialDefaultAutoReconnect Serial Default auto reconnect count
	SerialDefaultAutoReconnect = 0
)

type IEC102Provider struct {
	// Serial port configuration.
	serial.Config
	mu   sync.Mutex
	port io.ReadWriteCloser
}

func NewRTUClientProvider() *IEC102Provider {
	p := &IEC102Provider{}
	return p
}

func (sf *IEC102Provider) Connect() error {
	sf.mu.Lock()
	err := sf.connect()
	sf.mu.Unlock()
	return err
}

// Caller must hold the mutex before calling this method.
func (sf *IEC102Provider) connect() error {
	port, err := serial.Open(&sf.Config)
	if err != nil {
		return err
	}
	sf.port = port
	return nil
}

// IsConnected returns a bool signifying whether the client is connected or not.
func (sf *IEC102Provider) IsConnected() bool {
	sf.mu.Lock()
	b := sf.isConnected()
	sf.mu.Unlock()
	return b
}

// Caller must hold the mutex before calling this method.
func (sf *IEC102Provider) isConnected() bool {
	return sf.port != nil
}

// Close close current connection.
func (sf *IEC102Provider) Close() error {
	var err error
	sf.mu.Lock()
	if sf.port != nil {
		err = sf.port.Close()
		sf.port = nil
	}
	sf.mu.Unlock()
	return err
}
