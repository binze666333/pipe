package generator

import (
	"math"
	"math/rand"
	"time"
)

// 预定义 zone、bizid、env 的可能取值
var zones = []string{"gz", "sz", "bj", "sh"}
var bizids = []string{"10", "20", "30", "40"}
var envs = []string{"prod", "dev", "test"}

// GenerateData 生成单条数据
func GenerateData() map[string]interface{} {
	return map[string]interface{}{
		"zone":      zones[rand.Intn(len(zones))],
		"bizid":     bizids[rand.Intn(len(bizids))],
		"env":       envs[rand.Intn(len(envs))],
		"__value__": math.Round(rand.Float64()*1000*100) / 100, // 生成 0~1000 之间的浮点数，只保留小数点后两位
	}
}

// StartDataGeneration 持续定期生成单条数据
func StartDataGeneration(inputChan chan<- map[string]interface{}, interval time.Duration, done <-chan struct{}) {
	go func() {
		for {
			select {
			case <-done:
				return
			case inputChan <- GenerateData():
				time.Sleep(interval)
			}
		}
	}()
}
