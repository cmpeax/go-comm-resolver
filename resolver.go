package resolver

import (
	"sync"
)

var (
	obj  *Resolver
	once sync.Once
)

// 协议 接口
type IProtocol interface {
	Listen()                                    // 解析器 监听处理
	GetProtocolName() string                    // 获取绑定解析器的协议名
	Init(chan []byte, chan []byte)              // 初始化 定义
	DealData() map[uint8]func(interface{}) bool // 发送的处理函数
}

// 通讯 接口
type IComm interface {
	GetReader() chan []byte // 通讯层暴露 reader通道,向解析器投喂通讯层 接收的消息
	GetWriter() chan []byte // 通讯层暴露 writer通道,解析器投喂通讯层 将要发送的信息
}

func GetInstance() *Resolver {
	once.Do(func() {
		obj = &Resolver{}
	})
	return obj
}

// 设置解析器. 接收通讯层和协议层接口.
func (s *Resolver) SetResolver(comm IComm, tempProtocol IProtocol) {
	tempProtocol.Init(comm.GetReader(), comm.GetWriter())
	s.protocol = tempProtocol
}

type Resolver struct {
	protocol    IProtocol
	comm        *IComm
	sendHandler func(interface{}) uint8
}

func (s *Resolver) GetProtocolName() string {
	return s.protocol.GetProtocolName()
}

// 解析器 接管通讯层的信息 监听工作
func (s *Resolver) Listen() {
	s.protocol.Listen()
}

// 设置 发送时的转换函数,提取功能代号并
func (s *Resolver) SetSendHandler(fun func(interface{}) uint8) {
	s.sendHandler = fun
}

// 向 解析器发送数据.
func (s *Resolver) Send(data interface{}) {
	flag := s.sendHandler(data)

	if _, ok := s.protocol.DealData()[flag]; ok {
		s.protocol.DealData()[flag](data)
	}

}
