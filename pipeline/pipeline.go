package pipeline

import (
	"fmt"
	"pipe/processor"
)

// Pipeline 负责按照顺序执行各个处理器
type Pipeline struct {
	Processors []processor.Processor
}

// NewPipeline 创建新的数据处理流水线
func NewPipeline(procs ...processor.Processor) *Pipeline {
	return &Pipeline{Processors: procs}
}

// Process 将数据依次经过流水线中的各个 Processor
func (p *Pipeline) Process(data map[string]interface{}) {
	var err error
	for _, proc := range p.Processors {
		data, err = proc.Process(data)
		if err != nil {
			fmt.Println("处理器错误:", err)
			return
		}
		if data == nil {
			// 数据被过滤掉
			return
		}
	}
}
