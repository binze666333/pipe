package processor

import "fmt"

type FilterProcessor struct {
	allowedEnvs map[string]bool
}

func NewFilterProcessor(allowed []string) *FilterProcessor {
	allowedMap := make(map[string]bool)
	for _, env := range allowed {
		allowedMap[env] = true
	}
	return &FilterProcessor{allowedEnvs: allowedMap}
}

func (fp *FilterProcessor) Process(data map[string]interface{}) (map[string]interface{}, error) {
	if env, ok := data["env"].(string); ok {
		if !fp.allowedEnvs[env] {
			// 数据不满足条件，返回 nil 表示过滤掉
			return nil, nil
		}
	} else {
		return nil, fmt.Errorf("缺少 env 字段或类型错误")
	}
	return data, nil
}
