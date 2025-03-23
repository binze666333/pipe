package backend

import "fmt"

type MysqlBackend struct{}

func NewMysqlBackend() *MysqlBackend {
	return &MysqlBackend{}
}

func (pb *MysqlBackend) Send(data map[string]interface{}) error {
	fmt.Println("发送聚合结果：", data)
	return nil
}
