package processor

type Processor interface {
	Process(data map[string]interface{}) (map[string]interface{}, error)
}
