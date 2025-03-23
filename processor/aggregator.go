package processor

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type AggregatorProcessor struct {
	// 按组聚合
	groups           map[string]*AggregatorGroup
	interval         time.Duration
	outputChan       chan map[string]interface{}
	aggregatorFields []string
}

type AggregatorGroup struct {
	//data保存分组的键值对（不含 __value__）
	data  map[string]interface{}
	sum   float64
	count int
	min   float64
	max   float64
}

func NewAggregatorProcessor(interval time.Duration, aggregatorFields []string, outputChan chan map[string]interface{}) *AggregatorProcessor {
	return &AggregatorProcessor{
		groups:           make(map[string]*AggregatorGroup),
		interval:         interval,
		outputChan:       outputChan,
		aggregatorFields: aggregatorFields,
	}
}

// getGroupKey 根据数据计算分组 key，并返回对应的 groupData（分组字段键值对）
// 如果 aggregationFields 非空，则仅取这些字段；否则取除 __value__ 外的所有字段
func (ap *AggregatorProcessor) getGroupKey(data map[string]interface{}) string {
	//groupData := make(map[string]interface{})
	var keys []string
	if len(ap.aggregatorFields) > 0 {
		for _, k := range ap.aggregatorFields {
			keys = append(keys, k)
		}
	} else {
		// 使用除 __value__ 之外的所有字段作为分组字段
		for k := range data {
			if k == "__value__" {
				continue
			}
			keys = append(keys, k)
		}
	}
	// 保证 key 的顺序一致
	sort.Strings(keys)
	var parts []string
	for _, k := range keys {
		// 将每个字段格式化为 "key=value"
		part := fmt.Sprintf("%s=%v", k, data[k])
		parts = append(parts, part)
	}
	groupKey := strings.Join(parts, "|")
	return groupKey
}

// Process 更新相应分组的聚合数据
func (ap *AggregatorProcessor) Process(data map[string]interface{}) (map[string]interface{}, error) {
	groupKey := ap.getGroupKey(data)
	value, ok := data["__value__"].(float64)
	if !ok {
		return data, fmt.Errorf("缺少 __value__ 或类型错误")
	}

	group, exists := ap.groups[groupKey]
	if !exists {
		group = &AggregatorGroup{
			data:  data,
			sum:   value,
			count: 1,
			min:   value,
			max:   value,
		}
		ap.groups[groupKey] = group
	} else {
		group.sum += value
		group.count++
		if value < group.min {
			group.min = value
		}
		if value > group.max {
			group.max = value
		}
	}

	//聚合数据不需要__value__字段
	delete(data, "__value__")
	return data, nil
}

// StartAggregation 启动定时任务，将每个分组的聚合结果分拆为 4 条记录发送到 OutputChan
func (ap *AggregatorProcessor) StartAggregation() {
	go func() {
		ticker := time.NewTicker(ap.interval * time.Second)
		for range ticker.C {
			for _, group := range ap.groups {
				// 生成各统计指标的结果记录
				data := copyMap(group.data)
				data["__value__max"] = group.max
				ap.outputChan <- data

				data = copyMap(group.data)
				data["__value__min"] = group.min
				ap.outputChan <- data

				data = copyMap(group.data)
				data["__value__sum"] = group.sum
				ap.outputChan <- data

				data = copyMap(group.data)
				data["__value__count"] = group.count
				ap.outputChan <- data
			}
			// 清空分组，开始下一轮聚合
			ap.groups = make(map[string]*AggregatorGroup)
		}
	}()
}

func copyMap(m map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}
