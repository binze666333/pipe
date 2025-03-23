package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os/exec"
	"pipe/backend"
	"pipe/generator"
	"pipe/pipeline"
	"pipe/processor"
	"runtime"
	"strings"
	"time"
)

// 处理器和通道
var (
	inputChan     = make(chan map[string]interface{}, 100)
	aggOutputChan = make(chan map[string]interface{}, 100)
	pipe          *pipeline.Pipeline
	done          = make(chan struct{})
)

func main() {
	// 初始化配置
	initConfig()
	reloadConfig()

	//将聚合后的数据保存到后端
	var dataBackend backend.Backend = backend.NewPrintBackend()
	go func() {
		for {
			select {
			case aggData := <-aggOutputChan:
				dataBackend.Send(aggData)
			}
		}
	}()

	fmt.Println("数据处理服务启动...")
	select {}
}

// 初始化Viper并监听配置文件变化
func initConfig() {
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// 启动热重载
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件发生变化，重新加载配置...")
		reloadConfig()
	})
}

// 重新加载配置
func reloadConfig() {
	//停止生成数据，重启过程如果往已关闭的channel写入数据会panic
	close(done)
	done = make(chan struct{})

	//重新创建处理器
	var process []processor.Processor
	if viper.GetBool("filter.enabled") {
		process = append(process, processor.NewFilterProcessor(viper.GetStringSlice("filter.allowed_envs")))
	}
	if viper.GetBool("fill.enabled") {
		osType := runtime.GOOS
		cmd := exec.Command("uname", "-r")
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("获取系统版本失败:", err)
		}
		osVersion := strings.TrimSpace(string(output))
		process = append(process, processor.NewFillProcessor(osType, osVersion))
	}
	if viper.GetBool("aggregator.enabled") {
		aggregator := processor.NewAggregatorProcessor(viper.GetDuration("aggregator.interval"),
			viper.GetStringSlice("aggregator.fields"), aggOutputChan)
		process = append(process, aggregator)
		aggregator.StartAggregation()
	}

	// 重新初始化流水线
	pipe = pipeline.NewPipeline(process...)

	// 重新调整 Worker 数量
	restartWorkers()

	//模拟订阅外部队列数据
	generator.StartDataGeneration(inputChan, 2000*time.Millisecond, done)
}

// 重新启动 Worker 协程
func restartWorkers() {
	// 关闭旧Worker避免goroutine泄漏
	close(inputChan)
	inputChan = make(chan map[string]interface{}, 100) // 重新创建通道

	numWorkers := runtime.NumCPU()
	//fmt.Printf("重新启动 %d 个 goroutine 处理数据\n", numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			for data := range inputChan {
				//fmt.Printf("Worker-%d 处理数据: %+v\n", workerID, data)
				pipe.Process(data)
			}
		}(i)
	}
}
