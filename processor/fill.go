package processor

func NewFillProcessor(osType, osVersion string) *FillProcessor {
	return &FillProcessor{
		osType:    osType,
		osVersion: osVersion,
	}
}

type FillProcessor struct {
	osType    string
	osVersion string
}

func (fp *FillProcessor) Process(data map[string]interface{}) (map[string]interface{}, error) {
	// 模拟获取系统信息
	data["os_type"] = fp.osType
	data["os_version"] = fp.osVersion
	return data, nil
}
