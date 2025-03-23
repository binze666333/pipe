package backend

import "fmt"

type RedisBackend struct{}

func NewRedisBackend() *RedisBackend {
	return &RedisBackend{}
}

func (pb *RedisBackend) Send(data map[string]interface{}) error {
	fmt.Println("发送聚合结果：", data)
	return nil
}
