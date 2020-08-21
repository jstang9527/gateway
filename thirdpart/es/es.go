package es

import (
	"github.com/jstang9527/gateway/dto"
)

// SendToESChan 发给通道
func SendToESChan(at *dto.TestItem) {
	ch <- at
}
