package task

import (
	"context"
	"fmt"

	"github.com/ctwj/urldb/db/entity"
)

// PluginTaskProcessor 插件任务处理器
type PluginTaskProcessor struct {
	taskType string
	pluginName string
	executeFunc func(ctx context.Context, taskID uint, item *entity.TaskItem) error
}

// NewPluginTaskProcessor 创建插件任务处理器
func NewPluginTaskProcessor(pluginName string, taskType string, executeFunc func(ctx context.Context, taskID uint, item *entity.TaskItem) error) *PluginTaskProcessor {
	return &PluginTaskProcessor{
		taskType:    fmt.Sprintf("plugin.%s.%s", pluginName, taskType),
		pluginName:  pluginName,
		executeFunc: executeFunc,
	}
}

// Process 处理任务项
func (p *PluginTaskProcessor) Process(ctx context.Context, taskID uint, item *entity.TaskItem) error {
	if p.executeFunc != nil {
		return p.executeFunc(ctx, taskID, item)
	}
	return fmt.Errorf("no execute function defined for plugin task processor")
}

// GetTaskType 获取任务类型
func (p *PluginTaskProcessor) GetTaskType() string {
	return p.taskType
}