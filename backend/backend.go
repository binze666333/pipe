package backend

// Backend 定义了数据输出接口
type Backend interface {
	Send(data map[string]interface{}) error
}
