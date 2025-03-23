package backend

import "fmt"

type PrintBackend struct{}

func NewPrintBackend() *PrintBackend {
	return &PrintBackend{}
}

func (pb *PrintBackend) Send(data map[string]interface{}) error {
	fmt.Println("发送聚合结果：", data)
	return nil
}
