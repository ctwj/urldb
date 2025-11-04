package debug

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ctwj/urldb/utils"
)

// DebugFileWriter 调试文件写入器
type DebugFileWriter struct {
	filePath string
	file     *os.File
	mutex    sync.Mutex
	enabled  bool
}

// NewDebugFileWriter 创建新的调试文件写入器
func NewDebugFileWriter(filePath string) *DebugFileWriter {
	writer := &DebugFileWriter{
		filePath: filePath,
		enabled:  true,
	}

	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		utils.Error("Failed to create debug log directory: %v", err)
		writer.enabled = false
		return writer
	}

	// 打开文件
	if err := writer.openFile(); err != nil {
		utils.Error("Failed to open debug log file: %v", err)
		writer.enabled = false
	}

	return writer
}

// WriteEvent 写入事件到文件
func (dw *DebugFileWriter) WriteEvent(event DebugEvent) {
	if !dw.enabled {
		return
	}

	dw.mutex.Lock()
	defer dw.mutex.Unlock()

	// 确保文件已打开
	if dw.file == nil {
		if err := dw.openFile(); err != nil {
			utils.Error("Failed to open debug log file: %v", err)
			dw.enabled = false
			return
		}
	}

	// 格式化事件
	line := dw.formatEvent(event)

	// 写入文件
	if _, err := dw.file.WriteString(line + "\n"); err != nil {
		utils.Error("Failed to write debug event to file: %v", err)
		// 尝试重新打开文件
		dw.reopenFile()
	}
}

// WriteEventJSON 写入事件为JSON格式
func (dw *DebugFileWriter) WriteEventJSON(event DebugEvent) {
	if !dw.enabled {
		return
	}

	dw.mutex.Lock()
	defer dw.mutex.Unlock()

	// 确保文件已打开
	if dw.file == nil {
		if err := dw.openFile(); err != nil {
			utils.Error("Failed to open debug log file: %v", err)
			dw.enabled = false
			return
		}
	}

	// 序列化为JSON
	jsonData, err := json.Marshal(event)
	if err != nil {
		utils.Error("Failed to marshal debug event to JSON: %v", err)
		return
	}

	// 写入文件
	if _, err := dw.file.WriteString(string(jsonData) + "\n"); err != nil {
		utils.Error("Failed to write debug event JSON to file: %v", err)
		// 尝试重新打开文件
		dw.reopenFile()
	}
}

// Close 关闭文件写入器
func (dw *DebugFileWriter) Close() {
	dw.mutex.Lock()
	defer dw.mutex.Unlock()

	if dw.file != nil {
		dw.file.Close()
		dw.file = nil
	}
}

// RotateLogFile 轮转日志文件
func (dw *DebugFileWriter) RotateLogFile() error {
	dw.mutex.Lock()
	defer dw.mutex.Unlock()

	if dw.file != nil {
		dw.file.Close()
		dw.file = nil
	}

	// 重命名当前文件
	currentFile := dw.filePath
	backupFile := fmt.Sprintf("%s.%s", currentFile, time.Now().Format("20060102_150405"))

	if _, err := os.Stat(currentFile); err == nil {
		if err := os.Rename(currentFile, backupFile); err != nil {
			return fmt.Errorf("failed to rename log file: %v", err)
		}
	}

	// 重新打开文件
	return dw.openFile()
}

// formatEvent 格式化事件
func (dw *DebugFileWriter) formatEvent(event DebugEvent) string {
	var sb strings.Builder

	// 时间戳
	sb.WriteString(event.Timestamp.Format("2006-01-02 15:04:05.000"))

	// 级别
	sb.WriteString(" [")
	sb.WriteString(string(event.Level))
	sb.WriteString("]")

	// 插件名
	sb.WriteString(" [")
	sb.WriteString(event.PluginName)
	sb.WriteString("]")

	// 事件类型
	sb.WriteString(" ")
	sb.WriteString(string(event.EventType))

	// 消息
	sb.WriteString(": ")
	sb.WriteString(event.Message)

	// 详细信息
	if len(event.Details) > 0 {
		sb.WriteString(" {")
		first := true
		for k, v := range event.Details {
			if !first {
				sb.WriteString(", ")
			}
			sb.WriteString(k)
			sb.WriteString("=")
			sb.WriteString(v)
			first = false
		}
		sb.WriteString("}")
	}

	// 持续时间
	if event.Duration > 0 {
		sb.WriteString(" (")
		sb.WriteString(event.Duration.String())
		sb.WriteString(")")
	}

	return sb.String()
}

// openFile 打开文件
func (dw *DebugFileWriter) openFile() error {
	var err error
	dw.file, err = os.OpenFile(dw.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	// 写入文件头
	header := fmt.Sprintf("# Debug log started at %s\n", time.Now().Format(time.RFC3339))
	if _, err := dw.file.WriteString(header); err != nil {
		return err
	}

	return nil
}

// reopenFile 重新打开文件
func (dw *DebugFileWriter) reopenFile() {
	if dw.file != nil {
		dw.file.Close()
		dw.file = nil
	}

	if err := dw.openFile(); err != nil {
		utils.Error("Failed to reopen debug log file: %v", err)
		dw.enabled = false
	}
}

// GetFileSize 获取文件大小
func (dw *DebugFileWriter) GetFileSize() (int64, error) {
	dw.mutex.Lock()
	defer dw.mutex.Unlock()

	if dw.file == nil {
		return 0, fmt.Errorf("file not opened")
	}

	stat, err := dw.file.Stat()
	if err != nil {
		return 0, err
	}

	return stat.Size(), nil
}

// FileExists 检查文件是否存在
func (dw *DebugFileWriter) FileExists() bool {
	_, err := os.Stat(dw.filePath)
	return err == nil
}